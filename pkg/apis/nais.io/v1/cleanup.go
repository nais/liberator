package nais_io_v1

type Cleanup struct {
	// Enables automatic cleanup
	// Default: `true`
	Enabled bool `json:"enabled"`
	// Strategy sets how a deployment might be handled.
	// Setting this to an empty list is equivalent to setting `enabled: false`.
	// Default: `["abort-rollout", "downscale"]`.
	//
	// - `abort-rollout`: if new pods in a deployment are failing, but previous pods from the previous working
	//    revision are still running, Babylon can roll the deployment back to the working revision,
	//    aborting the rollout.
	//
	// - `downscale`: if all pods in a deployment are failing, Babylon will set replicaset to 0
	Strategy []CleanupStrategy `json:"strategy,omitempty"`
	// +kubebuilder:validation:Pattern=`^[0-9]+h$`
	// Default: `24h`
	GracePeriod string `json:"gracePeriod,omitempty"`
}

// +kubebuilder:validation:Enum=abort-rollout;downscale
type CleanupStrategy string
