package nais_io_v1

import (
	"testing"

	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func TestJobValidator_ValkeyUses(t *testing.T) {
	valkey := &unstructured.Unstructured{}
	valkey.SetAPIVersion("nais.io/v1")
	valkey.SetKind("Valkey")
	valkey.SetName("cache")
	valkey.SetNamespace("test-ns")

	job := &Naisjob{
		ObjectMeta: metav1.ObjectMeta{Name: "test-job", Namespace: "test-ns"},
		Spec:       NaisjobSpec{Uses: &Uses{Valkey: []ValkeyUse{{Name: "cache", Permission: "read"}}}},
	}

	t.Run("accepts existing namespace-local Valkey", func(t *testing.T) {
		validator := &JobValidator{Client: fakeKubeClient(valkey)}
		_, err := validator.ValidateCreate(t.Context(), job)
		require.NoError(t, err)
	})

	t.Run("rejects missing Valkey", func(t *testing.T) {
		validator := &JobValidator{Client: fakeKubeClient()}
		_, err := validator.ValidateCreate(t.Context(), job)
		require.ErrorContains(t, err, "Valkey 'cache' does not exist")
	})

	t.Run("rejects Valkey in another namespace", func(t *testing.T) {
		other := valkey.DeepCopy()
		other.SetNamespace("other-ns")
		validator := &JobValidator{Client: fakeKubeClient(other)}
		_, err := validator.ValidateCreate(t.Context(), job)
		require.ErrorContains(t, err, "Valkey 'cache' does not exist")
	})
}
