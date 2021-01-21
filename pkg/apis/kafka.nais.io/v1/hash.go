package kafka_nais_io_v1

import (
	"encoding/json"
	"fmt"

	hash "github.com/mitchellh/hashstructure"
)

func (in *Topic) Hash() (string, error) {
	type hashFields struct {
		Annotations map[string]string
		Spec        TopicSpec
	}
	annotations := in.Annotations
	if annotations == nil {
		annotations = map[string]string{}
	}
	data := hashFields{
		Annotations: map[string]string{
			RemoveDataAnnotation: annotations[RemoveDataAnnotation],
		},
		Spec: in.Spec,
	}
	marshalled, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	h, err := hash.Hash(marshalled, nil)
	return fmt.Sprintf("%x", h), err
}
