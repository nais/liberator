package kafka_nais_io_v1

import (
	"fmt"
	"hash/crc32"
	"strconv"
	"strings"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"

	"github.com/nais/liberator/pkg/intutil"
)

const (
	EventRolloutComplete       = "RolloutComplete"
	EventFailedPrepare         = "FailedPrepare"
	EventFailedSynchronization = "FailedSynchronization"

	RemoveDataAnnotation = "kafka.nais.io/removeDataWhenResourceIsDeleted"

	TeamNameLength            = 20
	AppNameLength             = 30
	AivenSyncFailureThreshold = time.Hour * 12
)

func init() {
	SchemeBuilder.Register(
		&Topic{},
		&TopicList{},
	)
}

// +kubebuilder:object:root=true
type TopicList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Topic `json:"items"`
}

// +kubebuilder:object:root=true
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:printcolumn:name="State",type="string",JSONPath=".status.synchronizationState"
// +kubebuilder:printcolumn:name="Fully Qualified Name",type="string",JSONPath=".status.fullyQualifiedName"
// +kubebuilder:printcolumn:name="Credentials expiry time",type="string",JSONPath=".status.credentialsExpiryTime"
type Topic struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              TopicSpec    `json:"spec"`
	Status            *TopicStatus `json:"status,omitempty"`
}

type Config struct {
	// CleanupPolicy is either "delete" or "compact" or both.
	// This designates the retention policy to use on old log segments.
	// Defaults to `delete`.
	// +kubebuilder:validation:Enum=delete;compact;"compact,delete"
	CleanupPolicy *string `json:"cleanupPolicy,omitempty"`
	// When a producer sets acks to "all" (or "-1"), `min.insync.replicas` specifies the minimum number of replicas
	// that must acknowledge a write for the write to be considered successful.
	// Defaults to `2`.
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=7
	MinimumInSyncReplicas *int `json:"minimumInSyncReplicas,omitempty"`
	// The default number of log partitions per topic.
	// Defaults to `1`.
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=1000000
	Partitions *int `json:"partitions,omitempty"`
	// The default replication factor for created topics.
	// Defaults to `3`.
	// +kubebuilder:validation:Minimum=2
	// +optional
	Replication *int `json:"replication,omitempty"`
	// Configuration controls the maximum size a partition can grow to before we will discard old log segments
	// to free up space if we are using the "delete" retention policy. By default there is no size limit only a time limit.
	// Since this limit is enforced at the partition level, multiply it by the number of partitions to compute the topic retention in bytes.
	// Defaults to `-1`.
	RetentionBytes *int `json:"retentionBytes,omitempty"`
	// The number of hours to keep a log file before deleting it.
	// Defaults to `168` hours (1 week).
	// +kubebuilder:validation:Maximum=2562047788015
	RetentionHours *int `json:"retentionHours,omitempty"`
	// The number of hours after which Kafka will force the log to roll even if the segment file isn't full to ensure
	// that retention can delete or compact old data.
	// Defaults to `168` hours (1 week).
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=8760
	SegmentHours *int `json:"segmentHours,omitempty"`
	// The largest record batch size allowed by Kafka (after compression if compression is enabled).
	// If this is increased and there are consumers older than 0.10.2, the consumers' fetch size must also be increased
	// so that they can fetch record batches this large. In the latest message format version, records are always grouped
	// into batches for efficiency. In previous message format versions, uncompressed records are not grouped into
	// batches and this limit only applies to a single record in that case.
	// Defaults to `1048588` bytes (1 MiB/mebibyte).
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=5242880
	MaxMessageBytes *int `json:"maxMessageBytes,omitempty"`
}

// TopicSpec is a specification of the desired behavior of the topic.
type TopicSpec struct {
	Pool   string    `json:"pool"`
	Config *Config   `json:"config,omitempty"`
	ACL    TopicACLs `json:"acl"`
}

type TopicStatus struct {
	SynchronizationState   string   `json:"synchronizationState,omitempty"`
	SynchronizationHash    string   `json:"synchronizationHash,omitempty"`
	SynchronizationTime    string   `json:"synchronizationTime,omitempty"`
	CredentialsExpiryTime  string   `json:"credentialsExpiryTime,omitempty"`
	Errors                 []string `json:"errors,omitempty"`
	Message                string   `json:"message,omitempty"`
	FullyQualifiedName     string   `json:"fullyQualifiedName,omitempty"`
	LatestAivenSyncFailure string   `json:"latestAivenSyncFailure,omitempty"`
}

type TopicACLs []TopicACL

// TopicACL describes the access granted for the topic.
type TopicACL struct {
	// Access type granted for a application.
	// Defaults to `readwrite`.
	// +kubebuilder:validation:Enum=read;write;readwrite
	Access string `json:"access"`
	// The name of the specified application
	Application string `json:"application"`
	// The team of the specified application
	Team string `json:"team"`
}

type User struct {
	Username    string
	Application string
	Team        string
}

func (in Topic) RemoveDataWhenDeleted() bool {
	if in.Annotations == nil {
		return false
	}
	b, err := strconv.ParseBool(in.Annotations[RemoveDataAnnotation])
	return b && err == nil
}

func (in Topic) FullName() string {
	return in.Namespace + "." + in.Name
}

// Generate name to use for ServiceUser.
// Suffix should be "*" in ACLs, or a counter (generation) % 100 for actual usernames.
func (in TopicACL) ServiceUserNameWithSuffix(suffix string) (string, error) {
	return ServiceUserNameWithSuffix(in.Team, in.Application, suffix)
}

func (in *Topic) NeedsSynchronization(hash string) bool {
	if in.Status == nil {
		return true
	}
	if strTime := in.Status.LatestAivenSyncFailure; strTime != "" {
		parsedTime, err := time.Parse(time.RFC3339, strTime)
		if err == nil {
			threshold := time.Now().Add(-AivenSyncFailureThreshold)
			timeDiff := threshold.Sub(parsedTime)
			if timeDiff > 0 {
				return true
			}
		}
	}
	return in.Status.SynchronizationHash != hash
}

// Apply default values to Topic Config where the nil-value is not what we want
func (cfg *Config) ApplyDefaults() {
	if cfg.CleanupPolicy == nil {
		cfg.CleanupPolicy = pointer.StringPtr("delete")
	}
	if cfg.MinimumInSyncReplicas == nil {
		cfg.MinimumInSyncReplicas = intutil.Intp(2)
	}
	if cfg.Partitions == nil {
		cfg.Partitions = intutil.Intp(1)
	}
	if cfg.Replication == nil {
		cfg.Replication = intutil.Intp(3)
	}
	if cfg.RetentionBytes == nil {
		cfg.RetentionBytes = intutil.Intp(-1)
	}
	if cfg.RetentionHours == nil {
		cfg.RetentionHours = intutil.Intp(168)
	}
	if cfg.SegmentHours == nil {
		cfg.SegmentHours = intutil.Intp(168)
	}
	if cfg.MaxMessageBytes == nil {
		cfg.MaxMessageBytes = intutil.Intp(1048588)
	}
}

func ServiceUserNameWithSuffix(teamName, appName, suffix string) (string, error) {
	hash, err := hashedName(teamName, appName)
	if err != nil {
		return "", fmt.Errorf("unable to hash team and application names: %w", err)
	}
	return fmt.Sprintf("%s_%s_%s_%s", shortTeamName(teamName), shortAppName(teamName, appName), hash, suffix), nil
}

func hashedName(teamName, appName string) (string, error) {
	if strings.Contains(teamName, "*") || strings.Contains(appName, "*") {
		return "*", nil
	}
	hasher := crc32.NewIEEE()
	basename := fmt.Sprintf("%s%s", teamName, appName)
	_, err := hasher.Write([]byte(basename))
	if err != nil {
		return "", err
	}
	hashStr := fmt.Sprintf("%08x", hasher.Sum32())
	return hashStr, nil
}

func shortTeamName(team string) string {
	return shorten(team, "team", TeamNameLength)
}

func shortAppName(teamName, appName string) string {
	return shorten(appName, teamName, AppNameLength)
}

func shorten(input, prefix string, maxlen int) string {
	if strings.Contains(input, "*") && len(input) <= maxlen {
		return input
	}
	if strings.HasPrefix(input, prefix) && input != prefix {
		input = input[len(prefix):]
	}
	if len(input) > 0 {
		for input[0] == '-' {
			input = input[1:]
		}
	}
	if len(input) > maxlen {
		input = input[:maxlen]
	}
	if len(input) > 1 {
		for input[len(input)-1] == '-' {
			input = input[:len(input)-1]
		}
	}
	return input
}
