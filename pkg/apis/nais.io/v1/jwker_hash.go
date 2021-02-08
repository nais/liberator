package nais_io_v1

import (
	"github.com/nais/liberator/pkg/hash"
)

func (in *JwkerSpec) Hash() (string, error) {
	return hash.Hash(in)
}
