package gogetter

import (
	"embed"
	"fmt"
	"io"
	"text/template"
)

//go:embed *.go.tpl
var embeddedTemplates embed.FS

var templateFuncs = template.FuncMap{}

type tplData struct {
	Structs map[string]*StructInfo

	PackageName string
	Version     string
	Imports     []string

	TargetStructName string
	SourceStructName string
}

func executeTemplate(w io.Writer, data tplData) error {
	tpl := template.New("gogetter")
	tpl.Funcs(templateFuncs)

	tpl, err := tpl.ParseFS(embeddedTemplates, "*.go.tpl")
	if err != nil {
		return fmt.Errorf("parse template: %w", err)
	}

	if err := tpl.ExecuteTemplate(w, "template.go.tpl", data); err != nil {
		return fmt.Errorf("execute template: %w", err)
	}

	return nil
}
