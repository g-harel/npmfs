package internal

import (
	"sort"

	"github.com/Masterminds/semver"
)

type semverSort []string

var _ sort.Interface = semverSort{}

func (s semverSort) Len() int {
	return len(s)
}

func (s semverSort) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s semverSort) Less(i, j int) bool {
	return semver.MustParse(s[i]).Compare(semver.MustParse(s[j])) > 0
}
