// +kubebuilder:skip
package aiven_io_v1alpha1

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/scheme"
)

// Types defined here because importing them directly from aiven-operator introduces dependency resolution hell
// Copied and simplified types as of v0.12.0

var (
	// GroupVersion is group version used to register these objects
	GroupVersion = schema.GroupVersion{Group: "aiven.io", Version: "v1alpha1"}
	// renamed to SchemeGroupVersion???

	// SchemeBuilder is used to add go types to the GroupVersionKind scheme
	SchemeBuilder = &scheme.Builder{GroupVersion: GroupVersion}

	// AddToScheme adds the types in this group-version to the given scheme.
	AddToScheme = SchemeBuilder.AddToScheme
)
