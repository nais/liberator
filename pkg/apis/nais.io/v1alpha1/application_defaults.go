package nais_io_v1alpha1

import (
	"github.com/imdario/mergo"
	nais_io_v1 "github.com/nais/liberator/pkg/apis/nais.io/v1"
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
func ApplyDefaults(app *Application) error {
	return mergo.Merge(app, getAppDefaults())
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
				Min:                    2,
				Max:                    4,
				CpuThresholdPercentage: 50,
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
					Cpu:    "500m",
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
					Rules: []nais_io_v1.AccessPolicyRule{},
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
