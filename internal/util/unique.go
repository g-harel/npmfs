package util

import (
	"sort"
)

// Unique sorts and de-duplicates the input slice.
func Unique(s []string) []string {
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
