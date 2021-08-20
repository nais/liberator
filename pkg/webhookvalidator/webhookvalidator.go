package webhookvalidator

import (
	"fmt"
	"reflect"
	"strings"

	"k8s.io/apimachinery/pkg/util/validation/field"
)

func NaisCompare(new, old interface{}, path *field.Path) error {
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

	return compareObjects(newValue, oldValue, path).ToAggregate()
}

func compareObjects(new, old reflect.Value, path *field.Path) (allErrs field.ErrorList) {
	if new.Kind() == reflect.Slice || new.Kind() == reflect.Array {
		// Must be handled as a special case when comparing objects in slices
		return sliceMutationError(new, old, path)
	}
	if new.Kind() != reflect.Struct {
		// Only structs (iow: their fields) can have the "immutable" tag set
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
		if !tags["immutable"] {
			// Recursively descend into the fields of the current struct
			if err := compareObjects(newField, oldField, newPath); err != nil {
				allErrs = append(allErrs, err...)
				continue
			}
		}

		if tags["immutable"] && valuesDiffer(newField, oldField) {
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

// jsonName returns the name of the field by default, override if `json:"name"` is set
func jsonName(field reflect.StructField) string {
	t := field.Tag.Get("json")
	tag := field.Name
	parts := strings.Split(t, ",")
	if parts[0] != "" {
		tag = parts[0]
	}

	return tag
}

func valuesDiffer(new, old reflect.Value) bool {
	return !reflect.DeepEqual(new.Interface(), old.Interface())
}

// sliceMutationError compares two slices, and returns errors if they differ
//	Does not support slices of pointers (pointers of any kind), potential future todo.
func sliceMutationError(new, old reflect.Value, path *field.Path) (allErrs field.ErrorList) {
	if new.Len() == 0 || old.Len() == 0 {
		// We don't want to compare unless both slices have elements within them
		return allErrs
	}
	if new.Type().Elem().Kind() != reflect.Struct {
		// Immutability tag can only be used on fields within a struct
		return allErrs
	}

	keys := keysToCheck(new.Index(0).Type())
	if len(keys) == 0 && doPanic {
		panic("No `nais:\"key\"` defined for type on field" + path.String())
	}

OUTER:
	for i := 0; i < new.Len(); i++ {
		// For each element in "new" slice
		newVal := new.Index(i)

		for j := 0; j < old.Len(); j++ {
			// compared with each element in "old" slice
			oldVal := old.Index(j)
			for _, key := range keys {
				// Use the comparison keys found with keysToCheck to compare the elements from the two slices
				if !valuesDiffer(newVal.FieldByName(key), oldVal.FieldByName(key)) {
					errs := compareObjects(newVal, oldVal, path.Child(fmt.Sprint(i)))
					if errs != nil {
						allErrs = append(allErrs, errs...)
					}
					continue OUTER
				}
			}
		}
	}

	return allErrs
}

func keysToCheck(typ reflect.Type) (ret []string) {
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		properties := propertyMap(field)
		if properties["key"] {
			ret = append(ret, field.Name)
		}
	}

	return ret
}
