package nais_io_v1alpha1_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/nais/liberator/pkg/apis/nais.io/v1alpha1"
	"github.com/nais/liberator/pkg/testutil"
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
	`.Spec.Azure.Application.Claims.Extra`, // TODO: Remove when these deprecated Azure fields are removed.
	`.Spec.Azure.Application.ReplyURLs`,
	`.Spec.Azure.Application.SinglePageApplication`,
	`.Spec.Strategy.RollingUpdate.MaxSurge.IntVal`,
	`.Spec.Strategy.RollingUpdate.MaxUnavailable`,
	`.Spec.Strategy.RollingUpdate.MaxUnavailable.IntVal`,
	`.Spec.Strategy.RollingUpdate.MaxUnavailable.StrVal`,
	`.Spec.Strategy.RollingUpdate.MaxUnavailable.Type`,
	`.Status`,
	`.Status.Conditions`,
	`.Status.CorrelationID`,
	`.Status.DeploymentRolloutStatus`,
	`.Status.Problems`,
	`.Status.RolloutCompleteTime`,
	`.Status.SynchronizationHash`,
	`.Status.SynchronizationState`,
	`.Status.SynchronizationTime`,
	`.Status.EffectiveImage`,
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
