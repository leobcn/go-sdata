package main

import (
	"fmt"
	"github.com/jawher/mow.cli"
	"github.com/zpatrick/go-parser"
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
	PrimaryKey string
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
			fmt.Printf("[ERROR] %s\n", err.Error())
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
		return fmt.Errorf("No struct named \"%s\" found in input file", *specs.Struct)
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
				return nil, fmt.Errorf("The primary key field must be of type \"string\"")
			}

			fields = append(fields, field)
		}
	}

	if len(fields) == 1 {
		return fields[0], nil
	}

	rawField := "`data:\"primary_key\"`"
	if len(fields) == 0 {
		return nil, fmt.Errorf("No fields tagged with %s in target struct", rawField)
	}

	return nil, fmt.Errorf("Mutiple fields tagged with %s in target struct", rawField)
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

	pkg := target.File.Package
	if *specs.Package != "" {
		pkg = *specs.Package
	}

	context := TemplateContext{
		Package:    pkg,
		Type:       target.Name,
		PrimaryKey: primaryKey,
	}

	if err := tmpl.Execute(outputFile, context); err != nil {
		return err
	}

	return outputFile.Close()
}

const DefaultStoreTemplate = `package {{ .Package }}

import (
    "encoding/json"
    "github.com/zpatrick/go-sdata/container"
)

type {{ .Type }}Store struct {
    container container.Container
    table     string
}

func New{{ .Type }}Store(container container.Container) *{{ .Type }}Store {
    return &{{ .Type }}Store{
        container: container,
        table: "{{ .Package }}_{{ .Type }}",
    }
}

func (this *{{ .Type }}Store) Init() error {
    return this.container.Init(this.table)
}

type {{ .Type }}StoreCreate struct {
    *{{ .Type }}Store
    data *{{ .Type }}
}

func (this *{{ .Type }}Store) Create(data *{{ .Type }}) *{{ .Type }}StoreCreate {
    return &{{ .Type }}StoreCreate{
        {{ .Type }}Store: this,
        data:       data,
    }
}

func (this *{{ .Type }}StoreCreate) Execute() error {
    bytes, err := json.Marshal(this.data)
    if err != nil {
        return err
    }

    return this.container.Insert(this.table, this.data.ID, bytes)
}

type {{ .Type }}StoreSelect struct {
    *{{ .Type }}Store
    query  string
    filter {{ .Type }}Filter
    all    bool
}

func (this *{{ .Type }}Store) Select(query string) *{{ .Type }}StoreSelect {
    return &{{ .Type }}StoreSelect{
        {{ .Type }}Store: this,
        query:      query,
    }
}

func (this *{{ .Type }}Store) SelectAll() *{{ .Type }}StoreSelect {
    return &{{ .Type }}StoreSelect{
        {{ .Type }}Store: this,
        all:        true,
    }
}

type {{ .Type }}Filter func(*{{ .Type }}) bool

func (this *{{ .Type }}StoreSelect) Where(filter {{ .Type }}Filter) *{{ .Type }}StoreSelect {
    this.filter = filter
    return this
}

func (this *{{ .Type }}StoreSelect) Execute() ([]*{{ .Type }}, error) {
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

		if this.filter != nil && this.filter(value) {
			results = append(results, value)
		}
    }

    return results, nil
}

type {{ .Type }}StoreSelectFirst struct {
    *{{ .Type }}StoreSelect
}

func (this *{{ .Type }}StoreSelect) FirstOrNil() *{{ .Type }}StoreSelectFirst {
    return &{{ .Type }}StoreSelectFirst{
        {{ .Type }}StoreSelect: this,
    }
}

func (this *{{ .Type }}StoreSelectFirst) Execute() (*{{ .Type }}, error) {
    results, err := this.{{ .Type }}StoreSelect.Execute()
    if err != nil {
        return nil, err
    }

    if len(results) > 0 {
        return results[0], nil
    }

    return nil, nil
}

type {{ .Type }}StoreDelete struct {
    *{{ .Type }}Store
    key string
}

func (this *{{ .Type }}Store) Delete(key string) *{{ .Type }}StoreDelete {
    return &{{ .Type }}StoreDelete{
        {{ .Type }}Store: this,
        key:        key,
    }
}

func (this *{{ .Type }}StoreDelete) Execute() error {
    return this.container.Delete(this.table, this.key)
}
`
