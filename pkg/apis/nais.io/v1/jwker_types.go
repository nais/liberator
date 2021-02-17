package nais_io_v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	EventRolloutComplete       = "RolloutComplete"
	EventFailedPrepare         = "FailedPrepare"
	EventFailedSynchronization = "FailedSynchronization"
)

func init() {
	SchemeBuilder.Register(
		&Jwker{},
		&JwkerList{},
	)
}

type JwkerSpec struct {
	AccessPolicy *AccessPolicy `json:"accessPolicy"` // fixme: access policy should not have rules required, but cluster and namespace. doesn't need external.
	SecretName   string        `json:"secretName"`
}

// JwkerStatus defines the observed state of Jwker
type JwkerStatus struct {
	SynchronizationTime       int64  `json:"synchronizationTime,omitempty"`
	SynchronizationState      string `json:"synchronizationState,omitempty"`
	SynchronizationHash       string `json:"synchronizationHash,omitempty"`
	SynchronizationSecretName string `json:"synchronizationSecretName,omitempty"`
}

// +genclient
// +kubebuilder:printcolumn:name="Secret",type="string",JSONPath=".spec.secretName"
// +kubebuilder:object:root=true

// Jwker is the Schema for the jwkers API
type Jwker struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   JwkerSpec   `json:"spec,omitempty"`
	Status JwkerStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// JwkerList contains a list of Jwker
type JwkerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Jwker `json:"items"`
}
