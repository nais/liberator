package nais_io_v1_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	v1 "github.com/nais/liberator/pkg/apis/nais.io/v1"
	"github.com/nais/liberator/pkg/events"
)

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

func TestSetCondition(t *testing.T) {
	expected := metav1.Condition{
		Type:    "StatusType",
		Status:  "True",
		Reason:  "Reezon",
		Message: "MessH",
	}

	status := NewTestStatus("NOSTATE")

	status.SetCondition(expected.Type, metav1.ConditionTrue, expected.Reason, expected.Message)

	// timestamp is always updated
	assert.NotEqual(t, time.Time{}, (*status.Conditions)[0].LastTransitionTime)
	expected.LastTransitionTime = (*status.Conditions)[0].LastTransitionTime

	assert.Len(t, *status.Conditions, 1)
	assert.Equal(t, []metav1.Condition{expected}, *status.Conditions)
}

func TestSetSynchronizationStateWithCondition(t *testing.T) {
	const typeStr = "SynchronizationState"
	const statusStr = "True"
	const okMsg = "well done, sir"
	const errMsg = "situation normal; all fudged up"

	// Initialize empty state
	status := NewTestStatus("NOSTATE")

	// Set condition for the first time
	status.SetSynchronizationStateWithCondition(events.RolloutComplete, okMsg)

	// Test against expectations
	expected := metav1.Condition{
		Type:               typeStr,
		Status:             statusStr,
		Message:            okMsg,
		Reason:             "RolloutComplete",
		LastTransitionTime: (*status.Conditions)[0].LastTransitionTime,
	}
	assert.Equal(t, []metav1.Condition{expected}, *status.Conditions)

	// Set condition for the second time
	status.SetSynchronizationStateWithCondition(events.FailedSynchronization, errMsg)

	// Test against expectations
	expected = metav1.Condition{
		Type:               typeStr,
		Status:             statusStr,
		Message:            errMsg,
		Reason:             "FailedSynchronization",
		LastTransitionTime: (*status.Conditions)[0].LastTransitionTime,
	}
	assert.Equal(t, []metav1.Condition{expected}, *status.Conditions)
}
