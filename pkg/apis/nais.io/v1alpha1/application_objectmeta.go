package nais_io_v1alpha1

import (
	nais_io_v1 "github.com/nais/liberator/pkg/apis/nais.io/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (in *Application) CreateObjectMeta() metav1.ObjectMeta {
	return nais_io_v1.CreateObjectMeta(in.Name, in.Namespace, in.CorrelationID(), in.Labels, in.OwnerReferences(in))
}

// We concatenate name, namespace and add a hash in order to avoid duplicate names when creating service accounts in common service accounts namespace.
// Also making sure to not exceed name length restrictions of 30 characters
func (in *Application) CreateAppNamespaceHash() string {
	return nais_io_v1.CreateAppNamespaceHash(in.Name, in.Namespace)
}

func (in *Application) CreateObjectMetaWithName(name string) metav1.ObjectMeta {
	m := in.CreateObjectMeta()
	m.Name = name
	return m
}

func (in *Application) OwnerReferences(app *Application) []metav1.OwnerReference {
	return []metav1.OwnerReference{app.GetOwnerReference()}
}
