package nais_io_v1alpha1

import (
	"encoding/base32"
	"encoding/binary"
	"hash/crc32"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (in *Application) CreateObjectMeta() metav1.ObjectMeta {
	meta := metav1.ObjectMeta{
		Name:      in.Name,
		Namespace: in.Namespace,
		Labels: map[string]string{},
		Annotations: map[string]string{
			DeploymentCorrelationIDAnnotation: in.CorrelationID(),
		},
		OwnerReferences: in.OwnerReferences(in),
	}

	for k, v := range in.Labels {
		meta.Labels[k] = v
	}

	meta.Labels["app"] = in.Name

	return meta
}

// We concatenate name, namespace and add a hash in order to avoid duplicate names when creating service accounts in common service accounts namespace.
// Also making sure to not exceed name length restrictions of 30 characters
func (in *Application) CreateAppNamespaceHash() string {
	name := in.Name
	namespace := in.Namespace
	if len(name) > 11 {
		name = in.Name[:11]
	}
	if len(namespace) > 10 {
		namespace = in.Namespace[:10]
	}
	appNameSpace := name + "-" + namespace

	checksum := crc32.ChecksumIEEE([]byte(in.Name + "-" + in.Namespace))
	bs := make([]byte, 4)
	binary.LittleEndian.PutUint32(bs, checksum)

	return appNameSpace + "-" + strings.ToLower(base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(bs))
}

func (in *Application) CreateObjectMetaWithName(name string) metav1.ObjectMeta {
	m := in.CreateObjectMeta()
	m.Name = name
	return m
}

func (in *Application) OwnerReferences(app *Application) []metav1.OwnerReference {
	return []metav1.OwnerReference{app.GetOwnerReference()}
}
