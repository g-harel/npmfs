package util_test

import (
	"testing"

	"github.com/g-harel/npmfs/internal/util"
)

func TestUnique(t *testing.T) {
	tt := map[string]struct {
		Input    []string
		Expected []string
	}{
		"empty": {
			Input:    []string{},
			Expected: []string{},
		},
		"unique": {
			Input:    []string{"a", "a", "b", "c", "c", "c"},
			Expected: []string{"a", "b", "c"},
		},
		"sort": {
			Input:    []string{"z", "y", "x", "a", "b", "c"},
			Expected: []string{"a", "b", "c", "x", "y", "z"},
		},
		"sort_and_unique": {
			Input:    []string{"aa", "a", "ccc", "bb", "a", "ccc", "b"},
			Expected: []string{"a", "aa", "b", "bb", "ccc"},
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			sliceEqual(t, tc.Expected, util.Unique(tc.Input))
		})
	}
}
