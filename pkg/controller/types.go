package controller

import (
	"fmt"

	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type SynchronizationState uint8

const (
	SynchronizationStateUnknown SynchronizationState = iota
	SynchronizationStateSuccessful
	SynchronizationStateFailed
)

func (s *SynchronizationState) String() string {
	switch *s {
	case SynchronizationStateUnknown:
		return "Unknown"
	case SynchronizationStateSuccessful:
		return "Successful"
	case SynchronizationStateFailed:
		return "Failed"
	}
	panic(fmt.Sprintf("Invalid SynchronizationState: %d", s))
}

func (s *SynchronizationState) UnmarshalText(text []byte) error {
	str := string(text)
	switch str {
	case "Successful":
		*s = SynchronizationStateSuccessful
	case "Failed":
		*s = SynchronizationStateFailed
	case "Unknown":
		*s = SynchronizationStateUnknown
	default:
		return fmt.Errorf("invalid SynchronizationState: %v", str)
	}
	return nil
}

func (s *SynchronizationState) MarshalText() (text []byte, err error) {
	return []byte(s.String()), nil
}

// +kubebuilder:object:generate=true

type NaisStatus struct {
	// SynchronizationHash is the hash of the object most recently successfully synchronized
	SynchronizationHash string `json:"synchronizationHash,omitempty"`
	// SynchronizationTime is the last time the Status subresource was updated
	SynchronizationTime *metav1.Time `json:"synchronizationTime,omitempty"`
	// SynchronizationState is the result of the most recent synchronization
	SynchronizationState SynchronizationState `json:"synchronizationState,omitempty"`
	// SynchronizedGeneration is the generation most recently successfully synchronized
	SynchronizedGeneration int64 `json:"synchronizedGeneration,omitempty"`
	// ObservedGeneration is the generation most recently observed
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`
	// Represents the latest available observations of the resource current state.
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

type ReconcileResult struct {
	// Requeue indicates whether the reconciliation should be requeued for a new attempt in a few seconds
	Requeue bool

	// Conditions should be a list of possible failure conditions
	//
	// Status == ConditionTrue means a failure has happened
	Conditions []metav1.Condition

	// State should indicate the overall state of the reconciliation
	//
	// If State == SynchronizationStateFailed, Conditions should contain at least one Condition with Status == ConditionTrue
	State SynchronizationState
}

type NaisResource interface {
	client.Object

	// ApplyDefaults must apply default values to the resource
	ApplyDefaults() error

	// Hash must return a hash of the relevant parts of the resource
	Hash() (string, error)

	// GetStatus must return a non-nil pointer to the resource status attached to this resource
	GetStatus() *NaisStatus
}

type Status interface {
	UpdateSynchronizationData(hash string, state SynchronizationState, generation int64)
	UpdateObservationData(generation int64)
	UpdateConditions(conditions []metav1.Condition, generation int64)
}

func (n *NaisStatus) UpdateSynchronizationData(hash string, state SynchronizationState, generation int64) {
	n.SynchronizationHash = hash
	n.SynchronizationState = state
	n.SynchronizedGeneration = generation
	now := metav1.Now()
	n.SynchronizationTime = &now
}

func (n *NaisStatus) UpdateObservationData(generation int64) {
	n.ObservedGeneration = generation
}

func (n *NaisStatus) UpdateConditions(conditions []metav1.Condition, generation int64) {
	for _, newCondition := range conditions {
		newCondition.ObservedGeneration = generation
		meta.SetStatusCondition(&n.Conditions, newCondition)
	}
}
