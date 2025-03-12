package nais_io_v1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

func init() {
	SchemeBuilder.Register(
		&Image{},
		&ImageList{},
	)
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type ImageList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Image `json:"items"`
}

// Image defines a Nais workload Image.
//
// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:printcolumn:name="Image",type="string",JSONPath=".spec.image"
// +kubebuilder:printcolumn:name="Last synchronization time",type="string",JSONPath=".status.synchronizationTime"
// +kubebuilder:validation:XValidation:rule="self.metadata.name.size() <= 63", message="metadata.name must be no more than 63 characters"
type Image struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ImageSpec   `json:"spec"`
	Status ImageStatus `json:"status,omitempty"`
}

type ImageSpec struct {
	// Image is the Docker image to deploy.
	// +required
	// +kubebuilder:validation:Pattern="^.*/.*:.*$"
	Image string `json:"image"`
}

type ImageStatus struct {
	// ObservedGeneration represents the .metadata.generation that was last observed.
	// For instance, if .metadata.generation is currently 12, but the .status.observedGeneration is 9,
	// the dependent resources are out of date with respect to the current state of the instance.
	// +optional
	// +kubebuilder:validation:Minimum=0
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`

	// SynchronizationTime is the last time the image was synchronized.
	// This should be when the dependent resource has been synchronized.
	// +required
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:Format=date-time
	SynchronizationTime metav1.Time `json:"synchronizationTime"`

	// Conditions represent the latest available observations of state
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}
