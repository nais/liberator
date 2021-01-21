package nais_io_v1

import (
	"encoding/json"
	"fmt"

	hash "github.com/mitchellh/hashstructure"
)

func (in *JwkerSpec) Hash() (string, error) {
	marshalled, err := json.Marshal(in)
	if err != nil {
		return "", err
	}
	h, err := hash.Hash(marshalled, nil)
	return fmt.Sprintf("%x", h), err
}
