package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"text/template"

	yaml2 "github.com/ghodss/yaml"
	"github.com/imdario/mergo"
	nais_io_v1 "github.com/nais/liberator/pkg/apis/nais.io/v1"
	nais_io_v1alpha1 "github.com/nais/liberator/pkg/apis/nais.io/v1alpha1"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"gopkg.in/yaml.v3"
	apiext "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-tools/pkg/crd"
	crd_markers "sigs.k8s.io/controller-tools/pkg/crd/markers"
	"sigs.k8s.io/controller-tools/pkg/loader"
	"sigs.k8s.io/controller-tools/pkg/markers"
)

// Generate documentation for NAIS CRD's

var defaultResource interface{}

var exampleResource interface{}

type DefaultableResource interface {
	ApplyDefaults() error
}

type DocumentableResource struct {
	Resource      DefaultableResource
	ExampleGetter func() interface{}
}

var supportedResources = map[string]DocumentableResource{
	"Application": {
		Resource:      &nais_io_v1alpha1.Application{},
		ExampleGetter: func() interface{} { return nais_io_v1alpha1.ExampleApplicationForDocumentation() },
	},
	"Naisjob": {
		Resource:      &nais_io_v1.Naisjob{},
		ExampleGetter: func() interface{} { return nais_io_v1.ExampleNaisjobForDocumentation() },
	},
	"Alert": {
		Resource:      &nais_io_v1.Alert{},
		ExampleGetter: func() interface{} { return nais_io_v1.ExampleAlertForDocumentation() },
	},
}

type Renderer func(w io.Writer, level int, jsonpath string, key string, parent, node apiext.JSONSchemaProps)

type Config struct {
	Directory         string
	BaseClass         string
	Group             string
	Kind              string
	ReferenceOutput   string
	ExampleOutput     string
	ReferenceTemplate string
	ExampleTemplate   string
	JSONSchema        string
}

type Doc struct {
	// Which cluster(s) or environments the feature is available in
	Availability string `marker:"Availability,optional"`
	// Adds Default values to documentation
	Default string `marker:"Default,optional"`
	// Links to documentation or other information
	Link []string `marker:"Link,optional"`
	// Immutable declares the field as immutable
	Immutable bool `marker:"Immutable,optional"`
	// Deprecated declares the field obsolete
	Deprecated bool `marker:"Deprecated,optional"`
	// Experimental declares the field as subject to instability, change, or removal
	Experimental bool `marker:"Experimental,optional"`
}

type ExtDoc struct {
	Availability string
	Default      string
	Description  string
	Enum         []string
	Level        int
	Link         []string
	Maximum      *float64
	Minimum      *float64
	Path         string
	Pattern      string
	Required     bool
	Title        string
	Type         string
	Immutable    bool
	Deprecated   bool
	Experimental bool
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
	pflag.StringVar(&cfg.ReferenceOutput, "reference-output", cfg.ReferenceOutput, "reference doc markdown output file")
	pflag.StringVar(&cfg.ExampleOutput, "example-output", cfg.ExampleOutput, "example yaml markdown output file")
	pflag.StringVar(&cfg.ReferenceTemplate, "reference-template", cfg.ReferenceTemplate, "template file for rendering reference doc")
	pflag.StringVar(&cfg.ExampleTemplate, "example-template", cfg.ExampleTemplate, "template file for rendering example doc")
	pflag.StringVar(&cfg.JSONSchema, "json-schema-output", cfg.JSONSchema, "if set, generate json schema to the provided file")
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
		return fmt.Errorf("schema generation failed; double check the syntax of doctags (+nais:* and +kubebuilder:*")
	}

	resourcer, ok := supportedResources[cfg.Kind]
	if !ok {
		return fmt.Errorf("kind '%s' is not supported; needs populator config in docgen.go", cfg.Kind)
	}
	err = resourcer.Resource.ApplyDefaults()
	if err != nil {
		return err
	}
	err = marshalToInterface(&defaultResource, resourcer.Resource)
	if err != nil {
		return err
	}
	err = marshalToInterface(&exampleResource, resourcer.ExampleGetter())
	if err != nil {
		return err
	}

	if cfg.JSONSchema != "" && len(pars.FlattenedSchemata) > 1 {
		fmt.Fprintln(os.Stderr, "More than one schema, skipping the json schema")
		cfg.JSONSchema = ""
	}

	for _, schemata := range pars.FlattenedSchemata {
		err = Write(WriteReferenceDoc, cfg.ReferenceTemplate, cfg.ReferenceOutput, schemata.Properties["spec"])
		if err != nil {
			return err
		}
		err = Write(WriteExampleDoc, cfg.ExampleTemplate, cfg.ExampleOutput, schemata)
		if err != nil {
			return err
		}

		if cfg.JSONSchema != "" {
			if err := writeJSONSchema(cfg.JSONSchema, schemata); err != nil {
				return err
			}
		}
	}

	return nil
}

func writeJSONSchema(path string, schemata apiext.JSONSchemaProps) error {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	schemata.Schema = apiext.JSONSchemaURL("http://json-schema.org/draft-07/schema")

	// Make some changes to the schema to make it even more useful for validation etc.
	schemata = setJSONSchemaDefault(schemata, "apiVersion", `"nais.io/v1alpha1"`)
	schemata = setJSONSchemaDefault(schemata, "kind", `"Application"`)

	schemata = setJSONSchemaRequired(schemata, ".", "kind", "metadata")
	schemata = setJSONSchemaRequired(schemata, "metadata", "labels")
	schemata = setJSONSchemaRequired(schemata, "metadata.labels", "team")

	schemata = addJSONSchemaProperty(schemata, "metadata.labels.team", apiext.JSONSchemaProps{
		Description: "Team name",
		Type:        "string",
	})

	enc := json.NewEncoder(f)
	enc.SetIndent("", "\t")
	return enc.Encode(schemata)
}

func marshalToInterface(dst, src interface{}) error {
	data, err := json.Marshal(src)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, dst)
}

func Write(renderer Renderer, tpl string, outFile string, base apiext.JSONSchemaProps) error {
	var err error
	w := os.Stdout
	if len(outFile) > 0 {
		w, err = os.OpenFile(outFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			return err
		}
	}
	mw := &multiwriter{w: w}

	templateEngine, err := template.ParseFiles(tpl)
	if err != nil {
		return err
	}

	err = templateEngine.Execute(mw, nil)
	if err != nil {
		return err
	}

	renderer(mw, 1, "", "", base, base)

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

func floatfmt(f *float64) string {
	if f == nil {
		return "+Inf"
	}
	return strconv.FormatFloat(*f, 'f', 0, 64)
}

func writeList(w io.Writer, list []string) {
	sort.Strings(list)
	max := len(list) - 1
	for i, item := range list {
		if len(item) > 0 {
			io.WriteString(w, fmt.Sprintf("`%s`", item))
		} else {
			io.WriteString(w, fmt.Sprintf("_(empty string)_"))
		}
		if i != max {
			io.WriteString(w, ", ")
		}
	}
	io.WriteString(w, "<br />\n")
}

func (m ExtDoc) formatStraight(w io.Writer) {
	io.WriteString(w, fmt.Sprintf("%s %s", strings.Repeat("#", m.Level), strings.TrimLeft(m.Path, ".")))
	io.WriteString(w, "\n")
	if len(m.Description) > 0 {
		io.WriteString(w, m.Description)
		io.WriteString(w, "\n\n")
	}
	if m.Experimental {
		io.WriteString(w, "!!! warning \"Experimental feature\"\n    This feature has not undergone much testing, and is subject to API change, instability, or removal.\n\n")
	}
	if m.Deprecated {
		io.WriteString(w, "!!! failure \"Deprecated\"\n    This feature is deprecated, preserved only for backwards compatibility.\n\n")
	}
	if len(m.Link) > 0 {
		io.WriteString(w, "Relevant information:\n\n")
		for _, link := range m.Link {
			io.WriteString(w, fmt.Sprintf("* [%s](%s)\n", link, link))
		}
		io.WriteString(w, "\n")
	}
	io.WriteString(w, linefmt("Type: `%s`", m.Type))
	io.WriteString(w, linefmt("Required: `%s`", strconv.FormatBool(m.Required)))
	if m.Immutable {
		io.WriteString(w, linefmt("Immutable: `%v`", m.Immutable))
	}
	if len(m.Default) > 0 {
		io.WriteString(w, linefmt("Default value: `%v`", m.Default))
	}
	if len(m.Availability) > 0 {
		io.WriteString(w, linefmt("Availability: %s", m.Availability))
	}
	if len(m.Pattern) > 0 {
		io.WriteString(w, linefmt("Pattern: `%s`", m.Pattern))
	}
	if m.Minimum != m.Maximum {
		min := floatfmt(m.Minimum)
		max := floatfmt(m.Maximum)
		switch {
		case m.Minimum == nil:
			io.WriteString(w, linefmt("Maximum value: `%s`", max))
		case m.Maximum == nil:
			io.WriteString(w, linefmt("Minimum value: `%s`", min))
		default:
			io.WriteString(w, linefmt("Value range: `%s`-`%s`", min, max))
		}
	}
	if len(m.Enum) > 0 {
		io.WriteString(w, "Allowed values: ")
		writeList(w, m.Enum)
	}
	io.WriteString(w, "\n")
}

func hasRequired(node apiext.JSONSchemaProps, key string) bool {
	for _, k := range node.Required {
		if k == key {
			return true
		}
	}

	if node.Items == nil {
		return false
	}

	for _, k := range node.Items.Schema.Required {
		if k == key {
			return true
		}
	}

	return false
}

func WriteExampleDoc(w io.Writer, level int, jsonpath string, key string, parent, node apiext.JSONSchemaProps) {
	js, _ := json.Marshal(exampleResource)
	ym, _ := yaml2.JSONToYAML(js)

	io.WriteString(w, "``` yaml\n")
	io.WriteString(w, string(ym))
	io.WriteString(w, "```\n")

	return
}

func WriteReferenceDoc(w io.Writer, level int, jsonpath string, key string, parent, node apiext.JSONSchemaProps) {
	if jsonpath == ".metadata" || jsonpath == ".status" {
		return
	}

	if len(node.Enum) > 0 {
		node.Type = "enum"
	}

	entry := &ExtDoc{
		Description: strings.TrimSpace(node.Description),
		Level:       level,
		Maximum:     node.Maximum,
		Minimum:     node.Minimum,
		Path:        jsonpath,
		Pattern:     node.Pattern,
		Required:    hasRequired(parent, key),
		Title:       key,
		Type:        node.Type,
	}

	// Override children when encountering an array
	if node.Type == "array" {
		node.Properties = node.Items.Schema.Properties
		jsonpath += "[]"
	}

	defaultValue, err := getValueFromStruct("spec"+jsonpath, defaultResource)
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
			entry.Availability = d.Availability
			entry.Link = d.Link
			entry.Default = d.Default
			entry.Deprecated = d.Deprecated
			entry.Experimental = d.Experimental
			entry.Immutable = d.Immutable
		} else {
			log.Errorf("unable to merge structs: %s", err)
		}
	}

	if len(jsonpath) > 0 {
		entry.formatStraight(w)

		example, err := getStructSubPath("spec"+jsonpath, exampleResource)
		if err == nil {
			io.WriteString(w, "??? example\n")
			io.WriteString(w, "    ``` yaml\n")
			buf := bytes.NewBuffer(nil)
			enc := yaml.NewEncoder(buf)
			enc.SetIndent(2)
			enc.Encode(example)
			scan := bufio.NewScanner(buf)
			for scan.Scan() {
				io.WriteString(w, "    "+scan.Text()+"\n")
			}
			io.WriteString(w, "    ```\n\n")
		}
	}

	if len(node.Properties) == 0 {
		return
	}

	keys := make([]string, 0)
	for k := range node.Properties {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		WriteReferenceDoc(w, level+1, jsonpath+"."+k, k, node, node.Properties[k])
	}
}

func getStructSubPath(keyWithDots string, object interface{}) (interface{}, error) {
	structure := make(map[string]interface{})
	var leaf interface{} = structure

	keySlice := strings.Split(keyWithDots, ".")
	v := reflect.ValueOf(object)

	resolve := func(v reflect.Value) reflect.Value {
		if v.Kind() == reflect.Ptr {
			return v.Elem()
		}
		return v
	}

	max := len(keySlice) - 1
	for i, key := range keySlice {
		key = strings.TrimRight(key, "[]")

		if len(key) == 0 {
			break
		}

		v = resolve(v)

		var added interface{}

		switch v.Kind() {
		case reflect.Map:
			drilldown := func() error {
				for _, k := range v.MapKeys() {
					if k.String() == key {
						v = v.MapIndex(k).Elem()
						return nil
					}
				}
				return fmt.Errorf("key not found")
			}

			err := drilldown()
			if err != nil {
				return nil, err
			}
		}

		v = resolve(v)

		switch {
		case v.Kind() == reflect.Slice:
			fallthrough
		case i == max:
			added = resolve(v).Interface()

		case v.Kind() == reflect.Map:
			added = make(map[string]interface{})
		}

		switch typedleaf := leaf.(type) {
		case map[string]interface{}:
			typedleaf[key] = added
		case []interface{}:
			typedleaf[0] = added
		}

		leaf = added
		if v.Kind() == reflect.Slice {
			break
		}
	}

	return structure, nil
}

func getValueFromStruct(keyWithDots string, object interface{}) (interface{}, error) {
	keySlice := strings.Split(keyWithDots, ".")
	v := reflect.ValueOf(object)

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

func runOnJSONSchemaProperty(root apiext.JSONSchemaProps, path string, f func(*apiext.JSONSchemaProps)) apiext.JSONSchemaProps {
	if path == "." {
		f(&root)
		return root
	}

	p := strings.Split(path, ".")
	obj := root.Properties[p[0]]
	if len(p) == 1 {
		f(&obj)
	} else {
		runOnJSONSchemaProperty(obj, strings.Join(p[1:], "."), f)
	}
	root.Properties[p[0]] = obj
	return root
}

func setJSONSchemaDefault(root apiext.JSONSchemaProps, path string, value string) apiext.JSONSchemaProps {
	return runOnJSONSchemaProperty(root, path, func(obj *apiext.JSONSchemaProps) {
		obj.Default = &apiext.JSON{
			Raw: []byte(value),
		}
	})
}

func setJSONSchemaRequired(root apiext.JSONSchemaProps, path string, values ...string) apiext.JSONSchemaProps {
	return runOnJSONSchemaProperty(root, path, func(obj *apiext.JSONSchemaProps) {
		obj.Required = append(obj.Required, values...)
	})
}

func addJSONSchemaProperty(root apiext.JSONSchemaProps, path string, prop apiext.JSONSchemaProps) apiext.JSONSchemaProps {
	ps := strings.Split(path, ".")
	return runOnJSONSchemaProperty(root, strings.Join(ps[:len(ps)-1], "."), func(obj *apiext.JSONSchemaProps) {
		if obj.Properties == nil {
			obj.Properties = make(map[string]apiext.JSONSchemaProps)
		}
		obj.Properties[ps[len(ps)-1]] = prop
	})
}
