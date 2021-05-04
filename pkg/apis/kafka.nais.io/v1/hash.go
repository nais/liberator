package kafka_nais_io_v1

import (
	"encoding/json"
	"fmt"

	"github.com/mitchellh/hashstructure"
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
	return hash(data)
}

func (in *AivenApplication) Hash() (string, error) {
	return hash(in.Spec)
}

func hash(data interface{}) (string, error) {
	marshalled, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	h, err := hashstructure.Hash(marshalled, nil)
	return fmt.Sprintf("%x", h), err
}
