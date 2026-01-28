package kafka_nais_io_v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func StreamDocExample() *Stream {
	return &Stream{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "example-stream",
			Namespace: "example-namespace",
		},
		Spec: StreamSpec{
			Pool: "example-pool",
			AdditionalUsers: []AdditionalStreamUser{
				{
					Username: "developer1",
				},
			},
		},
	}
}
