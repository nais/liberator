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
// +kubebuilder:resource:shortName={"aivenapp"}
// +kubebuilder:printcolumn:name="Name of secret",type=string,JSONPath=".spec.secretName"
// +kubebuilder:printcolumn:name="State",type="string",JSONPath=".status.synchronizationState",priority=10
// +kubebuilder:printcolumn:name="Synced",type="date",JSONPath=".status.synchronizationTime",priority=20
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp",priority=30
type AivenApplication struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              AivenApplicationSpec   `json:"spec,omitempty"`
	Status            AivenApplicationStatus `json:"status,omitempty"`
}

type AivenApplicationSpec struct {
	// SecretName is the name of the secret containing Aiven credentials
	SecretName string `json:"secretName"`
	// Kafka is a section configuring the kafka credentials to provision
	Kafka KafkaSpec `json:"kafka,omitempty"`
}

type KafkaSpec struct {
	// Pool is the Kafka pool (aka cluster) on Aiven this application uses
	Pool string `json:"pool"`
}

// +kubebuilder:subresource:status
type AivenApplicationStatus struct {
	// SynchronizationHash is the hash of the AivenApplication object
	SynchronizationHash string `json:"synchronizationHash,omitempty"`
	// SynchronizationSecretName is the SecretName set in the last successful synchronization
	SynchronizationSecretName string `json:"synchronizationSecretName,omitempty"`
	// SynchronizationState denotes whether the provisioning of the AivenApplication has been successfully completed or not
	SynchronizationState string `json:"synchronizationState,omitempty"`
	// SynchronizationTime is the last time the Status subresource was updated
	SynchronizationTime *metav1.Time `json:"synchronizationTime,omitempty"`
}

func (in *AivenApplication) GetOwnerReference() metav1.OwnerReference {
	return metav1.OwnerReference{
		APIVersion: in.APIVersion,
		Kind:       in.Kind,
		Name:       in.Name,
		UID:        in.UID,
	}
}
