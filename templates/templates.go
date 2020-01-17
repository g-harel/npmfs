package templates

import (
	"html/template"
	"log"
	"net/http"
)

// Renderer is used to store state before a template is executed.
type Renderer struct {
	statusCode int
	filenames []string
	context   interface{}
}

// Handler executes the templates and handles errors.
// Does not attempt to render the error page template to avoid possible infinite recursion.
func (r *Renderer) Handler(w http.ResponseWriter, _ *http.Request) {
	if r.statusCode > 0 {
		w.WriteHeader(r.statusCode)
	}

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
