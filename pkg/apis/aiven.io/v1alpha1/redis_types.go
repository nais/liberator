package aiven_io_v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func init() {
	SchemeBuilder.Register(
		&Redis{},
	)
}

// Types defined here because importing them directly from aiven-operator introduces dependency resolution hell
// Copied and simplified types as of v0.12.0

// +kubebuilder:object:root=true
type Redis struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              RedisSpec   `json:"spec,omitempty"`
	Status            RedisStatus `json:"status,omitempty"`
}

type RedisSpec struct {
	ServiceCommonSpec
}

type RedisStatus struct {
}
