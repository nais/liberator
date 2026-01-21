package nais_io_v1

import (
	"github.com/imdario/mergo"
	"k8s.io/utils/ptr"
)

// Application spec default values
const (
	DefaultBackoffLimit               = 6
	DefaultFailedJobsHistoryLimit     = 1
	DefaultProbePeriodSeconds         = 10
	DefaultProbeTimeoutSeconds        = 1
	DefaultProbeFailureThreshold      = 3
	DefaultSuccessfulJobsHistoryLimit = 3
)

// ApplyDefaults sets default values where they are missing from an Application spec.
func (in *Naisjob) ApplyDefaults() error {
	noBackoffLimit := in.Spec.BackoffLimit != nil && *in.Spec.BackoffLimit == 0

	err := mergo.Merge(in, getNaisjobDefaults())
	if err != nil {
		return err
	}

	if noBackoffLimit {
		in.Spec.BackoffLimit = ptr.To(int32(0))
	}

	return nil
}

func getNaisjobDefaults() *Naisjob {
	return &Naisjob{
		Spec: NaisjobSpec{
			Azure: &AzureNaisJob{
				Application: &AzureApplication{
					Enabled: false,
				},
			},
			BackoffLimit:           ptr.To(int32(DefaultBackoffLimit)),
			FailedJobsHistoryLimit: ptr.To(int32(DefaultFailedJobsHistoryLimit)),
			Liveness: &Probe{
				PeriodSeconds:    DefaultProbePeriodSeconds,
				Timeout:          DefaultProbeTimeoutSeconds,
				FailureThreshold: DefaultProbeFailureThreshold,
			},
			Resources: &ResourceRequirements{
				Limits: &ResourceSpec{
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
