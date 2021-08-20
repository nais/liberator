package webhookvalidator

import (
	"testing"

	"k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func TestDeepComparison(t *testing.T) {
	tests := map[string]struct {
		New        interface{}
		Old        interface{}
		TestErrors []string
	}{
		"Ensure equal leaf-like object pass": {
			New: SmallStruct{A: 2},
			Old: SmallStruct{A: 2},
		},
		"Ensure equal leaf-like object pass with modified non-immutable field": {
			New: SmallStruct{A: 2, B: 0},
			Old: SmallStruct{A: 2, B: 1337},
		},
		"Immutable leaf-like object fail if not equal": {
			New: SmallStruct{A: 1},
			Old: SmallStruct{A: 2},
			TestErrors: []string{
				"test.A",
			},
		},
		"Immutable nested object fails if immutable field not equal": {
			New: mediumStruct{
				Sub:          SmallStruct{A: 0},
				SubImmutable: SmallStruct{A: 1},
			},
			Old: mediumStruct{
				Sub:          SmallStruct{A: 0},
				SubImmutable: SmallStruct{A: 2},
			},
			TestErrors: []string{
				"test.SubImmutable",
			},
		},
		"Ensure nested equal objects pass": {
			New: mediumStruct{
				Sub:          SmallStruct{A: 0},
				SubImmutable: SmallStruct{A: 1},
			},
			Old: mediumStruct{
				Sub:          SmallStruct{A: 0},
				SubImmutable: SmallStruct{A: 1},
			},
		},
		"Nested object fails if immutable leaf field not equal": {
			New: mediumStruct{
				Sub:          SmallStruct{A: 0},
				SubImmutable: SmallStruct{A: 1},
			},
			Old: mediumStruct{
				Sub:          SmallStruct{A: 2},
				SubImmutable: SmallStruct{A: 1},
			},
			TestErrors: []string{
				"test.Sub.A",
			},
		},
		"Immutable pointer field pass if equal": {
			New: ptrStructImmutable{
				Pint: ptrInt(1),
			},
			Old: ptrStructImmutable{
				Pint: ptrInt(1),
			},
		},
		"Immutable pointer field fail if not equal": {
			New: ptrStructImmutable{
				Pint: ptrInt(1),
			},
			Old: ptrStructImmutable{
				Pint: ptrInt(3),
			},
			TestErrors: []string{
				"test.Pint",
			},
		},
		"Immutable struct pointer field fail if not equal": {
			New: ptrStructImmutable{
				Pint: ptrInt(1),
				PStruct: &SmallStruct{
					A: 5,
				},
			},
			Old: ptrStructImmutable{
				Pint: ptrInt(1),
				PStruct: &SmallStruct{
					A: 15,
				},
			},
			TestErrors: []string{
				"test.PStruct",
			},
		},
		"Immutable struct pointer field pass if equal": {
			New: ptrStructImmutable{
				Pint: ptrInt(1),
				PStruct: &SmallStruct{
					A: 5,
				},
			},
			Old: ptrStructImmutable{
				Pint: ptrInt(1),
				PStruct: &SmallStruct{
					A: 5,
				},
			},
		},
		"Inline struct fails if not equal": {
			New: inlineStruct{
				SmallStruct{
					A: 5,
				},
			},
			Old: inlineStruct{
				SmallStruct{
					A: 2,
				},
			},
			TestErrors: []string{
				"test.SmallStruct.A",
			},
		},
		"Inline struct pass if equal": {
			New: inlineStruct{
				SmallStruct{
					A: 5,
				},
			},
			Old: inlineStruct{
				SmallStruct{
					A: 5,
				},
			},
		},
		"Immutable slice and map pass if equal": {
			New: SmallStruct{
				Slice: []int{1, 2, 3},
				Map:   map[string]int{"a": 1, "b": 2, "c": 3},
			},
			Old: SmallStruct{
				Slice: []int{1, 2, 3},
				Map:   map[string]int{"a": 1, "b": 2, "c": 3},
			},
		},
		"Immutable map fail if not equal": {
			New: SmallStruct{
				Map: map[string]int{"a": 1, "b": 2, "c": 3},
			},
			Old: SmallStruct{
				Map: map[string]int{"a": 3, "b": 2, "c": 1},
			},
			TestErrors: []string{
				"test.Map",
			},
		},
		"Immutable slice fail if not equal": {
			New: SmallStruct{
				Slice: []int{1, 2, 3},
			},
			Old: SmallStruct{
				Slice: []int{3, 2, 1},
			},
			TestErrors: []string{
				"test.Slice",
			},
		},
		"Slice of struct with immutable fields fail if not equal on key": {
			New: mediumStruct{
				SliceStruct: []SmallStruct{
					{
						A: 2,
						C: 5,
					},
				},
			},
			Old: mediumStruct{
				SliceStruct: []SmallStruct{
					{
						A: 2,
						C: 8,
					},
				},
			},
			TestErrors: []string{
				"test.SliceStruct.0.C",
			},
		},
		"Slice of struct with immutable fields pass if equal on key": {
			New: mediumStruct{
				SliceStruct: []SmallStruct{
					{
						A: 2,
						C: 9,
					},
				},
			},
			Old: mediumStruct{
				SliceStruct: []SmallStruct{
					{
						A: 2,
						C: 9,
					},
				},
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			err := NaisCompare(tt.New, tt.Old, field.NewPath("test"))
			if len(tt.TestErrors) == 0 {
				if err != nil {
					t.Fatalf("expected error: %v", err)
				}
				return
			}

			if err == nil {
				t.Fatal("expected error, but got nil")
			}

			errors := err.(errors.Aggregate).Errors()
			if len(tt.TestErrors) != len(errors) {
				t.Errorf("expected %v errors, got %v", len(tt.TestErrors), len(errors))
			}

			found := map[string]bool{}
			for _, terr := range errors {
				found[terr.(*field.Error).Field] = true
			}
			for _, expected := range tt.TestErrors {
				if !found[expected] {
					t.Errorf("expected error: %q", expected)
				}
				delete(found, expected)
			}
			for val := range found {
				t.Errorf("got %q, but did not expect it", val)
			}
		})
	}
}

type SmallStruct struct {
	A     int `nais:"immutable,key"`
	B     int
	C     int            `nais:"immutable"`
	Slice []int          `nais:"immutable"`
	Map   map[string]int `nais:"immutable"`
}

type mediumStruct struct {
	Sub          SmallStruct
	SubImmutable SmallStruct `nais:"immutable"`
	SliceStruct  []SmallStruct
}

type inlineStruct struct {
	SmallStruct
}

type ptrStructImmutable struct {
	Pint    *int         `nais:"immutable"`
	PStruct *SmallStruct `nais:"immutable"`
}

type ptrStruct struct {
	Pint    *int
	PStruct *SmallStruct
}

func ptrInt(i int) *int {
	return &i
}
