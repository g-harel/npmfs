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
		Input         string
		ExpectedParts []string
		ExpectedLinks []string
	}{
		"root": {
			Input:         "",
			ExpectedParts: []string{},
			ExpectedLinks: []string{},
		},
		"root file": {
			Input:         "img.jpg",
			ExpectedParts: []string{"img.jpg"},
			ExpectedLinks: []string{""},
		},
		"single dir": {
			Input:         "test/",
			ExpectedParts: []string{"test"},
			ExpectedLinks: []string{""},
		},
		"nested file": {
			Input:         "test/img.jpg",
			ExpectedParts: []string{"test", "img.jpg"},
			ExpectedLinks: []string{"./", ""},
		},
		"nested dir": {
			Input:         "test/path/",
			ExpectedParts: []string{"test", "path"},
			ExpectedLinks: []string{"../", ""},
		},
		"deeply nested file": {
			Input:         "test/path/img.jpg",
			ExpectedParts: []string{"test", "path", "img.jpg"},
			ExpectedLinks: []string{"../", "./", ""},
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			parts, links := util.BreakPathRelative(tc.Input)

			t.Run("parts", func(t *testing.T) {
				sliceEqual(t, tc.ExpectedParts, parts)
			})

			t.Run("links", func(t *testing.T) {
				sliceEqual(t, tc.ExpectedLinks, links)
			})

			t.Run("leading slash ignored", func(t *testing.T) {
				parts, links := util.BreakPathRelative("/" + tc.Input)

				t.Run("parts", func(t *testing.T) {
					sliceEqual(t, tc.ExpectedParts, parts)
				})

				t.Run("links", func(t *testing.T) {
					sliceEqual(t, tc.ExpectedLinks, links)
				})
			})
		})
	}
}
