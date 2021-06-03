package nais_io_v1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

func (in *Naisjob) CreateObjectMeta() metav1.ObjectMeta {
	return CreateObjectMeta(in.Name, in.Namespace, in.CorrelationID(), in.Labels, in.OwnerReferences(in))
}

func (in *Naisjob) CreateAppNamespaceHash() string {
	return CreateAppNamespaceHash(in.Name, in.Namespace)
}

func (in *Naisjob) CreateObjectMetaWithName(name string) metav1.ObjectMeta {
	m := in.CreateObjectMeta()
	m.Name = name
	return m
}

func (in *Naisjob) OwnerReferences(naisjob *Naisjob) []metav1.OwnerReference {
	return []metav1.OwnerReference{naisjob.GetOwnerReference()}
}
