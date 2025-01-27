package nais_io_v1_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/nais/liberator/pkg/apis/nais.io/v1"
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
	`.Status`,
	`.Status.Conditions`,
	`.Status.CorrelationID`,
	`.Status.DeploymentRolloutStatus`,
	`.Status.Problems`,
	`.Status.RolloutCompleteTime`,
	`.Status.SynchronizationHash`,
	`.Status.SynchronizationState`,
	`.Status.SynchronizationTime`,
	`.Spec.AccessPolicy.Outbound.External.IPv4`,
	`.Spec.Azure.Application.Claims.Extra`, // TODO: Remove when these deprecated Azure fields are removed.
	`.Spec.Azure.Application.ReplyURLs`,
	`.Spec.Azure.Application.SinglePageApplication`,
	`.Spec.Azure.Application.Tenant`,
	`.Spec.Redis`,  // TODO: Add to example when feature is ready
	`.Spec.Valkey`, // TODO: Add to example when feature is ready
}

// Test that the example NaisJob contains examples for all fields encountered.
// Examples MUST contain a non-zero value to be valid, so no empty strings, false booleans, or zero ints.
func TestExampleNaisjobForDocumentation(t *testing.T) {
	naisjob := nais_io_v1.ExampleNaisjobForDocumentation()
	keys := testutil.ZeroFields(naisjob)

	for _, key := range keys {
		if !testutil.StringSliceContains(ignoredApplicationFields, key) {
			assert.Fail(t, key, "`%s` does not exist with a non-zero value in nais_io_v1alpha1.ExampleNaisjobForDocumentation", key)
		}
	}
}
