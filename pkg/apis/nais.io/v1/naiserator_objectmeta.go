package nais_io_v1

import (
	"encoding/base32"
	"encoding/binary"
	"hash/crc32"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func CreateObjectMeta(name, namespace, correlationID string, inLabels map[string]string, ownerReferences []metav1.OwnerReference) metav1.ObjectMeta {
	labels := map[string]string{}

	for k, v := range inLabels {
		labels[k] = v
	}

	labels["app"] = name

	return metav1.ObjectMeta{
		Name:      name,
		Namespace: namespace,
		Labels:    labels,
		Annotations: map[string]string{
			DeploymentCorrelationIDAnnotation: correlationID,
		},
		OwnerReferences: ownerReferences,
	}
}

// We concatenate name, namespace and add a hash in order to avoid duplicate names when creating service accounts in common service accounts namespace.
// Also making sure to not exceed name length restrictions of 30 characters
func CreateAppNamespaceHash(originalName, originalNamespace string) string {
	name := originalName
	namespace := originalNamespace
	if len(name) > 11 {
		name = originalName[:11]
	}
	if len(namespace) > 10 {
		namespace = originalNamespace[:10]
	}
	appNameSpace := name + "-" + namespace

	checksum := crc32.ChecksumIEEE([]byte(originalName + "-" + originalNamespace))
	bs := make([]byte, 4)
	binary.LittleEndian.PutUint32(bs, checksum)

	return appNameSpace + "-" + strings.ToLower(base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(bs))
}
