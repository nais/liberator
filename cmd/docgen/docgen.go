package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	apiext "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-tools/pkg/crd"
	crd_markers "sigs.k8s.io/controller-tools/pkg/crd/markers"
	"sigs.k8s.io/controller-tools/pkg/loader"
	"sigs.k8s.io/controller-tools/pkg/markers"
)

// Generate documentation for NAIS CRD's

type Config struct {
	Directory string
	BaseClass string
	Group     string
	Kind      string
	Output    string
}

type Node struct {
	Name     string
	Type     string
	Doc      string
	Array    bool
	JsonTag  string
	Required bool
	Markers  []string
	Children []*Node
}

func main() {
	err := run()
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
}

func run() error {
	log.SetLevel(log.DebugLevel)

	cfg := &Config{}
	pflag.StringVar(&cfg.Directory, "dir", cfg.Directory, "directory with packages")
	pflag.StringVar(&cfg.Group, "group", cfg.Group, "which group to generate documentation for")
	pflag.StringVar(&cfg.Kind, "kind", cfg.Kind, "which kind to generate documentation for")
	pflag.StringVar(&cfg.Output, "output", cfg.Output, "markdown output file")
	pflag.Parse()

	packages, err := loader.LoadRoots(cfg.Directory)
	if err != nil {
		return err
	}
	registry := &markers.Registry{}
	collector := &markers.Collector{
		Registry: registry,
	}
	err = crd_markers.Register(registry)
	if err != nil {
		return err
	}

	typechecker := &loader.TypeChecker{}
	pars := &crd.Parser{
		Collector: collector,
		Checker:   typechecker,
	}

	for _, pkg := range packages {
		pars.NeedPackage(pkg)
	}

	metav1Pkg := crd.FindMetav1(packages)
	if metav1Pkg == nil {
		return fmt.Errorf("no objects in the roots, since nothing imported metav1")
	}

	kubeKinds := crd.FindKubeKinds(pars, metav1Pkg)
	if len(kubeKinds) == 0 {
		return fmt.Errorf("no objects in the roots")
	}

	gk := schema.GroupKind{
		Group: cfg.Group,
		Kind:  cfg.Kind,
	}

	pars.NeedCRDFor(gk, nil)

	output := os.Stdout
	if len(cfg.Output) > 0 {
		output, err = os.OpenFile(cfg.Output, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	}
	mw := &multiwriter{w: output}
	for k, schemata := range pars.FlattenedSchemata {
		Degenerate(mw, 1, "", k.Name, schemata, schemata)
	}

	return mw.Error()
}

type multiwriter struct {
	w   io.Writer
	err error
}

func (m *multiwriter) Write(p []byte) (int, error) {
	if m.err != nil {
		return 0, m.err
	}
	n, err := m.w.Write(p)
	if err != nil {
		m.err = err
	}
	return n, err
}

func (m *multiwriter) Error() error {
	return m.err
}

func Degenerate(w io.Writer, level int, jsonpath string, key string, parent, node apiext.JSONSchemaProps) {
	if jsonpath == ".metadata" {
		return
	}

	if node.Type == "array" {
		Degenerate(w, level, jsonpath+"[]", key, parent, *node.Items.Schema)
		return
	}

	var required bool
	for _, k := range parent.Required {
		if k == key {
			required = true
			break
		}
	}

	_, _ = io.WriteString(w, fmt.Sprintf("%s %s\n", strings.Repeat("#", level), key))
	_, _ = io.WriteString(w, "\n")

	if len(node.Description) > 0 {
		_, _ = io.WriteString(w, strings.TrimSpace(node.Description))
		_, _ = io.WriteString(w, "\n\n")
	}

	_, _ = io.WriteString(w, fmt.Sprintf("* JSONPath: `%s`\n", jsonpath))
	_, _ = io.WriteString(w, fmt.Sprintf("* Type: `%s`\n", node.Type))
	_, _ = io.WriteString(w, fmt.Sprintf("* Required: `%s`\n", strconv.FormatBool(required)))

	if len(node.Pattern) > 0 {
		_, _ = io.WriteString(w, fmt.Sprintf("* Pattern: `%s`\n", node.Pattern))
	}
	if node.Minimum != nil {
		_, _ = io.WriteString(w, fmt.Sprintf("* Minimum value: `%0.f`\n", *node.Minimum))
	}
	if node.Maximum != nil {
		_, _ = io.WriteString(w, fmt.Sprintf("* Minimum value: `%0.f`\n", *node.Maximum))
	}
	if len(node.Enum) > 0 {
		_, _ = io.WriteString(w, fmt.Sprintf("* Allowed values:\n"))
		for _, v := range node.Enum {
			s := ""
			err := json.Unmarshal(v.Raw, &s)
			if err != nil {
				s = string(v.Raw)
			}
			if len(s) > 0 {
				_, _ = io.WriteString(w, fmt.Sprintf("  * `%s`\n", s))
			} else {
				_, _ = io.WriteString(w, fmt.Sprintf("  * (empty)\n"))
			}
		}
	}
	_, _ = io.WriteString(w, "\n")

	for k, n := range node.Properties {
		Degenerate(w, level+1, jsonpath+"."+k, k, node, n)
	}
}
