package nais_io_v1

type Cleanup struct {
	// Enables automatic cleanup
	// Default: `true`
	Enabled bool `json:"enabled"`
	// Rollback sets whether a deployment is rolled back or scaled down.
	// If `true` the deployment will be rolled back to the previous working version.
	// If `false` the deployment will be scaled down to zero replicas instead.
	// Default: `true`
	Rollback bool `json:"rollback"`
	// +kubebuilder:validation:Pattern=`^[0-9]+h$`
	// Default: `24h`
	GracePeriod string `json:"gracePeriod,omitempty"`
}
