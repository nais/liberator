package aiven_io_v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func init() {
	SchemeBuilder.Register(
		&OpenSearchACLConfig{}, &OpenSearchACLConfigList{},
	)
}

// Types defined here because importing them directly from aiven-operator introduces dependency resolution hell
// Copied and simplified types as of v0.41.0

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:deprecatedversion:warning="Simplified OpenSearchACLConfig type copied from aiven-operator, not to be installed in clusters!"
type OpenSearchACLConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              OpenSearchACLConfigSpec   `json:"spec,omitempty"`
	Status            OpenSearchACLConfigStatus `json:"status,omitempty"`
}

// OpenSearchACLConfigSpec is the complete ACL configuration for one OpenSearch
// service: aiven-operator removes any entry missing from Acls and reverts
// out-of-band changes, so writers must preserve entries they do not own.
type OpenSearchACLConfigSpec struct {
	// Identifies the project this resource belongs to
	Project string `json:"project"`

	// Specifies the name of the service that this resource belongs to
	ServiceName string `json:"serviceName"`

	// Enable OpenSearch ACLs. When disabled, authenticated service users have unrestricted access
	Enabled bool `json:"enabled"`

	// List of OpenSearch ACLs
	// +listType=map
	// +listMapKey=username
	Acls []OpenSearchACLConfigACL `json:"acls,omitempty"`
}

// OpenSearchACLConfigACL is one service user's ACL entry.
type OpenSearchACLConfigACL struct {
	// +kubebuilder:validation:MinLength=1
	// Username
	Username string `json:"username"`

	// +kubebuilder:validation:Required
	// OpenSearch rules
	Rules []OpenSearchACLConfigRule `json:"rules"`
}

// OpenSearchACLConfigRule grants a permission on an index pattern.
type OpenSearchACLConfigRule struct {
	// +kubebuilder:validation:MinLength=1
	// OpenSearch index pattern
	Index string `json:"index"`

	// +kubebuilder:validation:Enum=admin;deny;read;readwrite;write
	// OpenSearch permission
	Permission string `json:"permission"`
}

// +kubebuilder:object:generate=true
type OpenSearchACLConfigStatus struct {
	// Conditions represent the latest available observations of an OpenSearchACLConfig state
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
type OpenSearchACLConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []OpenSearchACLConfig `json:"items"`
}
