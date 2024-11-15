package service

import (
	"context"

	"github.com/aiven/aiven-go-client/v2"
)

type Interface interface {
	Get(ctx context.Context, project, service string) (*aiven.Service, error)
}

var _ Interface = &aiven.ServicesHandler{}
