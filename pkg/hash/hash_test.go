package hash_test

import (
	"github.com/nais/liberator/pkg/hash"
	"github.com/stretchr/testify/assert"
	"testing"
)

type someStruct struct {
	SomeStructKey      string
	SomeOtherStructKey string
}

func TestHash(t *testing.T) {
	input := struct {
		SomeKey      string
		SomeOtherKey int
		SomeStruct   someStruct
	}{
		SomeKey:      "some-value",
		SomeOtherKey: 0,
		SomeStruct: someStruct{
			SomeStructKey:      "some-other-value",
			SomeOtherStructKey: "some-other-value",
		},
	}

	actual, err := hash.Hash(input)
	assert.NoError(t, err)
	assert.Equal(t, "8a26a8b71bf71ecb", actual)
}
