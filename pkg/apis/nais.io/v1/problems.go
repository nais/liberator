package nais_io_v1

import (
	"time"
)

// Error: tilstand hvor Naiserator ikke kommer seg videre med deploy.
//   - feil i spec
//
// Warning: tilstand hvor flere deler av spec ikke henger riktig sammen med hverandre.
//   - du har konfigurert A, som krever enabled på B
//   - prøver å deploye til et domene som ikke finnes i dette clusteret
//   - har man konfigurert mer enn ett felt i en spec som kun støtter en av gangen?
//
// Deprecation: tilstand hvor vi planlegger å endre eller fjerne en feature.
//   - feature A utilgjengelig fra tidspunkt B
//   - bruker deploy API keys framfor oauth2
type ProblemKind string

const (
	ProblemKindDeprecation ProblemKind = "Deprecation"
	ProblemKindWarning     ProblemKind = "Warning"
	ProblemKindError       ProblemKind = "Error"
)

// Problems deal with errors, warnings and deprecations caused by invalid usage of the Application and NaisJob specs.
// They are user-facing and will be shown in various frontends, such as `kubectl describe app` and Nais console.
type Problem struct {
	// Full name of spec field that triggered the error, e.g. `.spec.image`.
	Source *string `json:"source,omitempty"`
	// Severity or kind of problem.
	Type ProblemKind `json:"type"`
	// If the problem is related to deprecation of some system, this field
	// MAY contain the end-of-life date for that particular system, formatted
	// as a ISO8601 date.
	EndOfLife *string `json:"endOfLife,omitempty"`
	// Human-readable message describing the problem.
	// The message will be visible in Nais console.
	Message string `json:"message"`
}

func (in *Status) ClearProblems() {
	in.Problems = nil
}

// Because Go.
func (in *Status) ensureHasProblemsSlice() {
	if in.Problems != nil {
		return
	}
	problems := make([]Problem, 0)
	in.Problems = &problems
}

// Use SetError for fail-fast permanent errors.
func (in *Status) SetError(message string) {
	in.ensureHasProblemsSlice()
	*in.Problems = append(*in.Problems, Problem{
		Type:    ProblemKindError,
		Message: message,
	})
}

// Use AddWarning to communicate something that might be wrongly configured,
// such as using spec fields that will not be used due to not being enabled.
// Another case might be that we have an external deprecation without a due date.
func (in *Status) AddWarning(specField string, message string) {
	in.ensureHasProblemsSlice()
	*in.Problems = append(*in.Problems, Problem{
		Type:    ProblemKindWarning,
		Source:  &specField,
		Message: message,
	})
}

// Use AddDeprecation for features that will be changed or removed at a well-defined in the future.
func (in *Status) AddDeprecation(specField string, message string, endOfLife time.Time) {
	in.ensureHasProblemsSlice()
	endOfLifeDate := endOfLife.Format(time.DateOnly)
	*in.Problems = append(*in.Problems, Problem{
		Type:      ProblemKindDeprecation, // or warning?
		Source:    &specField,
		EndOfLife: &endOfLifeDate,
		Message:   message,
	})
}
