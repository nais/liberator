package iam_cnrm_cloud_google_com_v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func init() {
	SchemeBuilder.Register(
		&IAMServiceAccount{},
		&IAMServiceAccountList{},
		&IAMPolicy{},
		&IAMPolicyList{},
		&IAMPolicyMember{},
		&IAMPolicyMemberList{},
	)
}

// +kubebuilder:object:root=true
type IAMServiceAccount struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              IAMServiceAccountSpec `json:"spec"`
}

// +kubebuilder:object:root=true
type IAMServiceAccountList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []IAMServiceAccount `json:"items"`
}

type IAMServiceAccountSpec struct {
	DisplayName string `json:"displayName"`
}

// +kubebuilder:object:root=true
type IAMPolicy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              IAMPolicySpec `json:"spec"`
}

type IAMPolicySpec struct {
	ResourceRef *ResourceRef `json:"resourceRef"`
	Bindings    []Bindings   `json:"bindings"`
}

type ResourceRef struct {
	ApiVersion string  `json:"apiVersion"`
	External   *string `json:"external,omitempty"`
	Kind       string  `json:"kind"`
	Name       *string `json:"name,omitempty"`
}

type Bindings struct {
	Role    string   `json:"role"`
	Members []string `json:"members"`
}

// +kubebuilder:object:root=true
type IAMPolicyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []IAMPolicy `json:"items"`
}

// +kubebuilder:object:root=true
type IAMPolicyMember struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              IAMPolicyMemberSpec   `json:"spec"`
	Status            IAMPolicyMemberStatus `json:"status,omitempty"`
}

type IAMPolicyMemberSpec struct {
	Member      string      `json:"member"`
	Role        string      `json:"role"`
	ResourceRef ResourceRef `json:"resourceRef"`
}

// IAMPolicyMemberStatus defines the observed state of IAMPolicyMember
type IAMPolicyMemberStatus struct {
	// Conditions represent the latest available observations of the IAM
	// policy's current state.
	Conditions []metav1.Condition `json:"conditions,omitempty"`
	// ObservedGeneration is the generation of the resource that was most recently observed by the Config Connector controller.
	// If this is equal to metadata.generation, then that means that the current reported status reflects the most recent desired state of the resource.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`
}

// +kubebuilder:object:root=true
type IAMPolicyMemberList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []IAMPolicyMember `json:"items"`
}
