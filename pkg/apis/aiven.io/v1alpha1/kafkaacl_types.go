package aiven_io_v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func init() {
	SchemeBuilder.Register(
		&KafkaACL{}, &KafkaACLList{},
	)
}

// Types defined here because importing them directly from aiven-operator introduces dependency resolution hell
// Copied and simplified types as of v0.41.0

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:deprecatedversion:warning="Simplified KafkaACL type copied from aiven-operator, not to be installed in clusters!"
type KafkaACL struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              KafkaACLSpec   `json:"spec,omitempty"`
	Status            KafkaACLStatus `json:"status,omitempty"`
}

type KafkaACLSpec struct {
	// Identifies the project this resource belongs to
	Project string `json:"project"`

	// Specifies the name of the service that this resource belongs to
	ServiceName string `json:"serviceName"`

	// +kubebuilder:validation:Enum=admin;read;readwrite;write
	// Kafka permission to grant (admin, read, readwrite, write)
	Permission string `json:"permission"`

	// Topic name pattern for the ACL entry
	Topic string `json:"topic"`

	// Username pattern for the ACL entry
	Username string `json:"username"`
}

// +kubebuilder:object:generate=true
type KafkaACLStatus struct {
	// Conditions represent the latest available observations of a KafkaACL state
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	// Kafka ACL ID
	ID string `json:"id,omitempty"`
}

// +kubebuilder:object:root=true
type KafkaACLList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KafkaACL `json:"items"`
}
