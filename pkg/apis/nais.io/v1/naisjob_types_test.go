package nais_io_v1_test

import (
	"testing"

	nais_io_v1 "github.com/nais/liberator/pkg/apis/nais.io/v1"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	// Change this value to accept re-synchronization of ALL naisjob resources when deploying a new version.
	naisjobHash     = "9de7bf5df82dace9"
	aivenGeneration = 0
)

func TestNaisjobHash(t *testing.T) {
	job := minimalNaisjob()
	err := job.ApplyDefaults()
	if err != nil {
		panic(err)
	}
	hash, err := job.Hash(aivenGeneration)
	assert.NoError(t, err)
	assert.Equalf(t, naisjobHash, hash, "Your Naisjob default value changes will trigger a FULL REDEPLOY of ALL NAISJOBS in ALL NAMESPACES across ALL CLUSTERS. If this is what you really want, change the `naisjobHash` constant in this test file to `%s`.", hash)
}

func minimalNaisjob() *nais_io_v1.Naisjob {
	return &nais_io_v1.Naisjob{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-app",
			Namespace: "test-namespace",
		},
		Spec: nais_io_v1.NaisjobSpec{
			Schedule: "test * * * * :)",
		},
		Status: nais_io_v1.Status{
			SynchronizationTime:     0,
			RolloutCompleteTime:     0,
			CorrelationID:           "test",
			DeploymentRolloutStatus: "test",
			SynchronizationState:    "test",
			SynchronizationHash:     "test",
		},
	}
}
