package nais_io_v1

import (
	"testing"

	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func TestJobValidator_OpenSearchUses(t *testing.T) {
	openSearch := &unstructured.Unstructured{}
	openSearch.SetAPIVersion("nais.io/v1")
	openSearch.SetKind("OpenSearch")
	openSearch.SetName("search")
	openSearch.SetNamespace("test-ns")

	job := &Naisjob{
		ObjectMeta: metav1.ObjectMeta{Name: "test-job", Namespace: "test-ns"},
		Spec:       NaisjobSpec{Uses: &Uses{OpenSearch: []OpenSearchUse{{Name: "search", Permission: "read"}}}},
	}

	t.Run("accepts existing namespace-local OpenSearch", func(t *testing.T) {
		validator := &JobValidator{Client: fakeKubeClient(openSearch)}
		_, err := validator.ValidateCreate(t.Context(), job)
		require.NoError(t, err)
	})

	t.Run("rejects missing OpenSearch", func(t *testing.T) {
		validator := &JobValidator{Client: fakeKubeClient()}
		_, err := validator.ValidateCreate(t.Context(), job)
		require.ErrorContains(t, err, "OpenSearch 'search' does not exist")
	})

	t.Run("rejects OpenSearch in another namespace", func(t *testing.T) {
		other := openSearch.DeepCopy()
		other.SetNamespace("other-ns")
		validator := &JobValidator{Client: fakeKubeClient(other)}
		_, err := validator.ValidateCreate(t.Context(), job)
		require.ErrorContains(t, err, "OpenSearch 'search' does not exist")
	})
}
