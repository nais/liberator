package hash

import (
	"encoding/json"
	"fmt"
	hash "github.com/mitchellh/hashstructure"
)

func Hash(input interface{}) (string, error) {
	marshalled, err := json.Marshal(input)
	if err != nil {
		return "", err
	}
	h, err := hash.Hash(marshalled, nil)
	return fmt.Sprintf("%x", h), err
}
