package testutil

import (
	"reflect"
	"sort"
)

// Test that an example resource contain examples for all fields encountered.
// Examples MUST contain a non-zero value to be valid, so no empty strings, false booleans, or zero ints.

// Return a list of JSON paths that have a default value, e.g. nil, zero-length slices,
// empty strings, false booleans, or zero ints.
// The function looks deeply within nested structures.
func ZeroFields(input interface{}) []string {
	seen := make(map[string]bool)

	observeMembers(seen, reflect.ValueOf(input), "")

	keys := make([]string, 0, len(seen))
	for k := range seen {
		if !seen[k] {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)

	return keys
}

func StringSliceContains(slice []string, key string) bool {
	for _, candidate := range slice {
		if key == candidate {
			return true
		}
	}
	return false
}

// Extra check for non-nil slices with zero elements.
func isZero(v reflect.Value) bool {
	zero := v.IsZero()
	if !zero && v.Kind() == reflect.Slice {
		return v.Len() == 0
	}
	return zero
}

// Recurse through a structure and record all members seen.
// Sets seen[key] to true if a non-zero value is encountered.
func observeMembers(seen map[string]bool, v reflect.Value, path string) {
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.Invalid:
		return

	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			f := v.Field(i)
			childPath := path + "." + v.Type().Field(i).Name
			observeMembers(seen, f, childPath)
		}
	case reflect.Map:
		for _, k := range v.MapKeys() {
			childPath := path + "." + k.String()
			observeMembers(seen, v.MapIndex(k), childPath)
		}
	case reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			f := v.Index(i)
			observeMembers(seen, f, path)
		}
	}

	seen[path] = seen[path] || !isZero(v)
}
