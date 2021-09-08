// Package v1 contains API Schema definitions for the kafka.nais.io v1 API group
// +kubebuilder:object:generate=true
// +groupName=networking.istio.io
// +versionName=v1alpha3
package networking_istio_io_v1alpha3

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/scheme"
)

var (
	// GroupVersion is group version used to register these objects
	GroupVersion = schema.GroupVersion{Group: "networking.istio.io", Version: "v1alpha3"}
	// renamed to SchemeGroupVersion???

	// SchemeBuilder is used to add go types to the GroupVersionKind scheme
	SchemeBuilder = &scheme.Builder{GroupVersion: GroupVersion}

	// AddToScheme adds the types in this group-version to the given scheme.
	AddToScheme = SchemeBuilder.AddToScheme
)
