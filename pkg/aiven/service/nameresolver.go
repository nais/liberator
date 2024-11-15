package service

import (
	"context"
	"fmt"

	"github.com/aiven/aiven-go-client/v2"
)

type NameResolver interface {
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

	candidates := []string{
		"kafka",
		fmt.Sprintf("%s-kafka", project),
	}

	lookupErrors := make([]error, 0, 2)
	for _, candidate := range candidates {
		svc, err := r.Get(ctx, project, candidate)
		if err == nil {
			r.cache[project] = svc.Name
			return svc.Name, nil
		} else if !aiven.IsNotFound(err) {
			lookupErrors = append(lookupErrors, err)
		}
	}

	if len(lookupErrors) > 0 {
		return "", fmt.Errorf("failed to lookup kafka service in project %s: %w", project, lookupErrors[0])
	}

	return "", fmt.Errorf("no kafka service found in project %s", project)
}

func NewCachedNameResolver(services Interface) *CachedNameResolver {
	return &CachedNameResolver{
		Interface: services,
		cache:     make(map[string]string, 0),
	}
}

var _ NameResolver = &CachedNameResolver{}
