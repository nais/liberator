package kafka_nais_io_v1

import (
	"fmt"
	aiven_nais_io_v1 "github.com/nais/liberator/pkg/apis/aiven.nais.io/v1"
	"strconv"

	"github.com/nais/liberator/pkg/namegen"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	EventRolloutComplete       = "RolloutComplete"
	EventFailedPrepare         = "FailedPrepare"
	EventFailedSynchronization = "FailedSynchronization"

	MaxServiceUserNameLength = 40

	Finalizer            = "kafkarator.kafka.nais.io"
	RemoveDataAnnotation = "kafka.nais.io/removeDataWhenResourceIsDeleted"
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
	CleanupPolicy         *string `json:"cleanupPolicy,omitempty"`
	MinimumInSyncReplicas *int    `json:"minimumInSyncReplicas,omitempty"`
	Partitions            *int    `json:"partitions,omitempty"`
	Replication           *int    `json:"replication,omitempty"`
	RetentionBytes        *int    `json:"retentionBytes,omitempty"`
	RetentionHours        *int    `json:"retentionHours,omitempty"`
}

type TopicSpec struct {
	Pool   string    `json:"pool"`
	Config *Config   `json:"config,omitempty"`
	ACL    TopicACLs `json:"acl"`
}

type TopicStatus struct {
	SynchronizationState  string   `json:"synchronizationState,omitempty"`
	SynchronizationHash   string   `json:"synchronizationHash,omitempty"`
	SynchronizationTime   string   `json:"synchronizationTime,omitempty"`
	CredentialsExpiryTime string   `json:"credentialsExpiryTime,omitempty"`
	Errors                []string `json:"errors,omitempty"`
	Message               string   `json:"message,omitempty"`
	FullyQualifiedName    string   `json:"fullyQualifiedName,omitempty"`
}

type TopicACLs []TopicACL

type TopicACL struct {
	// +kubebuilder:validation:Enum=read;write;readwrite
	Access      string `json:"access"`
	Application string `json:"application"`
	Team        string `json:"team"`
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

func (in *Topic) AppendFinalizer() {
	if in.Finalizers == nil {
		in.Finalizers = make([]string, 0)
	}
	for _, v := range in.Finalizers {
		if v == Finalizer {
			return
		}
	}
	in.Finalizers = append(in.Finalizers, Finalizer)
}

func (in *Topic) RemoveFinalizer() {
	finalizers := make([]string, 0, len(in.Finalizers))
	for _, v := range in.Finalizers {
		if v != Finalizer {
			finalizers = append(finalizers, v)
		}
	}
	in.Finalizers = finalizers
}

func (in Topic) FullName() string {
	return in.Namespace + "." + in.Name
}

func (in TopicACL) Username() string {
	username := in.Team + "." + in.Application
	username, err := namegen.ShortName(username, MaxServiceUserNameLength)
	if err != nil {
		panic(err)
	}
	return username
}

func (in TopicACL) ACLname() string {
	// TODO: Use new max length when Aivenator takes over creation of service users
	return fmt.Sprintf("%s*", aiven_nais_io_v1.ServiceUserPrefix(in.Application, in.Team, MaxServiceUserNameLength))
}

func (in TopicACL) User() User {
	return User{
		Username:    in.Username(),
		Application: in.Application,
		Team:        in.Team,
	}
}

func (in TopicACLs) Users() []User {
	users := make(map[User]interface{})
	result := make([]User, 0, len(in))
	for _, acl := range in {
		users[acl.User()] = new(interface{})
	}
	for k := range users {
		result = append(result, k)
	}
	return result
}

func (in *Topic) NeedsSynchronization(hash string) bool {
	if in.Status == nil {
		return true
	}
	return in.Status.SynchronizationHash != hash
}
