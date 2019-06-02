package paths_test

import (
	"testing"

	"github.com/g-harel/npmfs/internal/paths"
)

func TestBreakRelative(t *testing.T) {
	cases := map[string]struct {
		Parts []string
		Links []string
	}{
		"": {
			Parts: []string{},
			Links: []string{},
		},
		"/": {
			Parts: []string{},
			Links: []string{},
		},
		"img.jpg": {
			Parts: []string{"img.jpg"},
			Links: []string{""},
		},
		"/img.jpg": {
			Parts: []string{"img.jpg"},
			Links: []string{""},
		},
		"test/": {
			Parts: []string{"test"},
			Links: []string{""},
		},
		"/test/": {
			Parts: []string{"test"},
			Links: []string{""},
		},
		"test/img.jpg": {
			Parts: []string{"test", "img.jpg"},
			Links: []string{"./", ""},
		},
		"/test/img.jpg": {
			Parts: []string{"test", "img.jpg"},
			Links: []string{"./", ""},
		},
		"test/path/": {
			Parts: []string{"test", "path"},
			Links: []string{"../", ""},
		},
		"/test/path/": {
			Parts: []string{"test", "path"},
			Links: []string{"../", ""},
		},
		"test/path/img.jpg": {
			Parts: []string{"test", "path", "img.jpg"},
			Links: []string{"../", "./", ""},
		},
		"/test/path/img.jpg": {
			Parts: []string{"test", "path", "img.jpg"},
			Links: []string{"../", "./", ""},
		},
	}

	for path, result := range cases {
		parts, links := paths.BreakRelative(path)

		if len(parts) != len(result.Parts) {
			t.Fatalf("expected/received parts do not match \n%v\n%v\n%v", path, parts, result.Parts)
		}
		for i := range parts {
			if parts[i] != result.Parts[i] {
				t.Fatalf("expected/received parts do not match \n%v\n%v\n%v", path, parts, result.Parts)
			}
		}

		if len(links) != len(result.Links) {
			t.Fatalf("expected/received links do not match \n%v\n%v\n%v", path, links, result.Links)
		}
		for i := range links {
			if links[i] != result.Links[i] {
				t.Fatalf("expected/received links do not match \n%v\n%v\n%v", path, links, result.Links)
			}
		}
	}
}
