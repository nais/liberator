package kafka_nais_io_v1

import (
	"github.com/stretchr/testify/assert"
	"testing"
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

func Test_hash(t *testing.T) {
	tests := []struct {
		name    string
		data    interface{}
		want    string
		wantErr bool
	}{
		{name: "ClearAivenApplication", data: AivenApplication{}.Spec, want: "ae9606dc1d9142ad", wantErr: false},
		{name: "AivenApplicationWithSecretName", data: AivenApplication{
			Spec: AivenApplicationSpec{
				SecretName: "this-is-my-secret",
			},
		}.Spec, want: "a26742b533308093", wantErr: false},
		{name: "Empty spec", data: nil, want: "fd836cbf52075630", wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := hash(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("hash() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("hash() got = %v, want %v", got, tt.want)
			}
		})
	}
}
