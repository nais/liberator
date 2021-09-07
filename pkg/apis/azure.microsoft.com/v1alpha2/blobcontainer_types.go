// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package azure_microsoft_com_v1alpha2

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// PublicAccess enumerates the values for public access.
type PublicAccess string

const (
	// PublicAccessBlob ...
	PublicAccessBlob PublicAccess = "Blob"
	// PublicAccessContainer ...
	PublicAccessContainer PublicAccess = "Container"
	// PublicAccessNone ...
	PublicAccessNone PublicAccess = "None"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// BlobContainerSpec defines the desired state of BlobContainer
type BlobContainerSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Location string `json:"location"`
	// +kubebuilder:validation:Pattern=^[-\w\._\(\)]+$
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:Required
	ResourceGroup string         `json:"resourceGroup"`
	AccountName   string         `json:"accountName,omitempty"`
	AccessLevel   PublicAccess `json:"accessLevel,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:storageversion

// BlobContainer is the Schema for the blobcontainers API
// +kubebuilder:resource:shortName=bc
// +kubebuilder:printcolumn:name="Provisioned",type="string",JSONPath=".status.provisioned"
// +kubebuilder:printcolumn:name="Message",type="string",JSONPath=".status.message"
type BlobContainer struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   BlobContainerSpec `json:"spec,omitempty"`
	Status ASOStatus         `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// BlobContainerList contains a list of BlobContainer
type BlobContainerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []BlobContainer `json:"items"`
}

func init() {
	SchemeBuilder.Register(&BlobContainer{}, &BlobContainerList{})
}

func (bc *BlobContainer) IsSubmitted() bool {
	return bc.Status.Provisioned || bc.Status.Provisioning
}

func (bc *BlobContainer) IsProvisioned() bool {
	return bc.Status.Provisioned
}
