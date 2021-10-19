package v1alpha1

import (
	nais_io_v1 "github.com/nais/liberator/pkg/apis/nais.io/v1"
	nais_io_v1alpha1 "github.com/nais/liberator/pkg/apis/nais.io/v1alpha1"
)

// TODO: replace manual getters with generated code
// TODO: candidates are either `go generate` or a switch to Protobuf

func (in *Application) SetStatus(status *nais_io_v1.Status) {
	in.Status = *status
}

func (in *Application ) SetStatusConditions() {
	in.Status.SetStatusConditions()
}

func (in *Application) GetStatus() *nais_io_v1.Status {
	return &in.Status
}

func (in *Application) GetStrategy() *nais_io_v1.Strategy {
	return &nais_io_v1.Strategy{
		Type: nais_io_v1alpha1.DeploymentStrategyRollingUpdate,
	}
}

func (in *Application) GetReplicas() *nais_io_v1.Replicas {
	return in.Spec.Replicas
}

func (in *Application) GetCleanup() *nais_io_v1.Cleanup {
	return nil
}

func (in *Application) GetPrometheus() *nais_io_v1.PrometheusConfig {
	return in.Spec.Pod.Prometheus
}

func (in *Application) GetLogtransform() string {
	return ""
}

func (in *Application) GetLogformat() string {
	return ""
}

func (in *Application) GetPort() int {
	return nais_io_v1alpha1.DefaultAppPort
}

func (in *Application) GetService() *nais_io_v1.Service {
	return &nais_io_v1.Service{
		Port:     nais_io_v1alpha1.DefaultServicePort,
		Protocol: nais_io_v1alpha1.DefaultPortName,
	}
}

func (in *Application) GetLiveness() *nais_io_v1.Probe {
	return in.Spec.Pod.Liveness
}

func (in *Application) GetReadiness() *nais_io_v1.Probe {
	return in.Spec.Pod.Readiness
}

func (in *Application) GetStartup() *nais_io_v1.Probe {
	return in.Spec.Pod.Startup
}

func (in *Application) GetEnv() nais_io_v1.EnvVars {
	return in.Spec.Pod.Env
}

func (in *Application) GetImage() string {
	return in.Spec.Pod.Image
}

func (in *Application) GetCommand() []string {
	return in.Spec.Pod.Command
}

func (in *Application) GetFilesFrom() []nais_io_v1.FilesFrom {
	return in.Spec.Pod.FilesFrom
}

func (in *Application) GetEnvFrom() []nais_io_v1.EnvFrom {
	return in.Spec.Pod.EnvFrom
}

func (in *Application) GetPreStopHook() *nais_io_v1.PreStopHook {
	return nil
}

func (in *Application) GetPreStopHookPath() string {
	return ""
}

func (in *Application) GetResources() *nais_io_v1.ResourceRequirements {
	return in.Spec.Pod.Resources
}

func (in *Application) GetIngress() *IngressConfig{
	return 	&IngressConfig{
		Public: map[string]PublicIngressConfig{
			"default": {
				Enabled: true,
				Port:    8080,
				ServicePort: 80,
				Gateway: "istio-ingressgateway",
			},
		},
	}
}

func (in *Application) GetEgress() *EgressConfig{
	return in.Spec.Egress
}

func (in *Application) GetImagePolicy() *ImagePolicyConfig {
	return in.Spec.ImagePolicy
}

func (in *Application) GetAzureResourceGroup() string {
	return in.Spec.Azure.ResourceGroup
}

func (in *Application) GetPostgresDatabases() []*PostgreDatabaseConfig {
	return in.Spec.Azure.PostgreDatabases
}

func (in *Application) GetStorageAccounts() map[string]*StorageAccountConfig {
	return in.Spec.Azure.StorageAccount
}