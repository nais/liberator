package nais_io_v1

import (
	"k8s.io/api/batch/v1beta1"
)

// +kubebuilder:object:generate=false
type AzureInterface interface {
	GetApplication() *AzureApplication
	GetSidecar() *AzureSidecar
}

func (in *Naisjob) GetPort() int {
	return 0
}

func (in *Naisjob) GetStatus() *Status {
	return &in.Status
}

func (in *Naisjob) SetStatus(status *Status) {
	in.Status = *status
}

func (in *Naisjob) SetStatusConditions() {
	in.Status.SetStatusConditions()
}

func (in *Naisjob) GetConcurrencyPolicy() string {
	switch in.Spec.ConcurrencyPolicy {
	case string(v1beta1.ForbidConcurrent), string(v1beta1.ReplaceConcurrent):
		return in.Spec.ConcurrencyPolicy
	default:
		return string(v1beta1.AllowConcurrent)
	}
}

func (in *Naisjob) GetAccessPolicy() *AccessPolicy {
	return in.Spec.AccessPolicy
}

func (in *Naisjob) GetAzure() AzureInterface {
	return in.Spec.Azure
}

func (in *Naisjob) GetIngress() []Ingress {
	return nil
}

func (in *Azure) GetApplication() *AzureApplication {
	return in.Application
}

func (in *Azure) GetSidecar() *AzureSidecar {
	return in.Sidecar
}

func (in *AzureNaisJob) GetApplication() *AzureApplication {
	return in.Application
}

func (in *AzureNaisJob) GetSidecar() *AzureSidecar {
	return nil
}
