package kafka_nais_io_v1

import (
	"fmt"
	"hash/crc32"
	"strconv"
	"strings"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
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
type Topic struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              TopicSpec    `json:"spec"`
	Status            *TopicStatus `json:"status,omitempty"`
}

type Config struct {
	// CleanupPolicy is either "delete" or "compact" or both.
	// This designates the retention policy to use on old log segments.
	// +nais:doc:Default="delete"
	// +nais:doc:Link="https://kafka.apache.org/33/documentation.html#topicconfigs_cleanup.policy"
	// +kubebuilder:validation:Enum=delete;compact;"compact,delete"
	CleanupPolicy *string `json:"cleanupPolicy,omitempty"`

	// When a producer sets acks to "all" (or "-1"), `min.insync.replicas` specifies the minimum number of replicas
	// that must acknowledge a write for the write to be considered successful.
	// +nais:doc:Default="2"
	// +nais:doc:Link="https://kafka.apache.org/33/documentation.html#topicconfigs_min.insync.replicas"
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=7
	MinimumInSyncReplicas *int `json:"minimumInSyncReplicas,omitempty"`

	// The default number of log partitions per topic.
	// +nais:doc:Default="1"
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=1000000
	Partitions *int `json:"partitions,omitempty"`

	// The default replication factor for created topics.
	// +nais:doc:Default="3"
	// +nais:doc:Link="https://kafka.apache.org/33/documentation.html#replication"
	// +kubebuilder:validation:Minimum=2
	Replication *int `json:"replication,omitempty"`

	// Configuration controls the maximum size a partition can grow to before we will discard old log segments
	// to free up space if we are using the "delete" retention policy. By default there is no size limit only a time limit.
	// Since this limit is enforced at the partition level, multiply it by the number of partitions to compute the topic retention in bytes.
	// +nais:doc:Default="-1"
	// +nais:doc:Link="https://kafka.apache.org/33/documentation.html#topicconfigs_retention.bytes"
	RetentionBytes *int `json:"retentionBytes,omitempty"`

	// The number of hours to keep a log file before deleting it.
	// +nais:doc:Default="168"
	// +nais:doc:Link="https://kafka.apache.org/33/documentation.html#topicconfigs_retention.ms"
	// +kubebuilder:validation:Maximum=2147483648
	RetentionHours *int `json:"retentionHours,omitempty"`

	// When set, remote storage will be used to store log segments.
	// This value controls the size of the log that is kept before it is moved to remote storage.
	// Must be less than RetentionBytes
	// +nais:doc:Default="-2"
	// +nais:doc:Link="https://kafka.apache.org/documentation.html#topicconfigs_local.retention.bytes"
	LocalRetentionBytes *int `json:"localRetentionBytes,omitempty"`

	// When set, remote storage will be used to store log segments.
	// This value controls the number of hours to keep before it is moved to remote storage.
	// Must be less than RetentionHours.
	// +nais:doc:Default="-2"
	// +nais:doc:Link="https://kafka.apache.org/documentation.html#topicconfigs_local.retention.ms"
	// +kubebuilder:validation:Maximum=2147483648
	LocalRetentionHours *int `json:"localRetentionHours,omitempty"`

	// The number of hours after which Kafka will force the log to roll even if the segment file isn't full to ensure
	// that retention can delete or compact old data.
	// +nais:doc:Default="168"
	// +nais:doc:Link="https://kafka.apache.org/33/documentation.html#topicconfigs_segment.ms"
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=8760
	SegmentHours *int `json:"segmentHours,omitempty"`

	// The largest record batch size allowed by Kafka (after compression if compression is enabled).
	// If this is increased and there are consumers older than 0.10.2, the consumers' fetch size must also be increased
	// so that they can fetch record batches this large. In the latest message format version, records are always grouped
	// into batches for efficiency. In previous message format versions, uncompressed records are not grouped into
	// batches and this limit only applies to a single record in that case.
	// +nais:doc:Default="1048588"
	// +nais:doc:Link="https://kafka.apache.org/33/documentation.html#topicconfigs_max.message.bytes"
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=5242880
	MaxMessageBytes *int `json:"maxMessageBytes,omitempty"`

	// MinCompactionLagMs indicates the minimum time a message will remain uncompacted in the log
	// +nais:doc:Default="0"
	// +nais:doc:Link="https://kafka.apache.org/33/documentation.html#topicconfigs_min.compaction.lag.ms"
	// +kubebuilder:validation:Minimum=0
	MinCompactionLagMs *int `json:"minCompactionLagMs,omitempty"`

	// MaxCompactionLagMs indicates the maximum time a message will remain ineligible for compaction in the log
	// +nais:doc:Default="Inf"
	// +nais:doc:Link="https://kafka.apache.org/33/documentation.html#topicconfigs_max.compaction.lag.ms"
	// +kubebuilder:validation:Minimum=0
	MaxCompactionLagMs *int `json:"maxCompactionLagMs,omitempty"`

	// MinCleanableDirtyRatio indicates the minimum ratio of dirty log to retention size to initiate log compaction
	// +nais:doc:Default="50%"
	// +nais:doc:Link="https://kafka.apache.org/33/documentation.html#topicconfigs_min.cleanable.dirty.ratio"
	MinCleanableDirtyRatioPercent *intstr.IntOrString `json:"minCleanableDirtyRatioPercent,omitempty"`
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
			if parsedTime.Before(threshold) {
				return true
			}
		}
	}
	return in.Status.SynchronizationHash != hash
}

func (in *Topic) ApplyDefaults() error {
	in.Spec.Config = &Config{}
	in.Spec.Config.ApplyDefaults()
	return nil
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
	if cfg.LocalRetentionBytes == nil {
		cfg.LocalRetentionBytes = intutil.Intp(-2)
	}
	if cfg.LocalRetentionHours == nil {
		cfg.LocalRetentionHours = intutil.Intp(-2)
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
