package templates

import (
	"sort"
	"strings"
)

// Sort and de-duplicate input slice.
func cleanup(s []string) []string {
	m := map[string]interface{}{}
	for _, item := range s {
		m[item] = true
	}
	out := []string{}
	for key := range m {
		out = append(out, key)
	}
	sort.Strings(out)
	return out
}

// Splits up the given path and calculates a relative link for each part.
func breakPath(path string) (parts []string, links []string) {
	// Remove leading slash.
	if len(path) > 0 && path[0] == '/' {
		path = path[1:]
	}

	// Empty paths are ignored.
	if len(path) == 0 {
		return
	}

	parts = strings.Split(path, "/")

	// Generate a decreasing series of relative links (ex. "../../", "../", "").
	links = []string{}
	for i := len(parts); i >= 0; i-- {
		links = append(links, strings.Repeat("../", i))
	}

	// Change behavior if the path represents a directory.
	if path[len(path)-1] == '/' {
		// Remove last path part, which will always be empty because of trailing slash.
		return parts[:len(parts)-1], links[2:]
	}

	// Add a "./" entry in the before-last position in the returned links.
	links = append(links[:len(links)-1], "./", links[len(links)-1])
	return parts, links[2:]
}

/*

TODO test
fmt.Println(breakPath("/"))
fmt.Println([]string{""}, []string{""})
println()
fmt.Println(breakPath("/test/"))
fmt.Println([]string{""}, []string{"test"})
println()
fmt.Println(breakPath("/test/path/"))
fmt.Println([]string{"../", ""}, []string{"test", "path"})
println()
fmt.Println(breakPath("/img.jpg"))
fmt.Println([]string{""}, []string{"img.jpg"})
println()
fmt.Println(breakPath("/test/img.jpg"))
fmt.Println([]string{"./", ""}, []string{"test", "img.jpg"})
println()
fmt.Println(breakPath("/test/path/img.jpg"))
fmt.Println([]string{"../", "./", ""}, []string{"test", "path", "img.jpg"})

*/
