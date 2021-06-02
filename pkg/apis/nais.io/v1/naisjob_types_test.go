package nais_io_v1_test

import (
	"testing"

	nais_io_v1 "github.com/nais/liberator/pkg/apis/nais.io/v1"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestNaisjobHash(t *testing.T) {
	actual, err := minimalNaisjob().Hash()
	assert.NoError(t, err)
	assert.Equal(t, "39091b68783a6373", actual)
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
