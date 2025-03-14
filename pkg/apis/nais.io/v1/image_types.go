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
}
