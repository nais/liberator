package v1_test

import (
	"testing"

	data_nais_io_v1 "github.com/nais/liberator/pkg/apis/data.nais.io/v1"
	"github.com/stretchr/testify/assert"

	"github.com/nais/liberator/pkg/testutil"
)

var ignoredPostgresFields = []string{
	`.ObjectMeta.Annotations`,
	`.ObjectMeta.ClusterName`,
	`.ObjectMeta.CreationTimestamp`,
	`.ObjectMeta.CreationTimestamp.Time`,
	`.ObjectMeta.CreationTimestamp.Time.ext`,
	`.ObjectMeta.CreationTimestamp.Time.loc`,
	`.ObjectMeta.CreationTimestamp.Time.wall`,
	`.ObjectMeta.DeletionGracePeriodSeconds`,
	`.ObjectMeta.DeletionTimestamp`,
	`.ObjectMeta.Finalizers`,
	`.ObjectMeta.GenerateName`,
	`.ObjectMeta.Generation`,
	`.ObjectMeta.ManagedFields`,
	`.ObjectMeta.OwnerReferences`,
	`.ObjectMeta.ResourceVersion`,
	`.ObjectMeta.SelfLink`,
	`.ObjectMeta.UID`,
	`.Status`,
}

// Test that the example Postgres contains examples for all fields encountered.
// Examples MUST contain a non-zero value to be valid, so no empty strings, false booleans, or zero ints.
func TestExamplePostgresForDocumentation(t *testing.T) {
	postgres := data_nais_io_v1.ExamplePostgresForDocumentation()
	keys := testutil.ZeroFields(postgres)

	for _, key := range keys {
		if !testutil.StringSliceContains(ignoredPostgresFields, key) {
			assert.Fail(t, key, "`%s` does not exist with a non-zero value in data_nais_io_v1.ExamplePostgresForDocumentation", key)
		}
	}
}
