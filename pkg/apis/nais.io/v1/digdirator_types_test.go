package nais_io_v1_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	nais_io_v1 "github.com/nais/liberator/pkg/apis/nais.io/v1"
)

func TestMaskinportenClient_CalculateHash(t *testing.T) {
	actual, err := minimalMaskinportenClient().Hash()
	assert.NoError(t, err)
	assert.Equal(t, "fcd4a1835320374a", actual)
}

func TestIDPortenClient_Hash(t *testing.T) {
	actual, err := minimalIDPortenClient().Hash()
	assert.NoError(t, err)
	assert.Equal(t, "de6ecbc3b6cb148b", actual)
}

func TestMaskinportenClientScopeSpec_Hash(t *testing.T) {
	actual, err := minimalMaskinportenExtendedClient().Hash()
	assert.NoError(t, err)
	assert.Equal(t, "fcd4a1835320374a", actual)
}

func minimalMaskinportenClient() *nais_io_v1.MaskinportenClient {
	return &nais_io_v1.MaskinportenClient{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-app",
			Namespace: "test-namespace",
		},
		Spec: nais_io_v1.MaskinportenClientSpec{
			Scopes: nais_io_v1.MaskinportenScope{
				ConsumedScopes: nil,
				ExposedScopes:  nil,
			},
		},
	}
}

func minimalIDPortenClient() *nais_io_v1.IDPortenClient {
	return &nais_io_v1.IDPortenClient{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-app",
			Namespace: "test-namespace",
		},
		Spec: nais_io_v1.IDPortenClientSpec{
			RedirectURIs: []nais_io_v1.IDPortenURI{
				"https://test.com",
			},
			SecretName: "test",
		},
	}
}

func minimalMaskinportenExtendedClient() *nais_io_v1.MaskinportenClient {
	return &nais_io_v1.MaskinportenClient{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-app",
			Namespace: "test-namespace",
		},
		Spec: nais_io_v1.MaskinportenClientSpec{
			Scopes: nais_io_v1.MaskinportenScope{
				ConsumedScopes: nil,
				ExposedScopes:  nil,
			},
		},
	}
}
