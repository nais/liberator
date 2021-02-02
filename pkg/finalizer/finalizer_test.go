package finalizer_test

import (
	"github.com/nais/liberator/pkg/finalizer"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

func TestIsBeingDeleted(t *testing.T) {
	instance := someResourceInstance()

	t.Run("Resource without deletion marker should not be marked for deletion", func(t *testing.T) {
		assert.False(t, finalizer.IsBeingDeleted(instance))
	})

	t.Run("Resource with deletion marker should be marked for deletion", func(t *testing.T) {
		now := metav1.Now()
		instance.ObjectMeta.DeletionTimestamp = &now
		assert.True(t, finalizer.IsBeingDeleted(instance))
	})
}

func TestHasFinalizer(t *testing.T) {
	instance := someResourceInstance()
	finalizerName := "some-finalizer"

	t.Run("Resource should not have finalizer", func(t *testing.T) {
		assert.False(t, finalizer.HasFinalizer(instance, finalizerName))
	})

	t.Run("Resource with finalizer should have finalizer", func(t *testing.T) {
		instance.ObjectMeta.Finalizers = []string{finalizerName}
		assert.True(t, finalizer.HasFinalizer(instance, finalizerName))
	})
}

type someResource struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
}

func someResourceInstance() *someResource {
	return &someResource{
		ObjectMeta: metav1.ObjectMeta{
			Name:        "test-app",
			Namespace:   "test-namespace",
			ClusterName: "test-cluster",
		},
	}
}
