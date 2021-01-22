package crd_test

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/nais/liberator/pkg/crd"
)

// Test for one or more YAML files in this directory
func TestYamlDirectory(t *testing.T) {
	dir := crd.YamlDirectory()
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	fc := 0
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".yaml") {
			fc++
		}
	}
	if fc == 0 {
		t.Errorf("No YAML files found in CRD directory '%s'", dir)
		t.Fail()
	}
	if fc != len(files) {
		t.Errorf("Non-YAML files found in CRD directory '%s'", dir)
		t.Fail()
	}
}
