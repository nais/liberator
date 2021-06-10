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
func UniformResourceScopeName(resource metav1.Object, product, subScope string) string {
	return replace(fmt.Sprintf("%s:%s:%s.%s", resource.GetClusterName(), resource.GetNamespace(), resource.GetName(), ToScope(product, subScope)))
}

func replace(x string) string {
	return strings.ReplaceAll(x, "-", "")
}

func ToScope(product, subScope string) string {
	if !strings.Contains(subScope, "/") {
		return toScope(product, subScope, ":")
	}
	// able to use legacy scopes from on-prem in gcp
	return toScope(product, subScope, "/")
}

func toScope(product, subScope, separator string) string {
	return fmt.Sprintf("%s%s%s", product, separator, subScope)
}
