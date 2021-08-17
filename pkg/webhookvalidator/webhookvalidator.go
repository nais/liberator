package webhookvalidator

import (
	"fmt"
	"reflect"
	"strings"

	"k8s.io/apimachinery/pkg/util/validation/field"
)

func DeepComparison(new, old interface{}, path *field.Path) error {
	newValue := reflect.ValueOf(new)
	oldValue := reflect.ValueOf(old)
	if newValue.Kind() == reflect.Ptr {
		newValue = newValue.Elem()
	}
	if oldValue.Kind() == reflect.Ptr {
		oldValue = oldValue.Elem()
	}
	if newValue.Kind() != oldValue.Kind() {
		return fmt.Errorf("type mismatch")
	}

	cerr := compareObjects(newValue, oldValue, path)
	return cerr.ToAggregate()
}

func compareObjects(new, old reflect.Value, path *field.Path) (allErrs field.ErrorList) {
	if new.Kind() != reflect.Struct {
		return nil
	}

	newStruct := new.Type()
	// Iterate over all the fields of the current object being compared
	for i := 0; i < new.NumField(); i++ {
		newField := new.Field(i)

		if !newField.CanInterface() {
			// Ignore field if not exported
			continue
		}

		oldField := old.Field(i)

		if newField.Kind() == reflect.Ptr {
			// Derefence pointer if this current field is a pointer
			if newField.IsNil() || oldField.IsNil() {
				// TODO(thokra): Ensure that we allow the removal and addition of fields in nais.yaml spec
				continue
			}

			newField = newField.Elem()
			oldField = oldField.Elem()
		}

		tags := propertyMap(newStruct.Field(i))

		newPath := path.Child(jsonName(newStruct.Field(i)))
		if !tags["immutable"] && newField.Kind() == reflect.Struct {
			// If field is a struct, and it's not immutable, recursively descend and compare
			if err := compareObjects(newField, oldField, newPath); err != nil {
				allErrs = append(allErrs, err...)
				continue
			}
		}

		if tags["immutable"] && !reflect.DeepEqual(newField.Interface(), oldField.Interface()) {
			// If field is set to immutable, check if there's a change
			allErrs = append(allErrs, field.Invalid(newPath, newField.Interface(), "field is immutable"))
			continue
		}
	}
	return allErrs
}

// propertyMap creates a map["name of property"]boolean struct for easy look-up of `nais:"X"` tags
func propertyMap(field reflect.StructField) map[string]bool {
	tagss := strings.Split(strings.ToLower(field.Tag.Get("nais")), ",")
	tags := map[string]bool{}
	for _, t := range tagss {
		tags[t] = true
	}
	return tags
}

func jsonName(field reflect.StructField) string {
	t := field.Tag.Get("json")
	tag := field.Name
	parts := strings.Split(t, ",")
	if parts[0] != "" {
		tag = parts[0]
	}

	return tag
}
