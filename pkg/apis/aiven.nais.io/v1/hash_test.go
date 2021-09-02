package aiven_nais_io_v1

import "testing"

func Test_hash(t *testing.T) {
	tests := []struct {
		name     string
		aivenapp *AivenApplication
		want     string
		wantErr  bool
	}{
		{name: "ClearAivenApplication", aivenapp: &AivenApplication{}, want: "b2b67239ec384acc", wantErr: false},
		{name: "AivenApplicationWithSecretName", aivenapp: &AivenApplication{
			Spec: AivenApplicationSpec{
				SecretName: "this-is-my-secret",
			},
		}, want: "743f9f56c913003d", wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.aivenapp.Hash()
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
