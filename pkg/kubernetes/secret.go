package kubernetes

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func OpaqueSecret(objectMeta metav1.ObjectMeta, stringData map[string]string) corev1.Secret {
	return corev1.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		ObjectMeta: objectMeta,
		StringData: stringData,
		Type:       corev1.SecretTypeOpaque,
	}
}

type SecretLists struct {
	Used   corev1.SecretList
	Unused corev1.SecretList
}

// ListSecretsForApplication is like ListUsedAndUnusedSecretsForPods, but lists secrets and pods based on the given
// secret label selectors and the application object key, respectively.
//
// Pods are listed using the label "app", which should equal the given application object name.
// ReplicaSets are also listed to preserve rollback functionality, using the same labels.
//
// The predicates used to determine if a secret is used are non-exhaustive. We currently check for references in:
// - Pod volumes
// - For containers, init containers, and ephemeral containers:
//   - Individual environment variables (env)
//   - Sources for environment variables (envFrom)
//
// The caller should verify that secrets aren't referenced elsewhere, especially in custom resources.
func ListSecretsForApplication(ctx context.Context, reader client.Reader, application client.ObjectKey, secretLabels client.MatchingLabels) (SecretLists, error) {
	var lists SecretLists

	secrets := corev1.SecretList{}
	if err := reader.List(ctx, &secrets, client.InNamespace(application.Namespace), secretLabels); err != nil {
		return lists, err
	}

	pods, err := ListPodsForApplication(ctx, reader, application.Name, application.Namespace)
	if err != nil {
		return lists, err
	}

	replicaSets, err := ListReplicaSetsForApplication(ctx, reader, application.Name, application.Namespace)
	if err != nil {
		return lists, err
	}

	podSpecs := append(
		extractPodSpecs(&pods),
		extractPodSpecs(&replicaSets)...,
	)
	return partitionSecrets(secrets, podSpecs), nil
}

// ListUsedAndUnusedSecretsForPods partitions the given secrets into used and unused lists based
// on whether they are referenced by the given pods.
//
// The predicates used to determine if a secret is used are non-exhaustive. We currently check for references in:
// - Pod volumes
// - For containers, init containers, and ephemeral containers:
//   - Individual environment variables (env)
//   - Sources for environment variables (envFrom)
//
// The caller should verify that secrets aren't referenced elsewhere, especially in custom resources.
func ListUsedAndUnusedSecretsForPods(secrets corev1.SecretList, pods corev1.PodList) SecretLists {
	return partitionSecrets(secrets, extractPodSpecs(&pods))
}

func partitionSecrets(secrets corev1.SecretList, podSpecs []corev1.PodSpec) SecretLists {
	lists := SecretLists{
		Used: corev1.SecretList{
			Items: make([]corev1.Secret, 0),
		},
		Unused: corev1.SecretList{
			Items: make([]corev1.Secret, 0),
		},
	}

	for _, sec := range secrets.Items {
		if secretInPodSpecs(sec.Name, podSpecs) {
			lists.Used.Items = append(lists.Used.Items, sec)
		} else {
			lists.Unused.Items = append(lists.Unused.Items, sec)
		}
	}

	return lists
}

func secretInPodSpecs(secretName string, podSpecs []corev1.PodSpec) bool {
	for _, podSpec := range podSpecs {
		if secretInPodSpec(secretName, podSpec) {
			return true
		}
	}
	return false
}

func secretInPodSpec(secretName string, podSpec corev1.PodSpec) bool {
	return secretInVolumes(secretName, podSpec.Volumes) ||
		secretInContainers(secretName, podSpec.Containers) ||
		secretInContainers(secretName, podSpec.InitContainers) ||
		secretInContainers(secretName, asContainers(podSpec.EphemeralContainers))
}

func secretInVolumes(secretName string, volumes []corev1.Volume) bool {
	for _, volume := range volumes {
		if volume.Secret != nil && volume.Secret.SecretName == secretName {
			return true
		}
	}
	return false
}

func secretInContainers(secretName string, containers []corev1.Container) bool {
	for _, container := range containers {
		if secretInEnvVars(secretName, container.Env) || secretInEnvFromSources(secretName, container.EnvFrom) {
			return true
		}
	}
	return false
}

func secretInEnvVars(secretName string, envVars []corev1.EnvVar) bool {
	for _, env := range envVars {
		if env.ValueFrom != nil && env.ValueFrom.SecretKeyRef != nil {
			if env.ValueFrom.SecretKeyRef.Name == secretName {
				return true
			}
		}
	}
	return false
}

func secretInEnvFromSources(secretName string, envFrom []corev1.EnvFromSource) bool {
	for _, env := range envFrom {
		if env.SecretRef != nil && env.SecretRef.Name == secretName {
			return true
		}
	}
	return false
}

func asContainers(ephemeral []corev1.EphemeralContainer) []corev1.Container {
	containers := make([]corev1.Container, len(ephemeral))
	for i, e := range ephemeral {
		containers[i] = corev1.Container(e.EphemeralContainerCommon)
	}
	return containers
}
