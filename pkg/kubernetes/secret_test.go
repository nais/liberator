package kubernetes_test

import (
	"github.com/nais/liberator/pkg/kubernetes"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

func TestOpaqueSecret(t *testing.T) {
	objectMeta := kubernetes.ObjectMeta(
		"some-secret",
		"some-namespace",
		nil,
	)
	stringData := map[string]string{
		"some-key": "some-value",
	}

	secret := kubernetes.OpaqueSecret(objectMeta, stringData)

	assert.Equal(t, secret.TypeMeta.Kind, "Secret")
	assert.Equal(t, secret.TypeMeta.APIVersion, "v1")

	assert.Equal(t, secret.ObjectMeta, objectMeta)

	assert.Equal(t, secret.Type, corev1.SecretTypeOpaque)

	assert.Contains(t, secret.StringData, "some-key")
	assert.Len(t, secret.StringData, 1)
	assert.Equal(t, secret.StringData["some-key"], "some-value")
}

func TestListUsedAndUnusedSecretsForPods(t *testing.T) {
	usedEnvSecret := kubernetes.OpaqueSecret(kubernetes.ObjectMeta("used-env", "default", nil), nil)
	usedFileSecret := kubernetes.OpaqueSecret(kubernetes.ObjectMeta("used-file", "default", nil), nil)
	unusedSecret := kubernetes.OpaqueSecret(kubernetes.ObjectMeta("unused", "default", nil), nil)

	secretList := corev1.SecretList{
		TypeMeta: metav1.TypeMeta{
			Kind:       "SecretList",
			APIVersion: "v1",
		},
		Items: []corev1.Secret{usedEnvSecret, usedFileSecret, unusedSecret},
	}

	podWithEnvSecret := corev1.Pod{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Pod",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "pod-with-env-secret",
			Namespace: "default",
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:  "main",
					Image: "foo",
					EnvFrom: []corev1.EnvFromSource{
						{
							SecretRef: &corev1.SecretEnvSource{
								LocalObjectReference: corev1.LocalObjectReference{
									Name: "used-env",
								},
							},
						},
					},
				},
			},
		},
	}
	podWithFileSecret := corev1.Pod{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Pod",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "pod-with-file-secret",
			Namespace: "default",
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:  "main",
					Image: "foo",
				},
			},
			Volumes: []corev1.Volume{
				{
					Name: "foo",
					VolumeSource: corev1.VolumeSource{
						Secret: &corev1.SecretVolumeSource{
							SecretName: "used-file",
						},
					},
				},
			},
		},
	}

	podList := corev1.PodList{
		TypeMeta: metav1.TypeMeta{
			Kind:       "SecretList",
			APIVersion: "v1",
		},
		Items: []corev1.Pod{
			podWithEnvSecret,
			podWithFileSecret,
		},
	}

	secretLists := kubernetes.ListUsedAndUnusedSecretsForPods(secretList, podList)

	assert.Len(t, secretLists.Used.Items, 2)
	assert.Len(t, secretLists.Unused.Items, 1)

	assert.Equal(t, secretLists.Used.Items[0], usedEnvSecret)
	assert.Equal(t, secretLists.Used.Items[1], usedFileSecret)
	assert.Equal(t, secretLists.Unused.Items[0], unusedSecret)
}
