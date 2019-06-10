package util_test

import (
	"sort"
	"testing"

	"github.com/g-harel/npmfs/internal/util"
)

func TestSemverSort(t *testing.T) {
	cases := map[string][]string{
		// https://semver.org/#semantic-versioning-specification-semver
		"single element":          {"1.11.0", "1.10.0", "1.9.0"},
		"element priority":        {"2.1.1", "2.1.0", "2.0.0", "1.0.0"},
		"pre-release":             {"1.0.0", "1.0.0-alpha"},
		"pre-release identifiers": {"1.0.0", "1.0.0-rc.1", "1.0.0-beta.11", "1.0.0-beta.2", "1.0.0-beta", "1.0.0-alpha.beta", "1.0.0-alpha.1", "1.0.0-alpha"},
	}

	for name, input := range cases {
		t.Run(name, func(t *testing.T) {
			output := append([]string{}, input...)
			sort.Sort(util.SemverSort(output))

			for i := range input {
				if input[i] != output[i] {
					t.Fatalf("expected/received do not match\n%v\n%v", input, output)
				}
			}
		})
	}
}
