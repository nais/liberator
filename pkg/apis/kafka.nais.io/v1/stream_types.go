package kafka_nais_io_v1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

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

type StreamStatus struct {
	SynchronizationState      string   `json:"synchronizationState,omitempty"`
	SynchronizationHash       string   `json:"synchronizationHash,omitempty"`
	SynchronizationTime       string   `json:"synchronizationTime,omitempty"`
	Errors                    []string `json:"errors,omitempty"`
	Message                   string   `json:"message,omitempty"`
	FullyQualifiedTopicPrefix string   `json:"fullyQualifiedTopicPrefix,omitempty"`
}

type StreamSpec struct {
	Pool string `json:"pool"`
}
