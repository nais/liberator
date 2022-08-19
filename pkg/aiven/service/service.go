package service

import (
	"github.com/aiven/aiven-go-client"
)

type Interface interface {
	List(project string) ([]*aiven.Service, error)
}
