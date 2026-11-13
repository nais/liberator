package nais_io_v1alpha1

import (
	"testing"

	nais_io_v1 "github.com/nais/liberator/pkg/apis/nais.io/v1"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func TestApplicationValidator_ValkeyUses(t *testing.T) {
	valkey := &unstructured.Unstructured{}
	valkey.SetAPIVersion("nais.io/v1")
	valkey.SetKind("Valkey")
	valkey.SetName("cache")
	valkey.SetNamespace("test-ns")

	app := &Application{
		ObjectMeta: metav1.ObjectMeta{Name: "test-app", Namespace: "test-ns"},
		Spec:       ApplicationSpec{Uses: &nais_io_v1.Uses{Valkey: []nais_io_v1.ValkeyUse{{Name: "cache", Permission: "read"}}}},
	}

	t.Run("accepts existing namespace-local Valkey", func(t *testing.T) {
		validator := &ApplicationValidator{Client: fakeKubeClient(valkey)}
		_, err := validator.ValidateCreate(t.Context(), app)
		require.NoError(t, err)
	})

	t.Run("rejects missing Valkey", func(t *testing.T) {
		validator := &ApplicationValidator{Client: fakeKubeClient()}
		_, err := validator.ValidateCreate(t.Context(), app)
		require.ErrorContains(t, err, "Valkey 'cache' does not exist")
	})

	t.Run("rejects Valkey in another namespace", func(t *testing.T) {
		other := valkey.DeepCopy()
		other.SetNamespace("other-ns")
		validator := &ApplicationValidator{Client: fakeKubeClient(other)}
		_, err := validator.ValidateCreate(t.Context(), app)
		require.ErrorContains(t, err, "Valkey 'cache' does not exist")
	})
}
