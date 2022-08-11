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

func (in *Naisjob) GetImage() string {
	return in.Spec.Image
}

func (in *Naisjob) GetGCP() *GCP {
	return in.Spec.GCP
}

func (in *Naisjob) GetWebProxy() bool {
	return in.Spec.WebProxy
}

func (in *Naisjob) GetSkipCABundle() bool {
	return in.Spec.SkipCaBundle
}

func (in *Naisjob) GetSecureLogs() *SecureLogs {
	return in.Spec.SecureLogs
}

func (in *Naisjob) GetMaskinporten() *Maskinporten {
	return in.Spec.Maskinporten
}

func (in *Naisjob) GetInflux() *Influx {
	return in.Spec.Influx
}

func (in *Naisjob) GetKafka() *Kafka {
	return in.Spec.Kafka
}

func (in *Naisjob) GetOpenSearch() *OpenSearch {
	return in.Spec.OpenSearch
}

func (in *Naisjob) GetVault() *Vault {
	return in.Spec.Vault
}

func (in *Naisjob) GetLeaderElection() bool {
	return false
}

func (in *Naisjob) GetLiveness() *Probe {
	return in.Spec.Liveness
}

func (in *Naisjob) GetReadiness() *Probe {
	return in.Spec.Readiness
}

func (in *Naisjob) GetPrometheus() *PrometheusConfig {
	return nil
}

func (in *Naisjob) GetTerminationGracePeriodSeconds() *int64 {
	return in.Spec.TerminationGracePeriodSeconds
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
