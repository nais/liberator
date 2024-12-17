package kubernetes

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func ListReplicaSetsForApplication(ctx context.Context, reader client.Reader, name, namespace string) (appsv1.ReplicaSetList, error) {
	replicaSetList := appsv1.ReplicaSetList{}
	opts := []client.ListOption{
		client.MatchingLabels{
			"app": name,
		},
		client.InNamespace(namespace),
	}
	err := reader.List(ctx, &replicaSetList, opts...)
	if err != nil {
		return replicaSetList, err
	}

	return replicaSetList, nil
}
