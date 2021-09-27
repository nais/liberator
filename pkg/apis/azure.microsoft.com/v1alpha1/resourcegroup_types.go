// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package azure_microsoft_com_v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ResourceGroupSpec defines the desired state of ResourceGroup
type ResourceGroupSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Location string `json:"location"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// ResourceGroup is the Schema for the resourcegroups API
// +kubebuilder:resource:shortName=rg,path=resourcegroups
// +kubebuilder:printcolumn:name="Provisioned",type="string",JSONPath=".status.provisioned"
// +kubebuilder:printcolumn:name="Message",type="string",JSONPath=".status.message"
type ResourceGroup struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ResourceGroupSpec `json:"spec,omitempty"`
	Status ASOStatus         `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ResourceGroupList contains a list of ResourceGroup
type ResourceGroupList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ResourceGroup `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ResourceGroup{}, &ResourceGroupList{})
}

