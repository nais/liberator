//go:build tools

package tools

import (
	_ "github.com/vektra/mockery/v2"
	_ "sigs.k8s.io/controller-tools/cmd/controller-gen"
)
