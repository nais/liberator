package kafka_nais_io_v1_test

import (
	"testing"

	kafka_nais_io_v1 "github.com/nais/liberator/pkg/apis/kafka.nais.io/v1"
	"github.com/nais/liberator/pkg/testutil"
	"github.com/stretchr/testify/assert"
)

var ignoredStreamFields = []string{
	`.ObjectMeta.Annotations`,
	`.ObjectMeta.CreationTimestamp`,
	`.ObjectMeta.CreationTimestamp.Time`,
	`.ObjectMeta.DeletionGracePeriodSeconds`,
	`.ObjectMeta.DeletionTimestamp`,
	`.ObjectMeta.Finalizers`,
	`.ObjectMeta.GenerateName`,
	`.ObjectMeta.Generation`,
	`.ObjectMeta.Labels`,
	`.ObjectMeta.ManagedFields`,
	`.ObjectMeta.OwnerReferences`,
	`.ObjectMeta.ResourceVersion`,
	`.ObjectMeta.SelfLink`,
	`.ObjectMeta.UID`,
	`.Status`,
	`.TypeMeta`,
	`.TypeMeta.APIVersion`,
	`.TypeMeta.Kind`,
}

// Test that the example Stream contains examples for all fields encountered.
func TestStreamDocExample(t *testing.T) {
	stream := kafka_nais_io_v1.StreamDocExample()
	keys := testutil.ZeroFields(stream)

	for _, key := range keys {
		if !testutil.StringSliceContains(ignoredStreamFields, key) {
			assert.Fail(t, key, "`%s` does not exist with a non-zero value in kafka_nais_io_v1.StreamDocExample", key)
		}
	}
}
