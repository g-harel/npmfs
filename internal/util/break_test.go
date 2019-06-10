package util_test

import (
	"testing"

	"github.com/g-harel/npmfs/internal/util"
)

func sliceEqual(t *testing.T, expected, received []string) {
	if len(expected) != len(received) {
		t.Fatalf("expected/received do not match\n%v\n%v", expected, received)
	}
	for i := range expected {
		if expected[i] != received[i] {
			t.Fatalf("expected/received do not match\n%v\n%v", expected, received)
		}
	}
}

func TestBreakPathRelative(t *testing.T) {
	tt := map[string]struct {
		Path  string
		Parts []string
		Links []string
	}{
		"root": {
			Path:  "",
			Parts: []string{},
			Links: []string{},
		},
		"root file": {
			Path:  "img.jpg",
			Parts: []string{"img.jpg"},
			Links: []string{""},
		},
		"single dir": {
			Path:  "test/",
			Parts: []string{"test"},
			Links: []string{""},
		},
		"nested file": {
			Path:  "test/img.jpg",
			Parts: []string{"test", "img.jpg"},
			Links: []string{"./", ""},
		},
		"nested dir": {
			Path:  "test/path/",
			Parts: []string{"test", "path"},
			Links: []string{"../", ""},
		},
		"deeply nested file": {
			Path:  "test/path/img.jpg",
			Parts: []string{"test", "path", "img.jpg"},
			Links: []string{"../", "./", ""},
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			parts, links := util.BreakPathRelative(tc.Path)

			t.Run("parts", func(t *testing.T) {
				sliceEqual(t, tc.Parts, parts)
			})

			t.Run("links", func(t *testing.T) {
				sliceEqual(t, tc.Links, links)
			})

			t.Run("leading slash ignored", func(t *testing.T) {
				parts, links := util.BreakPathRelative("/" + tc.Path)

				t.Run("parts", func(t *testing.T) {
					sliceEqual(t, tc.Parts, parts)
				})

				t.Run("links", func(t *testing.T) {
					sliceEqual(t, tc.Links, links)
				})
			})
		})
	}
}
