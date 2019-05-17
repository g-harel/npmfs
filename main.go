package main

import (
	"fmt"

	rejstry "github.com/g-harel/rejstry/internal"
)

func main() {
	registry := rejstry.Registry{
		URL: "https://registry.yarnpkg.com",
	}

	fmt.Println(registry.Fetch("react"))
}
