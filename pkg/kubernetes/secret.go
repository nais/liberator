package kubernetes

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

// ListUsedAndUnusedSecretsForPods finds intersect between list of secrets and list of pods
// that uses (i.e. mounts or refers to the secret) these secrets,
// and separates the secret list into two lists; used and unused.
func ListUsedAndUnusedSecretsForPods(secrets corev1.SecretList, pods corev1.PodList) SecretLists {
	lists := SecretLists{
		Used: corev1.SecretList{
			Items: make([]corev1.Secret, 0),
		},
		Unused: corev1.SecretList{
			Items: make([]corev1.Secret, 0),
		},
	}

	for _, sec := range secrets.Items {
		if secretInPods(sec, pods) {
			lists.Used.Items = append(lists.Used.Items, sec)
		} else {
			lists.Unused.Items = append(lists.Unused.Items, sec)
		}
	}
	return lists
}

func secretInPods(secret corev1.Secret, pods corev1.PodList) bool {
	for _, pod := range pods.Items {
		if secretInPod(secret, pod) {
			return true
		}
	}
	return false
}

func secretInPod(secret corev1.Secret, pod corev1.Pod) bool {
	return secretRefInVolumes(secret, pod.Spec.Volumes) || secretRefInContainers(secret, pod.Spec.Containers)
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
