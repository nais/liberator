package nais_io_v1alpha1_test

import (
	"reflect"
	"sort"
	"testing"

	"github.com/nais/liberator/pkg/apis/nais.io/v1alpha1"
	"github.com/stretchr/testify/assert"
)

var ignoredFields = []string{
	`.ObjectMeta.Annotations`,
	`.ObjectMeta.ClusterName`,
	`.ObjectMeta.CreationTimestamp`,
	`.ObjectMeta.CreationTimestamp.Time`,
	`.ObjectMeta.CreationTimestamp.Time.ext`,
	`.ObjectMeta.CreationTimestamp.Time.wall`,
	`.ObjectMeta.Finalizers`,
	`.ObjectMeta.GenerateName`,
	`.ObjectMeta.Generation`,
	`.ObjectMeta.ManagedFields`,
	`.ObjectMeta.OwnerReferences`,
	`.ObjectMeta.ResourceVersion`,
	`.ObjectMeta.SelfLink`,
	`.ObjectMeta.UID`,
	`.Status`,
	`.Status.CorrelationID`,
	`.Status.DeploymentRolloutStatus`,
	`.Status.RolloutCompleteTime`,
	`.Status.SynchronizationHash`,
	`.Status.SynchronizationState`,
	`.Status.SynchronizationTime`,
}

// Test that the example Application contain examples for all fields encountered.
// Examples MUST contain a non-zero value to be valid, so no empty strings, false booleans, or zero ints.
func TestExampleApplicationForDocumentation(t *testing.T) {
	app := nais_io_v1alpha1.ExampleApplicationForDocumentation()
	seen := make(map[string]bool)

	observeMembers(seen, reflect.ValueOf(app), "")

	for _, k := range ignoredFields {
		delete(seen, k)
	}

	keys := make([]string, 0, len(seen))
	for k := range seen {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, key := range keys {
		assert.Truef(t, seen[key], "`%s` does not exist with a non-zero value in nais_io_v1alpha1.ExampleApplicationForDocumentation", key)
	}
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
