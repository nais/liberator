package kafka_nais_io_v1

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTopicHash(t *testing.T) {
	spec := Topic{}
	hash, err := spec.Hash()
	assert.NoError(t, err)
	assert.Equal(t, "45bb0c5791695f91", hash)

	spec.Annotations = map[string]string{
		RemoveDataAnnotation: "true",
	}
	hash, err = spec.Hash()
	assert.NoError(t, err)
	assert.Equal(t, "89b3de1b2598c91c", hash)
}
