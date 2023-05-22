package nais_io_v1_test

import (
	"testing"

	nais_io_v1 "github.com/nais/liberator/pkg/apis/nais.io/v1"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	// Change this value to accept re-synchronization of ALL AzureAdApplication resources when deploying a new version.
	expectedHash = "fb3ded02e59f9458"
)

func TestAzureAdApplication_Hash(t *testing.T) {
	actual, err := minimalApplication().Hash()
	assert.NoError(t, err)
	assert.Equal(t, expectedHash, actual)
}

func minimalApplication() *nais_io_v1.AzureAdApplication {
	return &nais_io_v1.AzureAdApplication{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-app",
			Namespace: "test-namespace",
		},
		Spec: nais_io_v1.AzureAdApplicationSpec{
			SecretName: "test",
		},
		Status: nais_io_v1.AzureAdApplicationStatus{
			PasswordKeyIds:            []string{"test"},
			CertificateKeyIds:         []string{"test"},
			ClientId:                  "test",
			ObjectId:                  "test",
			ServicePrincipalId:        "test",
			SynchronizationHash:       expectedHash,
			SynchronizationSecretName: "test",
		},
	}
}
