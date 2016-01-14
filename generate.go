package main

import (
	"fmt"
	"github.com/jawher/mow.cli"
	"github.com/zpatrick/go-parser"
	"os"
	"text/template"
)

/*
	TODO:
	- use byte string variable for default template: https://github.com/codegangsta/cli/blob/cffab77ecb4f963ced9e30344eb2b9282ef36887/help.go#L12
	- rename a bunch of stuff. <type>Data sucks - it represeings the business logic layer. The store layers is properly named. What woudl i want?
		- Options: UserData.Select(), UserLogic, UserBL, data access layer(DAL) UserDAL, object relational mapping , UserStore.Select(), User
		- maybe use a store metaphor? Warehouse, retail, outlet, factory, supplier
	- give json store
*/

type Specs struct {
	Input    *string
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
	Specs.Input = app.StringOpt("i input", "", "Path to the source file")
	Specs.Struct = app.StringOpt("s struct", "", "Name of the source struct")
	Specs.Output = app.StringOpt("o output", "stdout", "Path to the destination file")
	Specs.Package = app.StringOpt("p package", "", "Package name for the destination file")
	Specs.Template = app.StringOpt("t template", "data.template", "Path to the template file")

	app.Spec = "--input --struct [--output] [--package] [--template]"

	app.Action = func() {
		if err := Generate(Specs); err != nil {
			fmt.Printf("[ERROR] %s\n", err.Error())
			os.Exit(1)
		}
	}

	app.Run(os.Args)
}

func Generate(specs Specs) error {
	goFile, err := parser.ParseFile(*specs.Input)
	if err != nil {
		return err
	}

	if len(goFile.Structs) == 0 {
		return fmt.Errorf("No struct definitions found in source file")
	}

	if specs.Package == nil {
		*specs.Package = goFile.Package
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
	var fields []*parser.GoField

	for _, field := range target.Fields {
		if tag := field.Tag.Get("sdata"); tag == "primary_key" {
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
	// todo: template := template.Must(template.New("help").Funcs(funcMap).Parse(templ))
	// or template := template.ParseFiles(templateFile)
	// both juse use template.Execute
	templateFile := "data.template"
	if *specs.Template != "" {
		templateFile = *specs.Template
	}

	outputFile := os.Stdout
	if *specs.Output != "stdout" {
		out, err := os.OpenFile(*specs.Output, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0666)
		if err != nil {
			return err
		}

		outputFile = out
	}

	template, err := template.ParseFiles(templateFile)
	if err != nil {
		return err
	}

	context := TemplateContext{
		Package:    target.File.Package,
		Type:       target.Name,
		PrimaryKey: primaryKey,
	}

	if err := template.Execute(outputFile, context); err != nil {
		return err
	}

	return outputFile.Close()
}
