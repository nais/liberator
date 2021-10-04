package nais_io_v1

import (
	"time"

	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Status contains different NAIS status properties
type Status struct {
	SynchronizationTime     int64              `json:"synchronizationTime,omitempty"`
	RolloutCompleteTime     int64              `json:"rolloutCompleteTime,omitempty"`
	CorrelationID           string             `json:"correlationID,omitempty"`
	DeploymentRolloutStatus string             `json:"deploymentRolloutStatus,omitempty"`
	SynchronizationState    string             `json:"synchronizationState,omitempty"`
	SynchronizationHash     string             `json:"synchronizationHash,omitempty"`
	Conditions              *[]metav1.Condition `json:"conditions,omitempty"`
}

func (in *Status ) SetReadyCondition(condition metav1.ConditionStatus, reason string, message string) {
	if in.Conditions == nil {
		in.Conditions = &[]metav1.Condition{}
	}
	readyCondition := meta.FindStatusCondition(*in.Conditions, "READY")
	if readyCondition == nil {
		readyCondition = &metav1.Condition{
			Type:               "Ready",
			Status:             condition,
			LastTransitionTime: metav1.NewTime(time.Now()),
			ObservedGeneration: 1,
			Reason:             reason,
			Message:            message,
		}
	}else {
		readyCondition.Status = condition
		readyCondition.ObservedGeneration += 1
		readyCondition.Reason = reason
		readyCondition.Message = message
	}
	meta.SetStatusCondition(in.Conditions, *readyCondition)
}
