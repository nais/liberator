package events

// Machine readable status strings describing the current synchronization state of a custom resource.
const (
	FailedPrepare         = "FailedPrepare"
	FailedStatusUpdate    = "FailedStatusUpdate"
	FailedSynchronization = "FailedSynchronization"
	Retrying              = "Retrying"
	RolloutComplete       = "RolloutComplete"
	Synchronized          = "Synchronized"
)
