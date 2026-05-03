// Package views provides functions for rendering HTML templates and serving static files using Go's embed package. It includes embedded file systems for templates and static assets, and a RenderTemplate function that parses and executes templates with provided data. This package abstracts away the details of template rendering and allows for easy integration with HTTP handlers in the application.
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
