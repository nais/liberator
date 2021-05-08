package kafka_nais_io_v1

import (
	"github.com/nais/liberator/pkg/hash"
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
	return hash.Hash(data)
}

func (in *AivenApplication) Hash() (string, error) {
	return hash.Hash(in.Spec)
}
