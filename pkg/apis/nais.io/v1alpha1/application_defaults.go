package nais_io_v1alpha1

import (
	"github.com/imdario/mergo"
	nais_io_v1 "github.com/nais/liberator/pkg/apis/nais.io/v1"
	"github.com/nais/liberator/pkg/intutil"
)

// Application spec default values
const (
	DefaultPortName                 = "http"
	DefaultServicePort              = 80
	DefaultAppPort                  = 8080
	DefaultProbePeriodSeconds       = 10
	DefaultProbeTimeoutSeconds      = 1
	DefaultProbeFailureThreshold    = 3
	DeploymentStrategyRollingUpdate = "RollingUpdate"
	DeploymentStrategyRecreate      = "Recreate"
)

// ApplyDefaults sets default values where they are missing from an Application spec.
func (app *Application) ApplyDefaults() error {
	replicasIsZero := app.replicasDefined() && app.replicasIsZero()
	err := mergo.Merge(app, getAppDefaults())
	if err != nil {
		return err
	}

	if replicasIsZero {
		app.Spec.Replicas.Min = intutil.Intp(0)
		app.Spec.Replicas.Max = intutil.Intp(0)
	}

	return nil
}

func (app *Application) replicasDefined() bool {
	if app.Spec.Replicas != nil && app.Spec.Replicas.Min != nil && app.Spec.Replicas.Max != nil {
		return true
	}
	return false
}

func (app *Application) replicasIsZero() bool {
	return *app.Spec.Replicas.Min == 0 && *app.Spec.Replicas.Max == 0
}

func getAppDefaults() *Application {
	return &Application{
		Spec: ApplicationSpec{
			Azure: &nais_io_v1.Azure{
				Application: &nais_io_v1.AzureApplication{
					Enabled: false,
				},
			},
			Replicas: &nais_io_v1.Replicas{
				Min:                    intutil.Intp(2),
				Max:                    intutil.Intp(4),
				CpuThresholdPercentage: 50,
				DisableAutoScaling:     false,
			},
			Liveness: &nais_io_v1.Probe{
				PeriodSeconds:    DefaultProbePeriodSeconds,
				Timeout:          DefaultProbeTimeoutSeconds,
				FailureThreshold: DefaultProbeFailureThreshold,
			},
			Port: DefaultAppPort,
			Strategy: &nais_io_v1.Strategy{
				Type: DeploymentStrategyRollingUpdate,
			},
			Prometheus: &nais_io_v1.PrometheusConfig{
				Path: "/metrics",
			},
			Ingresses: []nais_io_v1.Ingress{},
			Resources: &nais_io_v1.ResourceRequirements{
				Limits: &nais_io_v1.ResourceSpec{
					Memory: "512Mi",
				},
				Requests: &nais_io_v1.ResourceSpec{
					Cpu:    "200m",
					Memory: "256Mi",
				},
			},
			Vault: &nais_io_v1.Vault{
				Enabled: false,
				Paths:   []nais_io_v1.SecretPath{},
			},
			Service: &nais_io_v1.Service{
				Port:     DefaultServicePort,
				Protocol: DefaultPortName,
			},
			SecureLogs: &nais_io_v1.SecureLogs{
				Enabled: false,
			},
			AccessPolicy: &nais_io_v1.AccessPolicy{
				Inbound: &nais_io_v1.AccessPolicyInbound{
					Rules: []nais_io_v1.AccessPolicyInboundRule{},
				},
				Outbound: &nais_io_v1.AccessPolicyOutbound{
					Rules:    []nais_io_v1.AccessPolicyRule{},
					External: []nais_io_v1.AccessPolicyExternalRule{},
				},
			},
			TokenX: &nais_io_v1.TokenX{
				Enabled:                 false,
				MountSecretsAsFilesOnly: false,
			},
		},
	}
}
