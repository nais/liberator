package nais_io_v1

import (
	"github.com/imdario/mergo"
)

// ApplyDefaults sets default values where they are missing from an Application spec.
func (alert *Alert) ApplyDefaults() error {
	return mergo.Merge(alert, getAlertDefaults())
}

func getAlertDefaults() *Alert {
	return &Alert{
		Spec:       AlertSpec{},
	}
}
