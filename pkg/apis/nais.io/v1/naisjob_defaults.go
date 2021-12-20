package nais_io_v1

import (
	"github.com/imdario/mergo"
)

// Application spec default values
const (
	DefaultBackoffLimit               = 6
	DefaultCompletions                = 1
	DefaultFailedJobsHistoryLimit     = 1
	DefaultParallelism                = 1
	DefaultProbePeriodSeconds         = 10
	DefaultProbeTimeoutSeconds        = 1
	DefaultProbeFailureThreshold      = 3
	DefaultSuccessfulJobsHistoryLimit = 3
)

// ApplyDefaults sets default values where they are missing from an Application spec.
func (job *Naisjob) ApplyDefaults() error {
	return mergo.Merge(job, getNaisjobDefaults())
}

func getNaisjobDefaults() *Naisjob {
	return &Naisjob{
		Spec: NaisjobSpec{
			Azure: &AzureNaisJob{
				Application: &AzureApplication{
					Enabled: false,
				},
			},
			BackoffLimit:           DefaultBackoffLimit,
			Completions:            DefaultCompletions,
			FailedJobsHistoryLimit: DefaultFailedJobsHistoryLimit,
			Liveness: &Probe{
				PeriodSeconds:    DefaultProbePeriodSeconds,
				Timeout:          DefaultProbeTimeoutSeconds,
				FailureThreshold: DefaultProbeFailureThreshold,
			},
			Parallelism: DefaultParallelism,
			Resources: &ResourceRequirements{
				Limits: &ResourceSpec{
					Cpu:    "500m",
					Memory: "512Mi",
				},
				Requests: &ResourceSpec{
					Cpu:    "200m",
					Memory: "256Mi",
				},
			},
			RestartPolicy:              "Never",
			SuccessfulJobsHistoryLimit: DefaultSuccessfulJobsHistoryLimit,
			Vault: &Vault{
				Enabled: false,
				Paths:   []SecretPath{},
			},
			SecureLogs: &SecureLogs{
				Enabled: false,
			},
			AccessPolicy: &AccessPolicy{
				Inbound: &AccessPolicyInbound{
					Rules: []AccessPolicyInboundRule{},
				},
				Outbound: &AccessPolicyOutbound{
					Rules:    []AccessPolicyRule{},
					External: []AccessPolicyExternalRule{},
				},
			},
		},
	}
}
