package main

import (
	"fmt"
	"github.com/jawher/mow.cli"
	"github.com/zpatrick/go-parser"
	"log"
	"os"
	"strings"
	"text/template"
)

type Specs struct {
	Path     *string
	Struct   *string
	Output   *string
	Package  *string
	Template *string
}

type TemplateContext struct {
	Package    string
	Type       string
	Name       string
	PrimaryKey string
	TypeImport string
}

func main() {
	Specs := Specs{}
	app := cli.App("sdata", "Data persistance made simple")
	Specs.Path = app.StringArg("PATH", "", "Path to the source file")
	Specs.Struct = app.StringArg("STRUCT", "", "Name of the source struct")
	Specs.Output = app.StringOpt("o output", "stdout", "Path to the destination file")
	Specs.Package = app.StringOpt("p package", "", "Package name for the destination file")
	Specs.Template = app.StringOpt("t template", "", "Path to the template file")

	app.Spec = "PATH STRUCT [--output] [--package] [--template]"

	app.Action = func() {
		if err := Generate(Specs); err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
	}

	app.Run(os.Args)
}

func Generate(specs Specs) error {
	goFile, err := parser.ParseFile(*specs.Path)
	if err != nil {
		return err
	}

	if len(goFile.Structs) == 0 {
		return fmt.Errorf("No struct definitions found in source file")
	}

	if specs.Package == nil {
		specs.Package = &goFile.Package
	}

	var target *parser.GoStruct
	for _, goStruct := range goFile.Structs {
		if goStruct.Name == *specs.Struct {
			target = goStruct
			break
		}
	}

	if target == nil {
		return fmt.Errorf("No struct named '%s' found in input file", *specs.Struct)
	}

	primaryKey, err := findPrimaryKey(target)
	if err != nil {
		return err
	}

	return writeDataFile(specs, target, primaryKey.Name)
}

func findPrimaryKey(target *parser.GoStruct) (*parser.GoField, error) {
	fields := []*parser.GoField{}

	for _, field := range target.Fields {
		if field.Tag == nil {
			continue
		}

		// if tag := field.Tag.Get("data"); tag == "primary_key" {
		if strings.Contains(field.Tag.Value, "data:\"primary_key\"") {
			if field.Type != "string" {
				return nil, fmt.Errorf("The primary key field must be of type 'string'")
			}

			fields = append(fields, field)
		}
	}

	rawField := "`data:\"primary_key\"`"
	if len(fields) == 0 {
		return nil, fmt.Errorf("No fields tagged with %s in target struct", rawField)
	} else if len(fields) > 1 {
		return nil, fmt.Errorf("Mutiple fields tagged with %s in target struct", rawField)
	}

	return fields[0], nil
}

func writeDataFile(specs Specs, target *parser.GoStruct, primaryKey string) error {
	var parser func() (*template.Template, error)

	if *specs.Template == "" {
		parser = func() (*template.Template, error) { return template.New("").Parse(DefaultStoreTemplate) }
	} else {
		parser = func() (*template.Template, error) { return template.ParseFiles(*specs.Template) }
	}

	tmpl, err := parser()
	if err != nil {
		return err
	}

	outputFile := os.Stdout
	if *specs.Output != "stdout" {
		out, err := os.OpenFile(*specs.Output, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0666)
		if err != nil {
			return err
		}

		outputFile = out
	}

	outputPackage := target.File.Package
	if *specs.Package != "" {
		outputPackage = *specs.Package
	}

	var typeImport string

	structType := target.Name
	if outputPackage != target.File.Package {
		structType = fmt.Sprintf("%s.%s", target.File.Package, target.Name)

		path, err := target.File.ImportPath()
		if err != nil {
			return err
		}

		typeImport = path
		//imports = append(imports, target.File.Package)
	}

	// todo: if *specs.Package != to source struct package, type should be <source_struct>.<target.Name>,
	// and we should import that file
	context := TemplateContext{
		Package:    outputPackage,
		Type:       structType,
		Name:       target.Name,
		PrimaryKey: primaryKey,
		TypeImport: typeImport,
	}

	if err := tmpl.Execute(outputFile, context); err != nil {
		return err
	}

	return outputFile.Close()
}

const DefaultStoreTemplate = `package {{ .Package }}

// Automatically generated by go-sdata. DO NOT EDIT!

import (
    "encoding/json"
    "github.com/zpatrick/go-sdata/container"
    {{ if .TypeImport }} "{{ .TypeImport }}" {{ end }}
)

type {{ .Name }}Store struct {
    container container.Container
    table     string
}

func New{{ .Name }}Store(container container.Container) *{{ .Name }}Store {
    return &{{ .Name }}Store{
        container: container,
        table:     "{{ .Type }}",
    }
}

func (this *{{ .Name }}Store) Init() error {
    return this.container.Init(this.table)
}

type {{ .Name }}StoreInsert struct {
    *{{ .Name }}Store
    data *{{ .Type }}
}

func (this *{{ .Name }}Store) Insert(data *{{ .Type }}) *{{ .Name }}StoreInsert {
    return &{{ .Name }}StoreInsert{
        {{ .Name }}Store: this,
        data:       data,
    }
}

func (this *{{ .Name }}StoreInsert) Execute() error {
    bytes, err := json.Marshal(this.data)
    if err != nil {
        return err
    }

    return this.container.Insert(this.table, this.data.ID, bytes)
}

type {{ .Name }}StoreSelect struct {
    *{{ .Name }}Store
    query  string
    filter {{ .Name }}Filter
    all    bool
}

func (this *{{ .Name }}Store) Select(query string) *{{ .Name }}StoreSelect {
    return &{{ .Name }}StoreSelect{
        {{ .Name }}Store: this,
        query:      query,
    }
}

func (this *{{ .Name }}Store) SelectAll() *{{ .Name }}StoreSelect {
    return &{{ .Name }}StoreSelect{
        {{ .Name }}Store: this,
        all:        true,
    }
}

type {{ .Name }}Filter func(*{{ .Type }}) bool

func (this *{{ .Name }}StoreSelect) Where(filter {{ .Name }}Filter) *{{ .Name }}StoreSelect {
    this.filter = filter
    return this
}

func (this *{{ .Name }}StoreSelect) Execute() ([]*{{ .Type }}, error) {
    var query func() (map[string][]byte, error)

    if this.all {
        query = func() (map[string][]byte, error) { return this.container.SelectAll(this.table) }
    } else {
        query = func() (map[string][]byte, error) { return this.container.Select(this.table, this.query) }
    }

    data, err := query()
    if err != nil {
        return nil, err
    }

    results := []*{{ .Type }}{}
    for _, d := range data {
        var value *{{ .Type }}

        if err := json.Unmarshal(d, &value); err != nil {
            return nil, err
        }

		if this.filter == nil || this.filter(value) {
			results = append(results, value)
		}
    }

    return results, nil
}

type {{ .Name }}StoreSelectFirst struct {
    *{{ .Name }}StoreSelect
}

func (this *{{ .Name }}StoreSelect) FirstOrNil() *{{ .Name }}StoreSelectFirst {
    return &{{ .Name }}StoreSelectFirst{
        {{ .Name }}StoreSelect: this,
    }
}

func (this *{{ .Name }}StoreSelectFirst) Execute() (*{{ .Type }}, error) {
    results, err := this.{{ .Name }}StoreSelect.Execute()
    if err != nil {
        return nil, err
    }

    if len(results) > 0 {
        return results[0], nil
    }

    return nil, nil
}

type {{ .Name }}StoreDelete struct {
    *{{ .Name }}Store
    key string
}

func (this *{{ .Name }}Store) Delete(key string) *{{ .Name }}StoreDelete {
    return &{{ .Name }}StoreDelete{
        {{ .Name }}Store: this,
        key:        key,
    }
}

func (this *{{ .Name }}StoreDelete) Execute() (bool, error) {
    return this.container.Delete(this.table, this.key)
}
`
