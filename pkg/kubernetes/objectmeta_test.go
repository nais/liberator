package kubernetes_test

import (
	"github.com/nais/liberator/pkg/kubernetes"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

func TestUniformResourceName(t *testing.T) {
	om := &metav1.ObjectMeta{
		Name:        "test-app",
		Namespace:   "test-namespace",
		ClusterName: "test-cluster",
	}
	expected := "test-cluster:test-namespace:test-app"
	assert.Equal(t, expected, kubernetes.UniformResourceName(om))
}

func TestObjectMeta(t *testing.T) {
	objectMeta := kubernetes.ObjectMeta(
		"some-secret",
		"some-namespace",
		map[string]string{
			"some-key": "some-value",
		},
	)

	assert.Equal(t, objectMeta.GetName(), "some-secret")

	assert.Equal(t, objectMeta.GetNamespace(), "some-namespace")

	assert.Contains(t, objectMeta.GetLabels(), "some-key")
	assert.Equal(t, objectMeta.GetLabels()["some-key"], "some-value")
	assert.Len(t, objectMeta.GetLabels(), 1)
}
