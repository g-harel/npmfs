package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/g-harel/rejstry/internal"
)

func handler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")

	if len(parts) < 4 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(http.StatusText(http.StatusNotFound)))
		return
	}

	if len(parts) == 4 {
		versions, err := internal.PackageVersions(parts[1], parts[2])
		if err != nil {
			fmt.Println(err)
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "%s", strings.Join(versions, "\n"))
		return
	}

	dir, err := ioutil.TempDir("", "")
	if err != nil {
		panic(err)
	}

	contents, err := internal.PackageContents(parts[1], parts[2], parts[3])
	if err != nil {
		panic(err)
	}
	defer contents.Close()

	err = internal.Extract(contents, internal.Downloader(dir))
	if err != nil {
		panic(err)
	}

	dir = path.Join(dir, "package")

	r.URL.Path = "/" + strings.Join(parts[4:], "/")
	http.FileServer(http.Dir(dir)).ServeHTTP(w, r)
}

func main() {
	http.HandleFunc("/", handler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
