package nais_io_v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (in *Naisjob) GetStatus() *Status {
	return &in.Status
}

func (in *Naisjob) SetStatus(status *Status) {
	in.Status = *status
}

func (in *Naisjob ) SetReadyCondition(condition metav1.ConditionStatus, reason string, message string) {
	in.Status.SetReadyCondition(condition, reason, message)
}
