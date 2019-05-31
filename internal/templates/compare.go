package templates

import (
	"html/template"
	"log"
	"net/http"
)

type compareLine struct {
	NumberA int
	NumberB int
	Content string
}

type compareDiff struct {
	PathA string
	PathB string
	Lines []compareLine
}

func Compare(w http.ResponseWriter, r *http.Request, name, versionA, versionB string) {
	tmpl, err := template.ParseFiles(
		"templates/layout.html",
		"templates/pages/compare.html",
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("ERROR parse template: %v", err)
		return
	}

	println(name, versionA, versionB)

	context := &struct {
		Package  string
		VersionA string
		VersionB string
		Diffs    []compareDiff
	}{
		Package:  "react",
		VersionA: "0.0.0",
		VersionB: "1.0.0",
		Diffs: []compareDiff{
			{"/test/path", "/test/path/renamed", []compareLine{
				{1, 1, "test"},
				{2, 0, "testA"},
				{0, 2, "testB"},
				{3, 3, "test"},
			}},
			{"/file", "/file", []compareLine{
				{1, 1, "a"},
				{2, 0, "a"},
				{3, 0, "a"},
				{4, 2, "shared"},
				{0, 0, "..."},
				{0, 3, "b"},
				{0, 4, "b"},
				{0, 5, "b"},
			}},
		},
	}

	err = tmpl.Execute(w, context)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("ERROR execute template: %v", err)
		return
	}
}
