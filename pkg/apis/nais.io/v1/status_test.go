package nais_io_v1_test

import (
	"testing"

	v1 "github.com/nais/liberator/pkg/apis/nais.io/v1"
	"github.com/nais/liberator/pkg/events"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type KStatusTests struct {
	State          string
	TrueConditions []string
}

func NewTestStatus(state string) *v1.Status {
	return &v1.Status{
		SynchronizationTime:     0,
		RolloutCompleteTime:     0,
		CorrelationID:           "",
		DeploymentRolloutStatus: "",
		SynchronizationState:    state,
		SynchronizationHash:     "",
		Conditions:              nil,
	}
}

func TestGenerateCorrectKStatus(t *testing.T) {

	testCases := []KStatusTests{
		{events.RolloutComplete, []string{"Ready"}},
		{events.FailedSynchronization, []string{"Stalled"}},
		{events.Synchronized, []string{"Reconciling"}},
		{events.FailedStatusUpdate, []string{"Reconciling"}},
		{events.Retrying, []string{"Reconciling"}},
		{events.FailedPrepare, []string{"Reconciling"}},
	}

	for _, test := range testCases {

		t.Run(test.State, func(t *testing.T) {
			t.Parallel()

			status := NewTestStatus(test.State)
			status.SetStatusConditions()

			var trueConditions []string
			for _, condition := range *status.Conditions {
				if condition.Status == metav1.ConditionTrue {
					trueConditions = append(trueConditions, condition.Type)
				}
			}

			assert.Equal(t, test.TrueConditions, trueConditions)
			assert.Equal(t, 3, len(*status.Conditions))
		})

	}

}
