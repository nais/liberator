package nais_io_v1alpha1

import (
	"testing"

	nais_io_v1 "github.com/nais/liberator/pkg/apis/nais.io/v1"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func TestApplicationValidator_OpenSearchUses(t *testing.T) {
	openSearch := &unstructured.Unstructured{}
	openSearch.SetAPIVersion("nais.io/v1")
	openSearch.SetKind("OpenSearch")
	openSearch.SetName("search")
	openSearch.SetNamespace("test-ns")

	app := &Application{
		ObjectMeta: metav1.ObjectMeta{Name: "test-app", Namespace: "test-ns"},
		Spec:       ApplicationSpec{Uses: &nais_io_v1.Uses{OpenSearch: []nais_io_v1.OpenSearchUse{{Name: "search", Permission: "read"}}}},
	}

	t.Run("accepts existing namespace-local OpenSearch", func(t *testing.T) {
		validator := &ApplicationValidator{Client: fakeKubeClient(openSearch)}
		_, err := validator.ValidateCreate(t.Context(), app)
		require.NoError(t, err)
	})

	t.Run("rejects missing OpenSearch", func(t *testing.T) {
		validator := &ApplicationValidator{Client: fakeKubeClient()}
		_, err := validator.ValidateCreate(t.Context(), app)
		require.ErrorContains(t, err, "OpenSearch 'search' does not exist")
	})

	t.Run("rejects OpenSearch in another namespace", func(t *testing.T) {
		other := openSearch.DeepCopy()
		other.SetNamespace("other-ns")
		validator := &ApplicationValidator{Client: fakeKubeClient(other)}
		_, err := validator.ValidateCreate(t.Context(), app)
		require.ErrorContains(t, err, "OpenSearch 'search' does not exist")
	})
}
