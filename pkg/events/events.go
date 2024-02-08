package events

// Machine readable status strings describing the current synchronization state of a NAIS application or job.
// These may be used both in the `.status.synchronizationState` and `.status.conditions[].reason` fields.

const (
	// Unable to read preliminary data from cluster. This error is always transient.
	FailedPrepare = "FailedPrepare"

	// Indicates invalid configuration in nais.yaml.
	FailedGenerate = "FailedGenerate"

	// Naiserator failed to persist resources into the cluster, and the error is permanent.
	FailedSynchronization = "FailedSynchronization"

	// Naiserator failed to persist resources into the cluster, and the error is transient.
	Retrying = "Retrying"

	// Naiserator has persisted resources into the cluster, and the rollout is in progress.
	Synchronized = "Synchronized"

	// Indicates that everything is complete - the most recently deployed version of the application has been activated.
	RolloutComplete = "RolloutComplete"

	// Emitted by events when the status field cannot be updated.
	FailedStatusUpdate = "FailedStatusUpdate"
)
