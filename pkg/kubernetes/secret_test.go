package kubernetes_test

import (
	"context"
	"testing"

	"github.com/nais/liberator/pkg/kubernetes"
	"github.com/nais/liberator/pkg/scheme"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
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

func TestListUsedAndUnusedSecretsForApp(t *testing.T) {
	appName := "some-app"
	namespace := "some-namespace"

	secretLabels := map[string]string{
		"some-key": "some-value",
	}
	makeSecret := func(name string) corev1.Secret {
		return kubernetes.OpaqueSecret(kubernetes.ObjectMeta(name, namespace, secretLabels), nil)
	}
	secretUsedByPod := makeSecret("used-by-pod")
	secretUsedByReplicaSet := makeSecret("used-by-replica-set")
	secretUnused := makeSecret("unused")
	secretUnrelated := kubernetes.OpaqueSecret(kubernetes.ObjectMeta("unrelated", namespace, nil), nil)

	makePodSpec := func(secretName string) corev1.PodSpec {
		return corev1.PodSpec{Containers: []corev1.Container{{
			Name:  "main",
			Image: "foo",
			EnvFrom: []corev1.EnvFromSource{{SecretRef: &corev1.SecretEnvSource{
				LocalObjectReference: corev1.LocalObjectReference{Name: secretName},
			}}},
		}}}
	}

	appLabels := map[string]string{
		"app": appName,
	}
	pod := &corev1.Pod{
		ObjectMeta: kubernetes.ObjectMeta("some-app-cafebabe-abcde", namespace, appLabels),
		Spec:       makePodSpec("used-by-pod"),
	}
	replicaSet := &appsv1.ReplicaSet{
		ObjectMeta: kubernetes.ObjectMeta("some-app-cafebabe", namespace, appLabels),
		Spec: appsv1.ReplicaSetSpec{
			Template: corev1.PodTemplateSpec{Spec: makePodSpec("used-by-replica-set")},
		},
	}

	s, err := scheme.All()
	require.NoError(t, err)

	cli := fake.NewClientBuilder().
		WithScheme(s).
		WithObjects(
			&secretUsedByPod,
			&secretUsedByReplicaSet,
			&secretUnused,
			&secretUnrelated,
			pod,
			replicaSet,
		).
		Build()

	secrets, err := kubernetes.ListSecretsForApplication(context.Background(), cli, client.ObjectKey{
		Name:      appName,
		Namespace: namespace,
	}, secretLabels)
	require.NoError(t, err)

	assert.Len(t, secrets.Used.Items, 2)
	assert.Equal(t, secrets.Used.Items[0], secretUsedByPod)
	assert.Equal(t, secrets.Used.Items[1], secretUsedByReplicaSet)

	assert.Len(t, secrets.Unused.Items, 1)
	assert.Equal(t, secrets.Unused.Items[0], secretUnused)
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
