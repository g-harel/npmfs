package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/api/v1/file", v1File)
	http.HandleFunc("/api/v1/files", v1Files)
	http.HandleFunc("/api/v1/versions", v1Versions)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("accepting connections at :%v", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
