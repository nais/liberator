package crd

import (
	"path/filepath"
	"runtime"
)

const relativePath = "../../config/crd/bases"

func YamlDirectory() string {
	_, filename, _, _ := runtime.Caller(0)
	return filepath.Clean(filepath.Join(filepath.Dir(filename), relativePath))
}
