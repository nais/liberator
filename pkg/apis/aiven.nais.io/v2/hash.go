package aiven_nais_io_v2

import "github.com/nais/liberator/pkg/hash"

func (in *AivenApplication) Hash() (string, error) {
	return hash.Hash(in.Spec)
}
