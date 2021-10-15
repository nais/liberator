package nais_io_v1alpha1_test

import (
	"testing"

	"github.com/nais/liberator/pkg/apis/nais.io/v1alpha1"
	"github.com/nais/liberator/pkg/testutil"
	"github.com/stretchr/testify/assert"
)

var ignoredApplicationFields = []string{
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
	`.Status.Conditions`,
	`.Status.CorrelationID`,
	`.Status.DeploymentRolloutStatus`,
	`.Status.RolloutCompleteTime`,
	`.Status.SynchronizationHash`,
	`.Status.SynchronizationState`,
	`.Status.SynchronizationTime`,
}

// Test that the example Application contains examples for all fields encountered.
// Examples MUST contain a non-zero value to be valid, so no empty strings, false booleans, or zero ints.
func TestExampleApplicationForDocumentation(t *testing.T) {
	app := nais_io_v1alpha1.ExampleApplicationForDocumentation()
	keys := testutil.ZeroFields(app)

	for _, key := range keys {
		if !testutil.StringSliceContains(ignoredApplicationFields, key) {
			assert.Fail(t, key, "`%s` does not exist with a non-zero value in nais_io_v1alpha1.ExampleApplicationForDocumentation", key)
		}
	}
}
