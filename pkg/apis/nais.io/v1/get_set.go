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

// As Jobs don't have ingresses, they also do not have redirects
func (in *Naisjob) GetRedirects() []Redirect {
	return nil
}

func (in *Naisjob) GetImage() string {
	return in.Spec.Image
}

func (in *Naisjob) GetEnv() EnvVars {
	return in.Spec.Env
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

func (in *Naisjob) GetRedis() []Redis {
	return in.Spec.Redis
}

func (in *Naisjob) GetVault() *Vault {
	return in.Spec.Vault
}

func (in *Naisjob) GetObservability() *Observability {
	return in.Spec.Observability
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

func (in *Naisjob) GetIDPorten() *IDPorten {
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

func (in *GCP) Instance() *CloudSqlInstance {
	if in.SqlInstances == nil || len(in.SqlInstances) < 1 {
		return nil
	}
	return &in.SqlInstances[0]
}

func (in *CloudSqlInstance) Database() *CloudSqlDatabase {
	if in.Databases == nil || len(in.Databases) < 1 {
		return nil
	}
	return &in.Databases[0]
}

func (in *Naisjob) GetLogin() *Login {
	return nil
}
