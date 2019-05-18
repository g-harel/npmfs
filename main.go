package main

import (
	"fmt"

	rejstry "github.com/g-harel/rejstry/internal"
)

func main() {
	registry := rejstry.Registry{
		// URL: "https://registry.yarnpkg.com",
		URL: "https://registry.npmjs.com",
	}

	pkg, err := registry.Fetch("react")
	if err != nil {
		panic(err)
	}

	dir, err := pkg.Download(pkg.Tags.Latest)
	if err != nil {
		panic(err)
	}

	fmt.Println(dir)
}
