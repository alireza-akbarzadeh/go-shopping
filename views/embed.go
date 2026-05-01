package views

import (
	"embed"
	"html/template"
	"net/http"
)

//go:embed templates/*.html
var TemplatesFS embed.FS

//go:embed static
var StaticFS embed.FS

// RenderTemplate parses and executes a template from the embedded files.
func RenderTemplate(w http.ResponseWriter, name string, data interface{}) error {
	tmpl, err := template.ParseFS(TemplatesFS, "templates/"+name)
	if err != nil {
		return err
	}
	return tmpl.Execute(w, data)
}
