package nais_io_v1alpha1

import (
	"testing"
)

func TestSnorlax_Enabled(t *testing.T) {
	tests := map[string]struct {
		s    *Snorlax
		want bool
	}{
		"default nil":   {want: false},
		"default empty": {s: &Snorlax{}, want: false},
		"disabled":      {s: &Snorlax{Strategy: SnorlaxDisabled}, want: false},
		"enabled":       {s: &Snorlax{Strategy: SnorlaxEnabled}, want: true},
		"unknown":       {s: &Snorlax{Strategy: "unknown"}, want: false},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			if got := tt.s.Enabled(); got != tt.want {
				t.Errorf("Snorlax.Enabled() = %v, want %v", got, tt.want)
			}
		})
	}
}
