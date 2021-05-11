package main

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/doc"
	"go/parser"
	"go/token"
	"io"
	"os"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	apiext "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/code-generator/third_party/forked/golang/reflect"
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
	pflag.StringVar(&cfg.Directory, "dir", cfg.Directory, "directory with package")
	pflag.StringVar(&cfg.BaseClass, "baseclass", cfg.BaseClass, "class to base doc on")
	pflag.StringVar(&cfg.Group, "group", cfg.Group, "")
	pflag.StringVar(&cfg.Kind, "kind", cfg.Kind, "")
	pflag.StringVar(&cfg.Output, "output", cfg.Output, "output file")
	pflag.Parse()

	return foo(cfg)

	fileSet := token.NewFileSet()
	packages, err := parser.ParseDir(fileSet, cfg.Directory, nil, parser.ParseComments|parser.AllErrors)

	if err != nil {
		return err
	}

	files := make([]*ast.File, 0)
	for _, pkg := range packages {
		for _, file := range pkg.Files {
			files = append(files, file)
		}
	}

	pkg, err := doc.NewFromFiles(fileSet, files, "")
	if err != nil {
		return err
	}

	output := os.Stdout
	if len(cfg.Output) > 0 {
		output, err = os.OpenFile(cfg.Output, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	}

	nodes := make(map[string]*Node, 0)

	for _, typ := range pkg.Types {
		nodes[typ.Name] = &Node{
			Name:     typ.Name,
			Type:     typ.Name,
			Doc:      typ.Doc,
			Markers:  make([]string, 0),
			Children: make([]*Node, 0),
		}
	}

	for _, typ := range pkg.Types {
		node := nodes[typ.Name]
		for _, spec := range typ.Decl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}
			structType, ok := typeSpec.Type.(*ast.StructType)
			if !ok {
				continue
			}
			for _, field := range structType.Fields.List {
				for _, ident := range field.Names {
					declField, ok := ident.Obj.Decl.(*ast.Field)
					if !ok {
						continue
					}
					child := &Node{}

					child.Name = ident.Name
					child.JsonTag = declField.Tag.Value
					child.Required = true

					st, err := reflect.ParseStructTags(strings.Trim(declField.Tag.Value, "`"))
					if err == nil {
						for _, tag := range st {
							if tag.Name == "json" {
								tokens := strings.Split(tag.Value, ",")
								child.JsonTag = tokens[0]
								for _, tok := range tokens {
									if tok == "omitempty" {
										child.Required = false
									}
								}
							}
						}
					} else {
						log.Errorf("error parsing json tag: %s", err)
					}

					doc := make([]string, 0)
					if declField.Doc != nil {
						for _, comment := range declField.Doc.List {
							doc = append(doc, comment.Text)
						}
					}

					child.Doc = formatDoc(doc)
					extractDocstrings(doc, child)

					switch t := declField.Type.(type) {
					case *ast.Ident:
						child.Type = t.Name
					case *ast.ArrayType:
						child.Array = true
						switch tx := t.Elt.(type) {
						case *ast.Ident:
							child.Type = tx.Name
						case *ast.SelectorExpr:
							child.Type = tx.X.(*ast.Ident).Name + "." + tx.Sel.Name
						}
					case *ast.StarExpr:
						switch tx := t.X.(type) {
						case *ast.Ident:
							child.Type = tx.Name
						case *ast.SelectorExpr:
							child.Type = tx.X.(*ast.Ident).Name + "." + tx.Sel.Name
						}
					default:
						child.Type = declField.Names[0].Name
					}

					node.Children = append(node.Children, child)
				}
			}
		}
	}

	Generate(output, 1, ".spec", nodes[cfg.BaseClass], nodes)

	return nil
}

func Generate(w io.Writer, level int, jsonpath string, node *Node, nodes map[string]*Node) error {
	if node == nil {
		return fmt.Errorf("unreachable class")
	}
	if level > 10 {
		return fmt.Errorf("too much")
	}

	if len(node.JsonTag) > 0 {
		jsonpath += "." + node.JsonTag
	}

	io.WriteString(w, fmt.Sprintf("%s %s\n", strings.Repeat("#", level), node.Name))
	io.WriteString(w, "\n")
	io.WriteString(w, fmt.Sprintf("* JSONPath: `%s`\n", jsonpath))
	io.WriteString(w, fmt.Sprintf("* Type: `%s`\n", node.Type))
	io.WriteString(w, fmt.Sprintf("* Required: `%s`\n", strconv.FormatBool(node.Required)))
	for _, marker := range node.Markers {
		io.WriteString(w, fmt.Sprintf("* `%s`\n", marker))
	}
	io.WriteString(w, "\n")

	if len(node.Doc) > 0 {
		io.WriteString(w, strings.TrimSpace(node.Doc))
		io.WriteString(w, "\n\n")
	}

	if node.Array {
		jsonpath += "[]"
	}

	for _, child := range node.Children {
		err := Generate(w, level+1, jsonpath, child, nodes)
		if false && err != nil {
			return err
		}
	}

	if node.Children == nil {
		this, ok := nodes[node.Type]
		if !ok {
			return fmt.Errorf("unable to reach class %s", node.Type)
		}
		for _, child := range this.Children {
			err := Generate(w, level+1, jsonpath, child, nodes)
			if false && err != nil {
				return err
			}
		}
	}

	return nil
}

func formatDoc(lines []string) string {
	doc := ""
	for i, line := range lines {
		lines[i] = strings.TrimSpace(strings.TrimLeft(line, "/"))
		if !strings.HasPrefix(lines[i], "+") {
			doc = doc + lines[i] + "\n"
		}
	}
	return strings.TrimSpace(doc)
}

func extractDocstrings(lines []string, node *Node) {
	mark := make([]string, 0)
	for _, line := range lines {
		line = strings.TrimSpace(strings.TrimLeft(line, "/"))
		if strings.HasPrefix(line, "+") {
			mark = append(mark, line)
		}
	}
	node.Markers = mark
}

func foo(cfg *Config) error {
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
