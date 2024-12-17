package kubernetes

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

func extractPodSpecs(objects metav1.ListInterface) []corev1.PodSpec {
	specs := make([]corev1.PodSpec, 0)

	switch list := objects.(type) {
	case *corev1.PodList:
		for _, pod := range list.Items {
			specs = append(specs, pod.Spec)
		}
	case *appsv1.ReplicaSetList:
		for _, replicaSet := range list.Items {
			specs = append(specs, replicaSet.Spec.Template.Spec)
		}
	}

	return specs
}
