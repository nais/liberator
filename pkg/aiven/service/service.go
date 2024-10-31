package service

import (
	"context"

	"github.com/aiven/aiven-go-client/v2"
)

type Interface interface {
	List(ctx context.Context, project string) ([]*aiven.Service, error)
}

var _ Interface = &aiven.ServicesHandler{}
