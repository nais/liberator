package aiven_nais_io_v1

import (
	"time"

	"github.com/nais/liberator/pkg/strings"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func init() {
	SchemeBuilder.Register(
		&AivenApplicationList{},
		&AivenApplication{},
	)
}

const (
	MaxServiceUserNameLength = 64
)

// +kubebuilder:object:root=true
type AivenApplicationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AivenApplication `json:"items"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:shortName={"aivenapp"}
// +kubebuilder:subresource:status
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

type OpenSearchSpec struct {
	// Use the `instance_name` that you specified in the [navikt/aiven-iac](https://github.com/navikt/aiven-iac) repository.
	Instance string `json:"instance,omitempty"`
	// Access level for opensearch user
	// +kubebuilder:validation:Enum=read;write;readwrite;admin
	Access string `json:"access,omitempty"`
	// SecretName is the name of the secret containing Aiven credentials for the OpensSearch serviceuser
	SecretName string `json:"secretName,omitempty"`
}

type ValkeySpec struct {
	// The last part of the name used when creating the instance (ie. valkey-<team>-<instance>)
	Instance string `json:"instance,omitempty"`
	// Access level for Valkey user
	// +kubebuilder:validation:Enum=read;write;readwrite;admin
	Access string `json:"access,omitempty"`
	// SecretName is the name of the secret containing Aiven credentials for the Valkey serviceuser
	SecretName string `json:"secretName,omitempty"`
}

type InfluxDBSpec struct {
	// Name of the InfluxDB instance (`influx-<team>`)
	Instance string `json:"instance,omitempty"`
}

type AivenApplicationSpec struct {
	// SecretName is the name of the secret containing Aiven credentials
	SecretName string `json:"secretName"`
	// A Protected secret will not be deleted by the janitor even when not in use
	Protected bool `json:"protected,omitempty"`
	// A timestamp that indicates time-to-expire-date for personal secrets.
	// Format RFC3339 = "2006-01-02T15:04:05Z07:00"
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:Format=date-time
	ExpiresAt *metav1.Time `json:"expiresAt,omitempty"`
	// Kafka is a section configuring the kafka credentials to provision
	Kafka *KafkaSpec `json:"kafka,omitempty"`
	// OpenSearch is a section configuring the OpenSearch credentials to provision
	OpenSearch *OpenSearchSpec `json:"openSearch,omitempty"`
	// Valkey is a section configuring the Valkey credentials to provision
	Valkey []*ValkeySpec `json:"valkey,omitempty"`
	// InfluxDB is a section configuring the InfluxDB credentials to provision
	InfluxDB *InfluxDBSpec `json:"influxDB,omitempty"`
}

type KafkaSpec struct {
	// Pool is the Kafka pool (aka cluster) on Aiven this application uses
	Pool string `json:"pool"`
}

type AivenApplicationConditionType string

const (
	AivenApplicationSucceeded    AivenApplicationConditionType = "Succeeded"
	AivenApplicationAivenFailure AivenApplicationConditionType = "AivenFailure"
	AivenApplicationLocalFailure AivenApplicationConditionType = "LocalFailure"
)

// AivenApplicationCondition describes the state of a deployment at a certain point.
type AivenApplicationCondition struct {
	// Type of condition.
	Type AivenApplicationConditionType `json:"type"`
	// Status of the condition, one of True, False, Unknown.
	Status corev1.ConditionStatus `json:"status"`
	// The last time this condition was updated.
	LastUpdateTime metav1.Time `json:"lastUpdateTime,omitempty"`
	// The reason for the condition's last transition.
	Reason string `json:"reason,omitempty"`
	// A human readable message indicating details about the transition.
	Message string `json:"message,omitempty"`
}

type AivenApplicationStatus struct {
	// SynchronizationHash is the hash of the AivenApplication object most recently successfully synchronized
	SynchronizationHash string `json:"synchronizationHash,omitempty"`
	// SynchronizationSecretName is the SecretName set in the last successful synchronization
	SynchronizationSecretName string `json:"synchronizationSecretName,omitempty"`
	// SynchronizationState denotes whether the provisioning of the AivenApplication has been successfully completed or not
	SynchronizationState string `json:"synchronizationState,omitempty"`
	// SynchronizationTime is the last time the Status subresource was updated
	SynchronizationTime *metav1.Time `json:"synchronizationTime,omitempty"`
	// SynchronizedGeneration is the generation most recently successfully synchronized by Aivenator
	SynchronizedGeneration int64 `json:"synchronizedGeneration,omitempty"`
	// ObservedGeneration is the generation most recently observed by Aivenator
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`
	// Represents the latest available observations of an AivenApplications' current state.
	Conditions []AivenApplicationCondition `json:"conditions,omitempty"`
}

func (in *AivenApplication) GetOwnerReference() metav1.OwnerReference {
	return metav1.OwnerReference{
		APIVersion: in.APIVersion,
		Kind:       in.Kind,
		Name:       in.Name,
		UID:        in.UID,
	}
}

func (in *AivenApplicationStatus) AddCondition(condition AivenApplicationCondition, dropTypes ...AivenApplicationConditionType) {
	var dropTypeStrings []string
	for _, dropType := range dropTypes {
		dropTypeStrings = append(dropTypeStrings, string(dropType))
	}
	condition.LastUpdateTime = metav1.Time{time.Now()}
	conditions := make([]AivenApplicationCondition, 0, len(in.Conditions))
	for _, c := range in.Conditions {
		if strings.ContainsString(dropTypeStrings, string(c.Type)) {
			continue
		}
		if c.Type != condition.Type {
			conditions = append(conditions, c)
		}
	}
	conditions = append(conditions, condition)
	in.Conditions = conditions
}

func (in *AivenApplicationStatus) GetConditionOfType(conditionType AivenApplicationConditionType) *AivenApplicationCondition {
	for _, condition := range in.Conditions {
		if condition.Type == conditionType {
			return &condition
		}
	}
	return nil
}

func (in *AivenApplication) FormatExpiresAt() string {
	return in.Spec.ExpiresAt.Format(time.RFC3339)
}

func (in *AivenApplication) SecretKey() client.ObjectKey {
	return client.ObjectKey{
		Namespace: in.GetNamespace(),
		Name:      in.Spec.SecretName,
	}
}
