package kafka_nais_io_v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
				CleanupPolicy:         new("delete"),
				DeleteRetentionHours:  new(24),
				MinimumInSyncReplicas: new(2),
				Partitions:            new(1),
				Replication:           new(3),
				RetentionBytes:        new(6000),
				RetentionHours:        new(168),
				LocalRetentionBytes:   new(1000),
				LocalRetentionHours:   new(68),
				SegmentHours:          new(168),
				MaxMessageBytes:       new(1048588),
				MaxCompactionLagMs:    new(60000),
				MinCompactionLagMs:    new(10000),
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
