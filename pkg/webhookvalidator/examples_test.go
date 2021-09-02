package webhookvalidator_test

import (
	"fmt"

	"github.com/nais/liberator/pkg/webhookvalidator"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func ExampleNaisCompare() {
	type Pet struct {
		// +nais:doc:Immutable=true
		Name string `nais:"key"`
		// +nais:doc:Immutable=true
		Species string `nais:"immutable"`
	}

	type Pets struct {
		Pets []Pet
	}

	new := Pets{
		Pets: []Pet{
			{Name: "Alberto", Species: "Dog"},
			{Name: "Sansa", Species: "Bird"},
		},
	}
	old := Pets{
		Pets: []Pet{
			{Name: "Alberto", Species: "Dog"},
			{Name: "Sansa", Species: "Cat"},
		},
	}

	err := webhookvalidator.NaisCompare(new, old, field.NewPath("animals"))
	fmt.Println(err)
	// Output: animals.Pets.1.Species: Invalid value: "Bird": field is immutable
}
