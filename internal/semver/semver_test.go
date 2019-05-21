package semver_test

import (
	"sort"
	"strings"
	"testing"

	"github.com/g-harel/rejstry/internal/semver"
)

func TestSort(t *testing.T) {
	cases := [][]string{
		// https://semver.org/
		{"1.11.0", "1.10.0", "1.9.0"},
		{"2.1.1", "2.1.0", "2.0.0", "1.0.0"},
		{"1.0.0", "1.0.0-alpha"},
		{"1.0.0", "1.0.0-rc.1", "1.0.0-beta.11", "1.0.0-beta.2", "1.0.0-beta", "1.0.0-alpha.beta", "1.0.0-alpha.1", "1.0.0-alpha"},
	}

	for _, input := range cases {
		t.Run(strings.Join(input, ", "), func(t *testing.T) {
			output := append([]string{}, input...)
			sort.Sort(semver.Sort(output))

			for i := range input {
				if input[i] != output[i] {
					t.Fatalf("expected/received do not match \n%v\n%v", input, output)
				}
			}
		})
	}
}
