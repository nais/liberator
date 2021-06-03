package nais_io_v1alpha1

import (
	nais_io_v1 "github.com/nais/liberator/pkg/apis/nais.io/v1"
	corev1 "k8s.io/api/core/v1"
)

func (in *Application) CreateEvent(reason, message, typeStr string) *corev1.Event {
	return nais_io_v1.CreateEvent(in.CreateObjectMeta(), in.GetObjectReference(), reason, message, typeStr)
}
