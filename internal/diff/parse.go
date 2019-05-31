package diff

import (
	"regexp"
	"strconv"
	"strings"
)

type Patch struct {
	PathA string
	PathB string
	Lines []PatchLine
}

type PatchLine struct {
	NumberA int
	NumberB int
	Content string
}

func parse(out string) ([]*Patch, error) {
	hunkPattern := regexp.MustCompile(`@@ -(\d+),\d+ \+(\d+),\d+ @@.*`)

	patches := []*Patch{}
	lines := strings.Split(out, "\n")

	patch := &Patch{}
	lineA := 0
	lineB := 0
	for _, line := range lines {
		// Detect start of new patch.
		if strings.HasPrefix(line, "diff --git ") {
			patch = &Patch{}
			patches = append(patches, patch)
			continue
		}

		// Detect file name in "a".
		if strings.HasPrefix(line, "--- a/") {
			patch.PathA = strings.TrimPrefix(line, "--- a/content/")
			continue
		}

		// Detect file name in "b".
		if strings.HasPrefix(line, "+++ b/") {
			patch.PathB = strings.TrimPrefix(line, "+++ b/content/")
			continue
		}

		// Detect start of hunk.
		if strings.HasPrefix(line, "@@ ") {
			match := hunkPattern.FindStringSubmatch(line)
			lineA, _ = strconv.Atoi(match[1])
			lineB, _ = strconv.Atoi(match[2])

			if len(patch.Lines) != 0 {
				patch.Lines = append(patch.Lines, PatchLine{0, 0, "@@"})
			}
			continue
		}

		// Detect added lines.
		if strings.HasPrefix(line, "+") {
			patch.Lines = append(patch.Lines, PatchLine{0, lineB, strings.TrimPrefix(line, "+")})
			lineB++
			continue
		}

		// Detect deleted lines.
		if strings.HasPrefix(line, "-") {
			patch.Lines = append(patch.Lines, PatchLine{lineA, 0, strings.TrimPrefix(line, "-")})
			lineA++
			continue
		}

		// Detect unchanged lines.
		if strings.HasPrefix(line, " ") {
			patch.Lines = append(patch.Lines, PatchLine{lineA, lineB, strings.TrimPrefix(line, " ")})
			lineA++
			lineB++
			continue
		}

		// Unrecognized lines are ignored.
	}

	return patches, nil
}
