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

func FilterUniformedName(resource metav1.Object, uniformedName, subScope string) string {
	if !strings.Contains(uniformedName, "/") {
		return filterClusterPrefix(resource, uniformedName)
	}
	// able to use legacy scopes from on-prem in gcp
	return subScope
}

func filterClusterPrefix(resource metav1.Object, uniformedScopeName string) string {
	return strings.TrimPrefix(uniformedScopeName, fmt.Sprintf("%s:", resource.GetClusterName()))
}
