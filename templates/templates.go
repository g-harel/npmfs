package templates

import (
	"html/template"
	"log"
	"net/http"
)

// Renderer is used to store state before a template is executed.
type Renderer struct {
	filenames []string
	context   interface{}
}

// Render is an application-aware helper to execute templates.
func (r *Renderer) Render(w http.ResponseWriter) {
	tmpl, err := template.ParseFiles(r.filenames...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("ERROR parse template: %v", err)
		return
	}

	err = tmpl.Execute(w, r.context)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("ERROR execute template: %v", err)
		return
	}
}
