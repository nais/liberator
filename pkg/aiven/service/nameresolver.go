package service

import (
	"context"
	"fmt"
)

type NameResolver interface {
	Interface
	ResolveKafkaServiceName(ctx context.Context, project string) (string, error)
}

type CachedNameResolver struct {
	Interface
	cache map[string]string
}

func (r *CachedNameResolver) ResolveKafkaServiceName(ctx context.Context, project string) (string, error) {
	if name, ok := r.cache[project]; ok {
		return name, nil
	}
	svcs, err := r.List(ctx, project)
	if err != nil {
		return "", err
	}

	candidates := []string{
		fmt.Sprintf("%s-kafka", project),
		"kafka",
	}

	for _, svc := range svcs {
		for _, candidate := range candidates {
			if svc.Name == candidate {
				r.cache[project] = svc.Name
				return svc.Name, nil
			}
		}
	}
	return "", fmt.Errorf("no kafka service found in project %s", project)
}

func NewCachedNameResolver(services Interface) *CachedNameResolver {
	return &CachedNameResolver{
		Interface: services,
		cache:     make(map[string]string, 0),
	}
}
