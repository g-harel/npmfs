package templates

import (
	"html/template"
	"log"
	"net/http"
)

type Renderer struct {
	filenames []string
	context   interface{}
}

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
