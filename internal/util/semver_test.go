package util_test

import (
	"math/rand"
	"sort"
	"testing"

	"github.com/g-harel/npmfs/internal/util"
)

func TestSemverSort(t *testing.T) {
	//
	// https://semver.org/#semantic-versioning-specification-semver
	tt := map[string]struct {
		Expected []string
	}{
		"single element": {
			Expected: []string{"1.11.0", "1.10.0", "1.9.0"},
		},
		"element priority": {
			Expected: []string{"2.1.1", "2.1.0", "2.0.0", "1.0.0"},
		},
		"pre-release": {
			Expected: []string{"1.0.0", "1.0.0-alpha"},
		},
		"pre-release identifiers": {
			Expected: []string{"1.0.0", "1.0.0-rc.1", "1.0.0-beta.11", "1.0.0-beta.2", "1.0.0-beta", "1.0.0-alpha.beta", "1.0.0-alpha.1", "1.0.0-alpha"},
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			// Shuffle and re-sort test case.
			input := append([]string{}, tc.Expected...)
			rand.Shuffle(len(input), func(i, j int) {
				input[i], input[j] = input[j], input[i]
			})
			sort.Sort(util.SemverSort(input))

			sliceEqual(t, tc.Expected, input)
		})
	}
}
