package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/g-harel/rejstry/internal/middleware"
)

func main() {
	http.HandleFunc("/api/v1/diff", middleware.Log(v1Diff))
	http.HandleFunc("/api/v1/file", middleware.Log(v1File))
	http.HandleFunc("/api/v1/files", middleware.Log(v1Files))
	http.HandleFunc("/api/v1/versions", middleware.Log(v1Versions))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("accepting connections at :%v", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
