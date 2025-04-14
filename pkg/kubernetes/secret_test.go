package kubernetes_test

import (
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

type secretInUseTest struct {
	name        string
	mutatePodFn func(secretName string, pod *corev1.Pod)
}

var secretInUseTests = []secretInUseTest{
	{
		name: "container-env-from",
		mutatePodFn: func(secretName string, pod *corev1.Pod) {
			pod.Spec.Containers[0].EnvFrom = secretEnvFromSources(secretName)
		},
	},
	{
		name: "container-env-var",
		mutatePodFn: func(secretName string, pod *corev1.Pod) {
			pod.Spec.Containers[0].Env = secretEnvVars(secretName)
		},
	},
	{
		name: "volume",
		mutatePodFn: func(secretName string, pod *corev1.Pod) {
			pod.Spec.Volumes = secretVolumes(secretName)
		},
	},
	{
		name: "init-container-env-from",
		mutatePodFn: func(secretName string, pod *corev1.Pod) {
			pod.Spec.InitContainers[0].EnvFrom = secretEnvFromSources(secretName)
		},
	},
	{
		name: "init-container-env-var",
		mutatePodFn: func(secretName string, pod *corev1.Pod) {
			pod.Spec.InitContainers[0].Env = secretEnvVars(secretName)
		},
	},
	{
		name: "ephemeral-container-env-from",
		mutatePodFn: func(secretName string, pod *corev1.Pod) {
			pod.Spec.EphemeralContainers[0].EnvFrom = secretEnvFromSources(secretName)
		},
	},
	{
		name: "ephemeral-container-env-var",
		mutatePodFn: func(secretName string, pod *corev1.Pod) {
			pod.Spec.EphemeralContainers[0].Env = secretEnvVars(secretName)
		},
	},
}

func TestListUsedAndUnusedSecretsForApp(t *testing.T) {
	app := client.ObjectKey{
		Name:      "some-app",
		Namespace: "some-namespace",
	}
	appLabels := map[string]string{"app": app.Name}
	secretLabels := map[string]string{"some-key": "some-value"}

	for _, tt := range secretInUseTests {
		t.Run(tt.name, func(t *testing.T) {
			secretUsedByPod := makeSecret("used-by-pod", secretLabels)
			secretUsedByReplicaSet := makeSecret("used-by-replica-set", secretLabels)
			secretUnused := makeSecret("unused", secretLabels)
			secretUnrelated := makeSecret("unrelated", nil)

			pod := makePod("some-app-cafebabe-abcde", appLabels)
			tt.mutatePodFn(secretUsedByPod.Name, pod)

			replicaSetPod := makePod("some-app-deadbeef-abcde", appLabels)
			tt.mutatePodFn(secretUsedByReplicaSet.Name, replicaSetPod)
			replicaSet := &appsv1.ReplicaSet{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "some-app-deadbeef",
					Namespace: app.Namespace,
					Labels:    appLabels,
				},
				Spec: appsv1.ReplicaSetSpec{
					Template: corev1.PodTemplateSpec{
						Spec: replicaSetPod.Spec,
					},
				},
			}

			s, err := scheme.All()
			require.NoError(t, err)

			cli := fake.NewClientBuilder().
				WithScheme(s).
				WithObjects(secretUsedByPod, secretUsedByReplicaSet, secretUnused, secretUnrelated, pod, replicaSet).
				Build()

			secrets, err := kubernetes.ListSecretsForApplication(t.Context(), cli, app, secretLabels)
			require.NoError(t, err)

			assert.ElementsMatch(t, secrets.Used.Items, []corev1.Secret{
				*secretUsedByPod,
				*secretUsedByReplicaSet,
			})
			assert.ElementsMatch(t, secrets.Unused.Items, []corev1.Secret{
				*secretUnused,
			})
		})
	}
}

func TestListUsedAndUnusedSecretsForPods(t *testing.T) {
	for _, tt := range secretInUseTests {
		t.Run(tt.name, func(t *testing.T) {
			secretUsed := makeSecret(tt.name, nil)
			secretUnused := makeSecret("unused", nil)
			pod := makePod(tt.name, nil)
			tt.mutatePodFn(secretUsed.Name, pod)

			lists := kubernetes.ListUsedAndUnusedSecretsForPods(corev1.SecretList{
				Items: []corev1.Secret{
					*secretUsed,
					*secretUnused,
				},
			}, corev1.PodList{
				Items: []corev1.Pod{
					*pod,
				},
			})

			assert.ElementsMatch(t, lists.Used.Items, []corev1.Secret{*secretUsed})
			assert.ElementsMatch(t, lists.Unused.Items, []corev1.Secret{*secretUnused})
		})
	}
}

func makeSecret(name string, labels map[string]string) *corev1.Secret {
	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: "some-namespace",
			Labels:    labels,
		},
		Type: corev1.SecretTypeOpaque,
	}
}

func makePod(name string, labels map[string]string) *corev1.Pod {
	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: "some-namespace",
			Labels:    labels,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:  "main",
					Image: "foo",
				},
			},
			InitContainers: []corev1.Container{
				{
					Name:  "init",
					Image: "bar",
				},
			},
			EphemeralContainers: []corev1.EphemeralContainer{
				{
					EphemeralContainerCommon: corev1.EphemeralContainerCommon{
						Name:  "ephemeral",
						Image: "baz",
					},
					TargetContainerName: "main",
				},
			},
		},
	}
}

func secretEnvFromSources(name string) []corev1.EnvFromSource {
	return []corev1.EnvFromSource{
		{
			SecretRef: &corev1.SecretEnvSource{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: name,
				},
			},
		},
	}
}

func secretEnvVars(name string) []corev1.EnvVar {
	return []corev1.EnvVar{
		{
			Name: name,
			ValueFrom: &corev1.EnvVarSource{
				SecretKeyRef: &corev1.SecretKeySelector{
					Key: name,
					LocalObjectReference: corev1.LocalObjectReference{
						Name: name,
					},
				},
			},
		},
	}
}

func secretVolumes(name string) []corev1.Volume {
	return []corev1.Volume{
		{
			Name: name,
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName: name,
				},
			},
		},
	}
}
