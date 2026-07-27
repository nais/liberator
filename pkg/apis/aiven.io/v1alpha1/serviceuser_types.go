package aiven_io_v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func init() {
	SchemeBuilder.Register(func(s *runtime.Scheme) error {
		s.AddKnownTypes(GroupVersion, &ServiceUser{}, &ServiceUserList{})
		return nil
	})
}

// Types defined here because importing them directly from aiven-operator introduces dependency resolution hell
// Copied and simplified types as of v0.41.0

// ServiceUser creates a service user for accessing Aiven services.
// The resource name becomes the username in Aiven.
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:deprecatedversion:warning="Simplified ServiceUser type copied from aiven-operator, not to be installed in clusters!"
type ServiceUser struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              ServiceUserSpec   `json:"spec,omitempty"`
	Status            ServiceUserStatus `json:"status,omitempty"`
}

type ServiceUserSpec struct {
	// Identifies the project this resource belongs to
	Project string `json:"project"`

	// Specifies the name of the service that this resource belongs to
	ServiceName string `json:"serviceName"`

	// Username of the service user on Aiven. Defaults to the resource name.
	// Set this to manage users whose names are not valid Kubernetes object names (e.g. containing underscores or uppercase).
	// Immutable after creation: to change it, delete and recreate the resource.
	// At most one resource may reference a given username.
	Username string `json:"username,omitempty"`

	// Secret configuration: where the operator publishes the user's connection details
	ConnInfoSecretTarget ConnInfoSecretTarget `json:"connInfoSecretTarget,omitempty"`

	// AccessControl Service type specific access control rules for user.
	// When this block is present, the operator manages the full access-control scope it contains.
	AccessControl *ServiceUserAccessControl `json:"accessControl,omitempty"`
}

// ConnInfoSecretTarget names the secret the operator publishes connection details to.
type ConnInfoSecretTarget struct {
	// Name of the secret resource to be created. By default, it is equal to the resource name
	Name string `json:"name"`
}

// ServiceUserAccessControl is the full desired Valkey ACL snapshot managed by the operator.
// When this block is present, omitted inner fields are treated as empty lists.
type ServiceUserAccessControl struct {
	// Key access rules.
	ValkeyACLKeys []string `json:"valkeyAclKeys,omitempty"`

	// Rules for individual commands. Order matters.
	ValkeyACLCommands []string `json:"valkeyAclCommands,omitempty"`

	// Command category rules. Order matters.
	ValkeyACLCategories []string `json:"valkeyAclCategories,omitempty"`

	// Glob-style patterns defining which pub/sub channels can be accessed.
	ValkeyACLChannels []string `json:"valkeyAclChannels,omitempty"`
}

// +kubebuilder:object:generate=true
type ServiceUserStatus struct {
	// Conditions represent the latest available observations of a ServiceUser state
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	// Type of the user account
	Type string `json:"type,omitempty"`
}

// +kubebuilder:object:root=true
type ServiceUserList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ServiceUser `json:"items"`
}
