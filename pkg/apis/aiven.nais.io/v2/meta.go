// Package v2 contains API Schema definitions for the aiven.nais.io v2 API group
// +kubebuilder:object:generate=true
// +groupName=aiven.nais.io
// +versionName=v2
package aiven_nais_io_v2

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/scheme"
)

var (
	// GroupVersion is group version used to register these objects
	GroupVersion = schema.GroupVersion{Group: "aiven.nais.io", Version: "v2"}
	// renamed to SchemeGroupVersion???

	// SchemeBuilder is used to add go types to the GroupVersionKind scheme
	SchemeBuilder = &scheme.Builder{GroupVersion: GroupVersion}

	// AddToScheme adds the types in this group-version to the given scheme.
	AddToScheme = SchemeBuilder.AddToScheme
)
