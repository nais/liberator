package nais_io_v1

import (
	corev1 "k8s.io/api/core/v1"
)

func (in *Naisjob) CreateEvent(reason, message, typeStr string) *corev1.Event {
	return CreateEvent(in.CreateObjectMeta(), in.GetObjectReference(), reason, message, typeStr)
}
