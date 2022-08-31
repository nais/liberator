package kubernetes_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/nais/liberator/pkg/kubernetes"
)

func TestUniformResourceName(t *testing.T) {
	clusterName := "test-cluster"
	om := &metav1.ObjectMeta{
		Name:      "test-app",
		Namespace: "test-namespace",
	}
	expected := "test-cluster:test-namespace:test-app"
	assert.Equal(t, expected, kubernetes.UniformResourceName(om, clusterName))
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
