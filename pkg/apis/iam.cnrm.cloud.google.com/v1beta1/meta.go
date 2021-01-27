// Package v1 contains API Schema definitions for the kafka.nais.io v1 API group
// +kubebuilder:object:generate=true
// +groupName=iam.cnrm.cloud.google.com
// +versionName=v1beta1
package iam_cnrm_cloud_google_com_v1beta1

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/scheme"
)

var (
	// GroupVersion is group version used to register these objects
	GroupVersion = schema.GroupVersion{Group: "iam.cnrm.cloud.google.com", Version: "v1beta1"}

	// SchemeBuilder is used to add go types to the GroupVersionKind scheme
	SchemeBuilder = &scheme.Builder{GroupVersion: GroupVersion}

	// AddToScheme adds the types in this group-version to the given scheme.
	AddToScheme = SchemeBuilder.AddToScheme
)
