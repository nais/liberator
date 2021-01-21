package kafka_nais_io_v1_test

import (
	"testing"

	"github.com/nais/liberator/pkg/apis/kafka.nais.io/v1"
	"github.com/stretchr/testify/assert"
)

func TestHash(t *testing.T) {
	spec := kafka_nais_io_v1.Topic{}
	hash, err := spec.Hash()
	assert.NoError(t, err)
	assert.Equal(t, "45bb0c5791695f91", hash)

	spec.Annotations = map[string]string{
		kafka_nais_io_v1.RemoveDataAnnotation: "true",
	}
	hash, err = spec.Hash()
	assert.NoError(t, err)
	assert.Equal(t, "89b3de1b2598c91c", hash)
}
