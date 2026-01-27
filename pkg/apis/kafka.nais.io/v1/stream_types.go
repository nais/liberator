package kafka_nais_io_v1

import (
	"fmt"

	"github.com/nais/liberator/pkg/hash"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func init() {
	SchemeBuilder.Register(
		&Stream{},
		&StreamList{},
	)
}

// +kubebuilder:object:root=true
type StreamList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Stream `json:"items"`
}

// +kubebuilder:object:root=true
// +kubebuilder:printcolumn:name="State",type="string",JSONPath=".status.synchronizationState"
type Stream struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              StreamSpec    `json:"spec"`
	Status            *StreamStatus `json:"status,omitempty"`
}

func (in *Stream) TopicPrefix() string {
	// `_` is not a valid character in namespace or name, so it's an excellent separator.
	// We keep the convention from Topics of using <team>.<app> for compatability, but use _ for the next separator.
	return fmt.Sprintf("%s.%s_stream_", in.Namespace, in.Name)
}

func (in *Stream) TopicWildcard() string {
	return fmt.Sprintf("%s*", in.TopicPrefix())
}

func (in *Stream) GetACLs() []TopicACL {
	acls := []TopicACL{
		{
			Access:      "admin",
			Application: in.GetName(),
			Team:        in.GetNamespace(),
		},
	}

	for _, user := range in.Spec.AdditionalUsers {
		acls = append(acls, TopicACL{
			Access:      "admin",
			Application: user,
			Team:        in.GetNamespace(),
		})
	}

	return acls
}

func (in *Stream) Hash() (string, error) {
	return hash.Hash(in.Spec)
}

func (in *Stream) NeedsSynchronization(hash string) bool {
	if in.Status == nil {
		return true
	}
	return in.Status.SynchronizationHash != hash
}

type StreamStatus struct {
	SynchronizationState      string   `json:"synchronizationState,omitempty"`
	SynchronizationHash       string   `json:"synchronizationHash,omitempty"`
	SynchronizationTime       string   `json:"synchronizationTime,omitempty"`
	Errors                    []string `json:"errors,omitempty"`
	Message                   string   `json:"message,omitempty"`
	FullyQualifiedTopicPrefix string   `json:"fullyQualifiedTopicPrefix,omitempty"`
}

type StreamSpec struct {
	Pool            string   `json:"pool"`
	AdditionalUsers []string `json:"additionalUsers,omitempty"`
}
