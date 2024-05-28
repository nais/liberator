package aiven_io_v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func init() {
	SchemeBuilder.Register(
		&OpenSearch{}, &OpenSearchList{},
	)
}

// Types defined here because importing them directly from aiven-operator introduces dependency resolution hell
// Copied and simplified types as of v0.12.0

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
type OpenSearch struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              OpenSearchSpec   `json:"spec,omitempty"`
	Status            OpenSearchStatus `json:"status,omitempty"`
}

type OpenSearchSpec struct {
	ServiceCommonSpec
}

// +kubebuilder:object:generate=true
type OpenSearchStatus struct {
	// Conditions represent the latest available observations of a service state
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	// Service state
	State string `json:"state,omitempty"`
}

// +kubebuilder:object:root=true
type OpenSearchList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []OpenSearch `json:"items"`
}
