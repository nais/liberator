package nais_io_v1

import (
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Status contains different Nais status properties
type Status struct {
	SynchronizationTime     int64               `json:"synchronizationTime,omitempty"`
	RolloutCompleteTime     int64               `json:"rolloutCompleteTime,omitempty"`
	CorrelationID           string              `json:"correlationID,omitempty"`
	DeploymentRolloutStatus string              `json:"deploymentRolloutStatus,omitempty"`
	SynchronizationState    string              `json:"synchronizationState,omitempty"`
	SynchronizationHash     string              `json:"synchronizationHash,omitempty"`
	EffectiveImage          string              `json:"effectiveImage,omitempty"`
	Problems                *[]Problem          `json:"problems,omitempty"`
	Conditions              *[]metav1.Condition `json:"conditions,omitempty"`
}

// SetSynchronizationStateWithCondition
// is a shorthand function for setting synchronization state.
// Additionally, the state and the human-readable message is stored in a condition.
func (in *Status) SetSynchronizationStateWithCondition(reason, message string) {
	const conditionType = "SynchronizationState"
	in.SynchronizationState = reason
	in.SetCondition(conditionType, metav1.ConditionTrue, reason, message)
}

// SetCondition is a wrapper around an upstream API that does more or less the same thing.
// The condition with the matching `typ` is either created or updated in .status.conditions[].
func (in *Status) SetCondition(typ string, status metav1.ConditionStatus, reason, message string) {
	if in.Conditions == nil {
		in.Conditions = &[]metav1.Condition{}
	}

	condition := &metav1.Condition{
		Type:    typ,
		Status:  status,
		Reason:  reason,
		Message: message,
	}

	meta.SetStatusCondition(in.Conditions, *condition)
}
