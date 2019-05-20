package internal

import (
	"sort"
	"strings"
	"testing"
)

func sliceEqual(t *testing.T, a, b []string) {
	return
}

func TestSort(t *testing.T) {
	// Test cases from https://semver.org/.
	cases := [][]string{
		{"1.9.0", "1.10.0", "1.11.0"},
		{"1.0.0", "2.0.0", "2.1.0", "2.1.1"},
		{"1.0.0-alpha", "1.0.0"},
		{"1.0.0-alpha", "1.0.0-alpha.1", "1.0.0-alpha.beta", "1.0.0-beta", "1.0.0-beta.2", "1.0.0-beta.11", "1.0.0-rc.1", "1.0.0"},
	}

	for _, input := range cases {
		t.Run(strings.Join(input, ", "), func(t *testing.T) {
			output := input[:]
			sort.Sort(semverSort(output))

			length := len(input)
			for i := range input {
				if input[i] != output[length-i-1] {
					t.Fatalf("expected '%v' at index '%v' (%v)", input[i], length-i-1, strings.Join(output, ", "))
				}
			}
		})
	}
}
