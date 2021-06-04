package nais_io_v1_test

import (
	"testing"

	nais_io_v1 "github.com/nais/liberator/pkg/apis/nais.io/v1"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	naisjobHash = "123724db529e6aa4"
)

func TestNaisjobHash(t *testing.T) {
	job := minimalNaisjob()
	err := nais_io_v1.ApplyNaisjobDefaults(job)
	if err != nil {
		panic(err)
	}
	hash, err := job.Hash()
	assert.NoError(t, err)
	assert.Equalf(t, naisjobHash, hash, "Your Naisjob default value changes will trigger a FULL REDEPLOY of ALL NAISJOBS in ALL NAMESPACES across ALL CLUSTERS. If this is what you really want, change the `naisjobHash` constant in this test file to `%s`.", hash)

}

func minimalNaisjob() *nais_io_v1.Naisjob {
	return &nais_io_v1.Naisjob{
		ObjectMeta: metav1.ObjectMeta{
			Name:        "test-app",
			Namespace:   "test-namespace",
			ClusterName: "test-cluster",
		},
		Spec: nais_io_v1.NaisjobSpec{
			Schedule: "test * * * * :)",
		},
		Status: nais_io_v1.NaisjobStatus{
			SynchronizationTime:     0,
			RolloutCompleteTime:     0,
			CorrelationID:           "test",
			DeploymentRolloutStatus: "test",
			SynchronizationState:    "test",
			SynchronizationHash:     "test",
		},
	}
}
