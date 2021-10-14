package nais_io_v1

import (
	"k8s.io/api/batch/v1beta1"
)

func (in *Naisjob) GetStatus() *Status {
	return &in.Status
}

func (in *Naisjob) SetStatus(status *Status) {
	in.Status = *status
}

func (in *Naisjob) GetConcurrencyPolicy() string {
	switch in.Spec.ConcurrencyPolicy {
	case string(v1beta1.ForbidConcurrent), string(v1beta1.ReplaceConcurrent):
		return in.Spec.ConcurrencyPolicy
	default:
		return string(v1beta1.AllowConcurrent)
	}
}
