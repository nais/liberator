package kafka_nais_io_v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
)

func ExampleTopicForDocumentation() *Topic {
	return &Topic{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Topic",
			APIVersion: "kafka.nais.io/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "mytopic",
			Namespace: "myteam",
			Labels: map[string]string{
				"team": "myteam",
			},
		},
		Spec: TopicSpec{
			Pool: "dev-nais-dev",
			Config: &Config{
				CleanupPolicy:         ptr.To("delete"),
				MinimumInSyncReplicas: ptr.To(2),
				Partitions:            ptr.To(1),
				Replication:           ptr.To(3),
				RetentionBytes:        ptr.To(-1),
				RetentionHours:        ptr.To(168),
				SegmentHours:          ptr.To(168),
				MaxMessageBytes:       ptr.To(1048588),
			},
			ACL: TopicACLs{
				{
					Access:      "read",
					Application: "consumer",
					Team:        "otherteam",
				},
				{
					Access:      "write",
					Application: "producer",
					Team:        "myteam",
				},
				{
					Access:      "readwrite",
					Application: "processor",
					Team:        "myteam",
				},
			},
		},
	}
}
