package crd_test

import (
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"testing"

	"github.com/nais/liberator/pkg/crd"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"sigs.k8s.io/yaml"
)

func TestWorkloadPostgresBranchSchema(t *testing.T) {
	for _, name := range []string{"nais.io_applications.yaml", "nais.io_naisjobs.yaml"} {
		t.Run(name, func(t *testing.T) {
			data, err := os.ReadFile(filepath.Join(crd.YamlDirectory(), name))
			if err != nil {
				t.Fatal(err)
			}
			var resource apiextensionsv1.CustomResourceDefinition
			if err := yaml.Unmarshal(data, &resource); err != nil {
				t.Fatal(err)
			}
			for _, version := range resource.Spec.Versions {
				t.Run(version.Name, func(t *testing.T) {
					use := version.Schema.OpenAPIV3Schema.Properties["spec"].Properties["uses"].Properties["postgres"].Items.Schema
					if use == nil {
						t.Fatal("missing uses.postgres schema")
					}
					branch, found := use.Properties["branch"]
					if !found || branch.Type != "string" || slices.Contains(use.Required, "branch") {
						t.Fatal("branch must be an optional string for both workload kinds")
					}
					if branch.MinLength == nil || *branch.MinLength != 1 || branch.MaxLength == nil || *branch.MaxLength != 63 {
						t.Fatal("branch must have the same name limits as PostgresBranch")
					}
					pattern, err := regexp.Compile(branch.Pattern)
					if err != nil {
						t.Fatal(err)
					}
					for _, name := range []string{"main", "pr-123", "a", strings.Repeat("a", 63)} {
						if !pattern.MatchString(name) {
							t.Errorf("valid branch %q rejected", name)
						}
					}
					for _, name := range []string{"", "INVALID", "other/branch", "-branch", "branch-"} {
						if pattern.MatchString(name) {
							t.Errorf("invalid branch %q accepted", name)
						}
					}
				})
			}
		})
	}
}
