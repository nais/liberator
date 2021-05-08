package kubernetes

import (
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func UniformResourceName(resource metav1.Object) string {
	return fmt.Sprintf("%s:%s:%s", resource.GetClusterName(), resource.GetNamespace(), resource.GetName())
}

func ObjectMeta(name, namespace string, labels map[string]string) metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Name:      name,
		Namespace: namespace,
		Labels:    labels,
	}
}
func UniformResourceScopeName(resource metav1.Object, scope string) string {
	return fmt.Sprintf("%s:%s:%s", resource.GetClusterName(), resource.GetNamespace(), scope)
}
