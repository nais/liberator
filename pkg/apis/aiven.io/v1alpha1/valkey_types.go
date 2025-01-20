package aiven_io_v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func init() {
	SchemeBuilder.Register(
		&Valkey{}, &ValkeyList{},
	)
}

// +kubebuilder:object:generate=true
type ValkeyStatus struct {
	// Conditions represent the latest available observations of a service state
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	// Service state
	State string `json:"state,omitempty"`
}

// Types defined here because importing them directly from aiven-operator introduces dependency resolution hell
// Copied and simplified types as of v0.12.0

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:deprecatedversion:warning="Simplified Valkey type copied from aiven-operator, not to be installed in clusters!"
type Valkey struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              ValkeySpec   `json:"spec,omitempty"`
	Status            ValkeyStatus `json:"status,omitempty"`
}

type ValkeySpec struct {
	ServiceCommonSpec `json:",inline"`
}

// +kubebuilder:object:root=true
type ValkeyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Valkey `json:"items"`
}
