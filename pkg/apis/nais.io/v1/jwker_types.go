package nais_io_v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
	// ClientID is the OAuth2 Client ID associated with this Jwker
	ClientID string `json:"clientID,omitempty"`
	// KeyIDs are the JWK Key IDs (kid) associated with this Jwker
	KeyIDs []string `json:"keyIDs,omitempty"`
	// ObservedGeneration is the most recent generation observed by the controller.
	// If this is equal to metadata.generation, then that means that the current reported status reflects the most recent desired state of the resource.
	ObservedGeneration        int64       `json:"observedGeneration,omitempty"`
	SynchronizationTimestamp  metav1.Time `json:"synchronizationTimestamp,omitempty"`
	SynchronizationState      string      `json:"synchronizationState,omitempty"`
	SynchronizationSecretName string      `json:"synchronizationSecretName,omitempty"`
}

// +genclient
// +kubebuilder:printcolumn:name="Created",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:printcolumn:name="Synced",type="date",JSONPath=".status.synchronizationTimestamp"
// +kubebuilder:printcolumn:name="State",type="string",JSONPath=".status.synchronizationState"
// +kubebuilder:printcolumn:name="Client ID",type="string",JSONPath=".status.clientID",priority=1
// +kubebuilder:printcolumn:name="Secret",type="string",JSONPath=".status.synchronizationSecretName",priority=1
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

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
