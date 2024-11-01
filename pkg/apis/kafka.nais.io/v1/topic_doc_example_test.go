package kafka_nais_io_v1_test

import (
	"testing"

	kafka_nais_io_v1 "github.com/nais/liberator/pkg/apis/kafka.nais.io/v1"
	"github.com/stretchr/testify/assert"

	"github.com/nais/liberator/pkg/testutil"
)

var ignoredTopicFields = []string{
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
	`.Spec.Config.MinCleanableDirtyRatioPercent`,
	`.Status`,
	`.Status.SynchronizationState`,
	`.Status.SynchronizationHash`,
	`.Status.SynchronizationTime`,
	`.Status.CredentialsExpiryTime`,
	`.Status.Errors`,
	`.Status.Message`,
	`.Status.FullyQualifiedName`,
	`.Status.LatestAivenSyncFailure`,
}

// Test that the example Application contains examples for all fields encountered.
// Examples MUST contain a non-zero value to be valid, so no empty strings, false booleans, or zero ints.
func TestExampleTopicForDocumentation(t *testing.T) {
	topic := kafka_nais_io_v1.ExampleTopicForDocumentation()
	keys := testutil.ZeroFields(topic)

	for _, key := range keys {
		if !testutil.StringSliceContains(ignoredTopicFields, key) {
			assert.Fail(t, key, "`%s` does not exist with a non-zero value in kafka_nais_io_v1.ExampleTopicForDocumentation", key)
		}
	}
}
