package main

import (
	"fmt"

	"github.com/g-harel/rejstry/internal"
)

func main() {
	dir, err := internal.DownloadPackage("https://registry.npmjs.com", "react", "16.8.6")
	if err != nil {
		panic(err)
	}

	fmt.Println(dir)
}
