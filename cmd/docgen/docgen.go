package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"github.com/imdario/mergo"
	nais_io_v1alpha1 "github.com/nais/liberator/pkg/apis/nais.io/v1alpha1"
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

var defaultApplication interface{}

type Config struct {
	Directory string
	BaseClass string
	Group     string
	Kind      string
	Output    string
}

type Doc struct {
	// Sample string(s) that can be used in this field
	Sample []string `marker:"Sample,optional"`
	// Which cluster(s) or environments the feature is available in
	Availability string `marker:"Availability,optional"`
}

type ExtDoc struct {
	Doc
	Level       int
	Title       string
	Description string
	Path        string
	Type        string
	Required    bool
	Default     string
	Pattern     string
	Minimum     *float64
	Maximum     *float64
	Enum        []string
}

// Hijack the "example" field for custom documentation fields
func (m Doc) ApplyToSchema(schema *apiext.JSONSchemaProps) error {
	d := &Doc{}
	if schema.Example != nil {
		err := json.Unmarshal(schema.Example.Raw, d)
		if err != nil {
			return err
		}
	}
	err := mergo.Merge(d, m)
	if err != nil {
		return err
	}
	b, err := json.Marshal(d)
	if err != nil {
		return err
	}
	schema.Example = &apiext.JSON{Raw: b}
	return nil
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

	err = registry.Define("nais:doc", markers.DescribesField, Doc{})
	if err != nil {
		return fmt.Errorf("register marker: %w", err)
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

	if len(pars.FlattenedSchemata) == 0 {
		return fmt.Errorf("no schemas to print")
	}

	output := os.Stdout
	if len(cfg.Output) > 0 {
		output, err = os.OpenFile(cfg.Output, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	}
	mw := &multiwriter{w: output}

	app := nais_io_v1alpha1.Application{}
	err = nais_io_v1alpha1.ApplyDefaults(&app)
	if err != nil {
		return err
	}
	data, err := json.Marshal(app.Spec)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &defaultApplication)
	if err != nil {
		return err
	}

	for _, schemata := range pars.FlattenedSchemata {
		Degenerate(mw, 1, "", "NAIS application", schemata.Properties["spec"], schemata.Properties["spec"])
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

func linefmt(format string, args ...interface{}) string {
	format = fmt.Sprintf(format, args...)
	if len(format) == 0 {
		format = "_no value_"
	}
	format = strings.ReplaceAll(format, "``", "_no value_")
	return format + "<br />\n"
}

func (m ExtDoc) formatStraight(w io.Writer) {
	io.WriteString(w, fmt.Sprintf("%s %s\n", strings.Repeat("#", m.Level), m.Title))
	io.WriteString(w, "\n")
	if len(m.Description) > 0 {
		io.WriteString(w, m.Description)
		io.WriteString(w, "\n\n")
	}
	io.WriteString(w, linefmt("Path: `%s`", m.Path))
	io.WriteString(w, linefmt("Type: `%s`", m.Type))
	io.WriteString(w, linefmt("Required: `%s`", strconv.FormatBool(m.Required)))
	switch {
	case len(m.Default) > 0:
		io.WriteString(w, linefmt("Default value: `%v`", m.Default))
	case len(m.Sample) > 1:
		io.WriteString(w, linefmt("Example values:"))
		io.WriteString(w, fmt.Sprintf("\n"))
		for _, sample := range m.Sample {
			io.WriteString(w, fmt.Sprintf("  * `%s`\n", sample))
		}
	case len(m.Sample) == 1:
		io.WriteString(w, linefmt("Example value: `%s`", m.Sample[0]))
	}
	if len(m.Availability) > 0 {
		io.WriteString(w, linefmt("Availability: %s", m.Availability))
	}
	if len(m.Pattern) > 0 {
		io.WriteString(w, linefmt("Pattern: `%s`", m.Pattern))
	}
	if m.Minimum != m.Maximum {
		io.WriteString(w, linefmt("Minimum value: `%d`", m.Minimum))
		io.WriteString(w, linefmt("Minimum value: `%d`", m.Maximum))
	}
	if len(m.Enum) > 0 {
		io.WriteString(w, linefmt("Allowed values:"))
		for _, v := range m.Enum {
			if len(v) > 0 {
				io.WriteString(w, fmt.Sprintf("  * `%s`\n", v))
			} else {
				io.WriteString(w, fmt.Sprintf("  * _no value_\n"))
			}
		}
	}
	io.WriteString(w, "\n")
}

func cell(s string) string {
	return "| " + s + " "
}

func (m ExtDoc) formatTable(w io.Writer) {
	io.WriteString(w, fmt.Sprintf("%s %s\n", strings.Repeat("#", m.Level), m.Title))
	io.WriteString(w, "\n")
	if len(m.Description) > 0 {
		io.WriteString(w, m.Description)
		io.WriteString(w, "\n\n")
	}
	io.WriteString(w, "| Path | Type | Required | Default | Example | Availability | Supported values |\n")
	io.WriteString(w, "|------|------|----------|---------|---------|--------------|------------------|\n")
	io.WriteString(w, cell(m.Path))
	io.WriteString(w, cell(m.Type))
	io.WriteString(w, cell(strconv.FormatBool(m.Required)))
	io.WriteString(w, cell(m.Default))
	if len(m.Sample) > 0 {
		io.WriteString(w, cell(m.Sample[0]))
	} else {
		io.WriteString(w, cell(""))
	}
	io.WriteString(w, cell(m.Availability))
	switch {
	case m.Minimum != m.Maximum:
		io.WriteString(w, cell(fmt.Sprintf("%d-%d", m.Minimum, m.Maximum)))
	case len(m.Enum) > 0:
		io.WriteString(w, cell(strings.Join(m.Enum, ", ")))
	case len(m.Pattern) > 0:
		io.WriteString(w, cell(m.Pattern))
	}
	io.WriteString(w, "\n")
}

func hasRequired(node apiext.JSONSchemaProps, key string) bool {
	for _, k := range node.Required {
		if k == key {
			return true
		}
	}
	return false
}

func Degenerate(w io.Writer, level int, jsonpath string, key string, parent, node apiext.JSONSchemaProps) {
	if jsonpath == ".metadata" || jsonpath == ".status" {
		return
	}

	// Override children when encountering an array
	if node.Type == "array" {
		node.Properties = node.Items.Schema.Properties
		jsonpath += "[]"
	}

	if len(node.Enum) > 0 {
		node.Type = "enum"
	}

	entry := &ExtDoc{
		Title:       key,
		Path:        jsonpath,
		Required:    hasRequired(parent, key),
		Description: strings.TrimSpace(node.Description),
		Type:        node.Type,
	}

	defaultValue, err := getValueFromStruct(strings.Trim(jsonpath, "."), defaultApplication)
	if err == nil {
		entry.Default = fmt.Sprintf("%v", defaultValue)
	}

	if len(node.Enum) > 0 {
		entry.Enum = make([]string, 0, len(entry.Enum))
		for _, v := range node.Enum {
			s := ""
			err := json.Unmarshal(v.Raw, &s)
			if err != nil {
				s = string(v.Raw)
			}
			entry.Enum = append(entry.Enum, s)
		}
	}

	if node.Example != nil {
		d := &Doc{}
		err := json.Unmarshal(node.Example.Raw, d)
		if err == nil {
			err = mergo.Merge(entry, d)
		}
		if err == nil {
			log.Errorf("unable to merge structs: %s", err)
		}
	}

	entry.formatStraight(w)

	if len(node.Properties) == 0 {
		return
	}

	keys := make([]string, 0)
	for k := range node.Properties {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		Degenerate(w, level+1, jsonpath+"."+k, k, node, node.Properties[k])
	}
}

func getValueFromStruct(keyWithDots string, object interface{}) (interface{}, error) {
	keySlice := strings.Split(keyWithDots, ".")
	v := reflect.ValueOf(object)
	// iterate through field names ,ignore the first name as it might be the current instance name
	// you can make it recursive also if want to support types like slice,map etc along with struct

	for _, key := range keySlice {
		if len(key) == 0 {
			break
		}
		for v.Kind() == reflect.Ptr {
			v = v.Elem()
		}
		if v.Kind() != reflect.Map {
			return nil, fmt.Errorf("only accepts maps; got %T", v)
		}
		getKey := func() error {
			for _, k := range v.MapKeys() {
				if k.String() == key {
					v = v.MapIndex(k).Elem()
					return nil
				}
			}
			return fmt.Errorf("key not found")
		}
		err := getKey()
		if err != nil {
			return nil, err
		}
	}

	if !v.IsValid() {
		return nil, fmt.Errorf("no value")
	}

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.String:
	case reflect.Int:
	case reflect.Bool:
	case reflect.Float64:
	default:
		return nil, fmt.Errorf("only scalar values supported")
	}

	return v.Interface(), nil
}
