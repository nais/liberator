package finalizer

import (
	"github.com/nais/liberator/pkg/strings"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func IsBeingDeleted(instance metav1.Object) bool {
	return !instance.GetDeletionTimestamp().IsZero()
}

func HasFinalizer(instance metav1.Object, finalizerName string) bool {
	return strings.ContainsString(instance.GetFinalizers(), finalizerName)
}
