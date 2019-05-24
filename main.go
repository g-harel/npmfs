package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/g-harel/rejstry/internal/middleware"
	"github.com/g-harel/rejstry/internal/page"
)

func main() {
	// Rendered templates.
	http.Handle("/package/", middleware.Log(page.Package()))

	// Static files.
	static := http.FileServer(http.Dir("static"))
	http.Handle("/static/", middleware.Log(http.StripPrefix("/static/", static)))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("accepting connections at :%v", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
