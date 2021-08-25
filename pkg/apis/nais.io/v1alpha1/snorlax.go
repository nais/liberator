package nais_io_v1alpha1

type SnorlaxStrategy string

const (
	SnorlaxEnabled  SnorlaxStrategy = "enabled"
	SnorlaxDisabled SnorlaxStrategy = "disabled"
)

type Snorlax struct {
	// Strategy blah blah. Defaults to disabled.
	// +kubebuilder:validation:Enum=enabled;disabled;""
	Strategy SnorlaxStrategy `json:"strategy,omitempty"`

	// If specified, the app will not be suspended within the time specified.
	// This will not automatically start the application.
	// +kubebuilder:validation:Pattern=`((((\d+,)+\d+|(\d+(\/|-)\d+)|\d+|\*) ?){5,7})?`
	NoSuspendCron string `json:"noSuspendCron,omitempty"`
}

func (s *Snorlax) Enabled() bool {
	if s == nil {
		return false
	}
	return s.Strategy == SnorlaxEnabled
}
