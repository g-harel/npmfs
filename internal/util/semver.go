package util

import (
	"sort"
	"strconv"
	"strings"
)

// SemverSort is a helper type to order semver version slices in decreasing order.
// All elements of the slice are assumed to be valid semver versions.
type SemverSort []string

var _ sort.Interface = SemverSort{}

func (s SemverSort) Len() int {
	return len(s)
}

func (s SemverSort) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s SemverSort) Less(i, j int) bool {
	splitA := semverSplit(s[i])
	splitB := semverSplit(s[j])

	for i := 0; ; i++ {
		// When major, minor, and patch are equal, a pre-release version has lower precedence than a normal version.
		// A larger set of pre-release fields has a higher precedence than a smaller set, if all of the preceding identifiers are equal.
		// Pre-release fields start at index i == 3.
		if len(splitA) <= i {
			return i <= 3
		}
		if len(splitB) <= i {
			return i > 3
		}

		res := semverIdentifierCompare(splitA[i], splitB[i])
		if res == 0 {
			// Equal identifiers means the next ones should be compared.
			continue
		}
		return res > 0
	}
}

// SemverSplit separates version strings into identifiers (1.2.3-alpha.1 => [1, 2, 3, alpha, 1]).
// First split around the first "-" for version/pre-release.
// Then split around "." for individual identifiers.
func semverSplit(version string) []string {
	dashSplit := strings.SplitN(version, "-", 2)
	parts := strings.Split(dashSplit[0], ".")
	if len(dashSplit) > 1 {
		parts = append(parts, strings.Split(dashSplit[1], ".")...)
	}
	return parts
}

// SemverIdentifierCompare compares two identifiers to determine precedence.
// Identifiers consisting of only digits are compared numerically.
// Identifiers with letters or hyphens are compared lexically in ASCII sort order
// Numeric identifiers always have lower precedence than non-numeric identifiers
func semverIdentifierCompare(a, b string) int {
	aNum, aStr := strconv.Atoi(a)
	bNum, bStr := strconv.Atoi(b)

	if aStr != nil && bStr != nil {
		return strings.Compare(a, b)
	}
	if aStr != nil {
		return 1
	}
	if bStr != nil {
		return -1
	}

	return aNum - bNum
}
