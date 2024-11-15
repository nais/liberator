package nais_io_v1alpha1

import nais_io_v1 "github.com/nais/liberator/pkg/apis/nais.io/v1"

// TODO: replace manual getters with generated code
// TODO: candidates are either `go generate` or a switch to Protobuf

var _ nais_io_v1.AivenInterface = &Application{}

func (in *Application) SetStatus(status *nais_io_v1.Status) {
	in.Status = *status
}

func (in *Application) GetStatus() *nais_io_v1.Status {
	return &in.Status
}

func (in *Application) GetStrategy() *nais_io_v1.Strategy {
	return in.Spec.Strategy
}

func (in *Application) GetReplicas() *nais_io_v1.Replicas {
	return in.Spec.Replicas
}

func (in *Application) GetPrometheus() *nais_io_v1.PrometheusConfig {
	return in.Spec.Prometheus
}

func (in *Application) GetLogtransform() string {
	return in.Spec.Logtransform
}

func (in *Application) GetLogformat() string {
	return in.Spec.Logformat
}

func (in *Application) GetPort() int {
	return in.Spec.Port
}

func (in *Application) GetService() *nais_io_v1.Service {
	return in.Spec.Service
}

func (in *Application) GetLiveness() *nais_io_v1.Probe {
	return in.Spec.Liveness
}

func (in *Application) GetReadiness() *nais_io_v1.Probe {
	return in.Spec.Readiness
}

func (in *Application) GetStartup() *nais_io_v1.Probe {
	return in.Spec.Startup
}

func (in *Application) GetEnv() nais_io_v1.EnvVars {
	return in.Spec.Env
}

func (in *Application) GetImage() string {
	return in.Spec.Image
}

func (in *Application) GetCommand() []string {
	return in.Spec.Command
}

func (in *Application) GetFilesFrom() []nais_io_v1.FilesFrom {
	return in.Spec.FilesFrom
}

func (in *Application) GetEnvFrom() []nais_io_v1.EnvFrom {
	return in.Spec.EnvFrom
}

func (in *Application) GetPreStopHook() *nais_io_v1.PreStopHook {
	return in.Spec.PreStopHook
}

func (in *Application) GetPreStopHookPath() string {
	return in.Spec.PreStopHookPath
}

func (in *Application) GetResources() *nais_io_v1.ResourceRequirements {
	return in.Spec.Resources
}

func (in *Application) GetAccessPolicy() *nais_io_v1.AccessPolicy {
	return in.Spec.AccessPolicy
}

func (in *Application) GetAzure() nais_io_v1.AzureInterface {
	return in.Spec.Azure
}

func (in *Application) GetIngress() []nais_io_v1.Ingress {
	return in.Spec.Ingresses
}

func (in *Application) GetRedirects() []nais_io_v1.Redirect {
	return in.Spec.Redirects
}

func (in *Application) GetLeaderElection() bool {
	return in.Spec.LeaderElection
}

func (in *Application) GetIDPorten() *nais_io_v1.IDPorten {
	return in.Spec.IDPorten
}

func (in *Application) GetGCP() *nais_io_v1.GCP {
	return in.Spec.GCP
}

func (in *Application) GetWebProxy() bool {
	return in.Spec.WebProxy
}

func (in *Application) GetSkipCABundle() bool {
	return in.Spec.SkipCaBundle
}

func (in *Application) GetSecureLogs() *nais_io_v1.SecureLogs {
	return in.Spec.SecureLogs
}

func (in *Application) GetMaskinporten() *nais_io_v1.Maskinporten {
	return in.Spec.Maskinporten
}

func (in *Application) GetTokenX() *nais_io_v1.TokenX {
	return in.Spec.TokenX
}

func (in *Application) GetInflux() *nais_io_v1.Influx {
	return in.Spec.Influx
}

func (in *Application) GetKafka() *nais_io_v1.Kafka {
	return in.Spec.Kafka
}

func (in *Application) GetOpenSearch() *nais_io_v1.OpenSearch {
	return in.Spec.OpenSearch
}

func (in *Application) GetRedis() []nais_io_v1.Redis {
	return in.Spec.Redis
}

func (in *Application) GetVault() *nais_io_v1.Vault {
	return in.Spec.Vault
}

func (in *Application) GetTerminationGracePeriodSeconds() *int64 {
	return in.Spec.TerminationGracePeriodSeconds
}

func (in *Application) GetFrontend() *nais_io_v1.Frontend {
	return in.Spec.Frontend
}

func (in *Application) GetObservability() *nais_io_v1.Observability {
	return in.Spec.Observability
}

func (in *Application) GetTTL() string {
	return in.Spec.TTL
}

func (in *Application) GetLogin() *nais_io_v1.Login {
	return in.Spec.Login
}
