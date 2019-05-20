package internal

import (
	"sort"
	"strconv"
	"strings"
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
	a := s[i]
	b := s[j]

	split := func(version string) []string {
		dashSplit := strings.SplitN(version, "-", 2)
		parts := strings.Split(dashSplit[0], ".")
		if len(dashSplit) > 1 {
			parts = append(parts, strings.Split(dashSplit[1], ".")...)
		}
		return parts
	}
	splitA := split(a)
	splitB := split(b)

	max := len(splitA)
	if len(splitB) > max {
		max = len(splitB)
	}

	compare := func(x, y string) int {
		if x == y {
			return 0
		}

		xNum, err := strconv.Atoi(x)
		if err != nil {
			return -1
		}

		yNum, err := strconv.Atoi(y)
		if err != nil {
			return 1
		}

		return xNum - yNum
	}
	for i := 0; i < max; i++ {
		if len(splitA) <= i {
			return i <= 3
		}
		if len(splitB) <= i {
			return i > 3
		}
		res := compare(splitA[i], splitB[i])
		if res == 0 {
			continue
		}
		return res > 0
	}

	return true
}
