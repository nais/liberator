package kubernetes

import (
	"fmt"
	"strings"

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
	return replace(fmt.Sprintf("%s:%s:%s.%s", resource.GetClusterName(), resource.GetNamespace(), resource.GetName(), scope))
}

func replace(x string) string {
	return strings.ReplaceAll(x, "-", "")
}

func FilterUniformedName(resource metav1.Object, subScope string) string {
	if !strings.Contains(subScope, "/") {
		return filterClusterPrefix(resource, subScope)
	}
	// able to use legacy scopes from on-prem in gcp
	return subScope
}

func filterClusterPrefix(resource metav1.Object, subScope string) string {
	return strings.TrimPrefix(UniformResourceScopeName(resource, subScope), fmt.Sprintf("%s:", replace(resource.GetClusterName())))
}
