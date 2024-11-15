package hash

import (
	"encoding/json"
	"fmt"

	hash "github.com/mitchellh/hashstructure"
)

func Hash(input interface{}) (string, error) {
	h, err := IntHash(input)
	return fmt.Sprintf("%x", h), err
}

func IntHash(input interface{}) (uint64, error) {
	marshalled, err := json.Marshal(input)
	if err != nil {
		return 0, err
	}
	return hash.Hash(marshalled, nil)
}
