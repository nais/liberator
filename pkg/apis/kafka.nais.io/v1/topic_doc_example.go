package kafka_nais_io_v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
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
			Pool: "nav-integration-test",
			Config: &Config{
				CleanupPolicy:         pointer.StringPtr("delete"),
				MinimumInSyncReplicas: pointer.IntPtr(2),
				Partitions:            pointer.IntPtr(1),
				Replication:           pointer.IntPtr(3),
				RetentionBytes:        pointer.IntPtr(-1),
				RetentionHours:        pointer.IntPtr(168),
				SegmentHours:          pointer.IntPtr(168),
				MaxMessageBytes:       pointer.IntPtr(1048588),
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
