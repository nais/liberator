package controller

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type NaisConditionType string

const (
	// ReconcileSucceeded is set when the reconciliation has been successfully completed
	ReconcileSucceeded NaisConditionType = "Succeeded"
	// ReconcileFailedLocally is set when the reconciliation has failed due to a problem in the reconciler that is local to the cluster
	ReconcileFailedLocally NaisConditionType = "LocalFailure"

	// Add more detailed condition types here for other failure modes outside the reconciler/cluster

	// ReconcileFailedInAiven is set when the reconciliation has failed due to a problem in Aiven
	ReconcileFailedInAiven NaisConditionType = "AivenFailure"
	// ReconcileFailedInAzure is set when the reconciliation has failed due to a problem in Azure
	ReconcileFailedInAzure NaisConditionType = "AzureFailure"
)

type NaisCondition struct {
	// Type of condition.
	Type NaisConditionType `json:"type"`
	// Status of the condition, one of True, False, Unknown.
	Status corev1.ConditionStatus `json:"status"`
	// The last time this condition was updated.
	LastUpdateTime metav1.Time `json:"lastUpdateTime,omitempty"`
	// The reason for the condition's last transition.
	Reason string `json:"reason,omitempty"`
	// A human-readable message indicating details about the transition.
	Message string `json:"message,omitempty"`
}

type NaisStatus struct {
	// SynchronizationHash is the hash of the object most recently successfully synchronized
	SynchronizationHash string `json:"synchronizationHash,omitempty"`
	// SynchronizationState denotes whether the reconciliation has been successfully completed or not
	SynchronizationState string `json:"synchronizationState,omitempty"`
	// SynchronizationTime is the last time the Status subresource was updated
	SynchronizationTime *metav1.Time `json:"synchronizationTime,omitempty"`
	// SynchronizedGeneration is the generation most recently successfully synchronized
	SynchronizedGeneration int64 `json:"synchronizedGeneration,omitempty"`
	// ObservedGeneration is the generation most recently observed
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`
	// Represents the latest available observations of the resource current state.
	Conditions []NaisCondition `json:"conditions,omitempty"`
}

type ReconcileResult struct {
	DeleteFinalized bool
	Skipped         bool
	Requeue         bool
	Status          *NaisStatus
	Error           error
}
