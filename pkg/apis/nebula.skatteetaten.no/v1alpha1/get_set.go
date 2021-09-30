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
	return  &nais_io_v1.PrometheusConfig{
		Path: "/metrics",
	}
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
	return &nais_io_v1.Probe{
		PeriodSeconds:    nais_io_v1alpha1.DefaultProbePeriodSeconds,
		Timeout:          nais_io_v1alpha1.DefaultProbeTimeoutSeconds,
		FailureThreshold: nais_io_v1alpha1.DefaultProbeFailureThreshold,
	}
}

func (in *Application) GetReadiness() *nais_io_v1.Probe {
	return nil
}

func (in *Application) GetStartup() *nais_io_v1.Probe {
	return nil
}

func (in *Application) GetEnv() nais_io_v1.EnvVars {
	return nais_io_v1.EnvVars{}
}

func (in *Application) GetImage() string {
	return in.Spec.Pod.Image
}

func (in *Application) GetCommand() []string {
	return []string{}
}

func (in *Application) GetFilesFrom() []nais_io_v1.FilesFrom {
	return []nais_io_v1.FilesFrom{}
}

func (in *Application) GetEnvFrom() []nais_io_v1.EnvFrom {
	return []nais_io_v1.EnvFrom{}
}

func (in *Application) GetPreStopHook() *nais_io_v1.PreStopHook {
	return nil
}

func (in *Application) GetPreStopHookPath() string {
	return ""
}

func (in *Application) GetResources() *nais_io_v1.ResourceRequirements {
	return &nais_io_v1.ResourceRequirements{
		Limits: &nais_io_v1.ResourceSpec{
			Cpu:    in.Spec.Pod.Resource.Limits.Cpu().String(),
			Memory: in.Spec.Pod.Resource.Limits.Memory().String(),
		},
		Requests: &nais_io_v1.ResourceSpec{
			Cpu:    in.Spec.Pod.Resource.Requests.Cpu().String(),
			Memory: in.Spec.Pod.Resource.Requests.Memory().String(),
		},
	}
}
