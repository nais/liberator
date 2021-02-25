package kubernetes

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func ListPods(ctx context.Context, reader client.Reader, opts ...client.ListOption) (corev1.PodList, error) {
	podList := corev1.PodList{}
	err := reader.List(ctx, &podList, opts...)
	if err != nil {
		return podList, err
	}
	return podList, nil
}

func ListPodsForApplication(ctx context.Context, reader client.Reader, name, namespace string) (corev1.PodList, error) {
	matchingLabels := client.MatchingLabels{
		"app": name,
	}
	namespaceSelector := client.InNamespace(namespace)
	return ListPods(ctx, reader, matchingLabels, namespaceSelector)
}
