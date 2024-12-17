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

// ListSecretsForApplication finds all secrets matching the given labels and partitions them into SecretLists.
// A secret is considered "used" if it is mounted or referred to by a pod or any replica set that matches the
// given application.
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

// ListUsedAndUnusedSecretsForPods finds intersect between list of secrets and list of pods
// that uses (i.e. mounts or refers to the secret) these secrets,
// and separates the secret list into two lists; used and unused.
//
// This is a low-level method that requires you to perform the list operations for the secrets and pod objects yourself.
// In most cases, you should prefer using ListSecretsForApplication instead.
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
		if secretInPodSpecs(sec, podSpecs) {
			lists.Used.Items = append(lists.Used.Items, sec)
		} else {
			lists.Unused.Items = append(lists.Unused.Items, sec)
		}
	}

	return lists
}

func secretInPodSpecs(secret corev1.Secret, podSpecs []corev1.PodSpec) bool {
	for _, podSpec := range podSpecs {
		if secretInPodSpec(secret, podSpec) {
			return true
		}
	}
	return false
}

func secretInPodSpec(secret corev1.Secret, podSpec corev1.PodSpec) bool {
	return secretRefInVolumes(secret, podSpec.Volumes) || secretRefInContainers(secret, podSpec.Containers)
}

func secretRefInVolumes(secret corev1.Secret, volumes []corev1.Volume) bool {
	for _, volume := range volumes {
		if volume.Secret != nil && volume.Secret.SecretName == secret.Name {
			return true
		}
	}
	return false
}

func secretRefInContainers(secret corev1.Secret, containers []corev1.Container) bool {
	for _, container := range containers {
		if secretRefInEnvFromSource(secret, container.EnvFrom) {
			return true
		}
	}
	return false
}

func secretRefInEnvFromSource(secret corev1.Secret, envFromSource []corev1.EnvFromSource) bool {
	for _, envFrom := range envFromSource {
		if envFrom.SecretRef != nil && envFrom.SecretRef.Name == secret.Name {
			return true
		}
	}
	return false
}
