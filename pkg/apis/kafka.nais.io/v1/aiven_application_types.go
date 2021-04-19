package kafka_nais_io_v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func init() {
	SchemeBuilder.Register(
		&AivenApplicationList{},
		&AivenApplication{},
	)
}

// +kubebuilder:object:root=true
type AivenApplicationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AivenApplication `json:"items"`
}

// +kubebuilder:object:root=true
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:printcolumn:name="State",type="string",JSONPath=".status.synchronizationState"
// +kubebuilder:printcolumn:name="Name of secret",type=string,JSONPath=".status.secretName",priority=10
// +kubebuilder:printcolumn:name="Credentials expire",type=string,JSONPath=".spec.serviceUsers[*].credentialsExpirationTime",priority=20
type AivenApplication struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              AivenApplicationSpec    `json:"spec"`
	Status            *AivenApplicationStatus `json:"status,omitempty"`
}

type ApplicationConfig struct {
	RotationIntervalDays int `json:"rotationIntervalDays,omitempty"`
}

type ServiceUser struct {
	// +kubebuilder:validation:MaxLength=64
	Name string `json:"name"`
	// +kubebuilder:validation:Format=date-time
	CredentialsExpirationTime string `json:"credentialsExpirationTime,omitempty"`
}

type AivenApplicationSpec struct {
	Pool         string             `json:"pool"`
	ServiceUsers []ServiceUser      `json:"serviceUsers,omitempty"`
	Config       *ApplicationConfig `json:"config,omitempty"`
}

// +kubebuilder:subresource:status
type AivenApplicationStatus struct {
	// +kubebuilder:validation:Format=date-time
	SynchronizationTime  string   `json:"synchronizationTime,omitempty"`
	SynchronizationState string   `json:"synchronizationState,omitempty"`
	SynchronizationHash  string   `json:"synchronizationHash,omitempty"`
	Errors               []string `json:"errors,omitempty"`
	Message              string   `json:"message,omitempty"`
	SecretName           string   `json:"secretName,omitempty"`
	CurrentUser          string   `json:"currentUser,omitempty"`
}
