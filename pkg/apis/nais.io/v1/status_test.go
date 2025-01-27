package nais_io_v1_test

import (
	"testing"
	"time"

	nais_io_v1alpha1 "github.com/nais/liberator/pkg/apis/nais.io/v1alpha1"
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

func TestProblem(t *testing.T) {
	var app nais_io_v1alpha1.Application
	app.Status.SetError("error")
	app.Status.AddWarning(".spec.image", "warning")
	app.Status.AddDeprecation(".spec.image", "deprecation", time.Time{})

	problems := *app.Status.Problems

	assert.Equal(t, "error", problems[0].Message)
	assert.Equal(t, "warning", problems[1].Message)
	assert.Equal(t, "deprecation", problems[2].Message)

	app.Status.ClearProblems()

	assert.Nil(t, app.Status.Problems)
}
