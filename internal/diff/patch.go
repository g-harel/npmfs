package diff

import (
	"regexp"
	"strconv"
	"strings"
)

// Patch represents a single file diff.
// File was created/deleted when one of "PathA" or "PathB" is empty.
// File was renamed when "PathA" and "PathB" are not equal.
type Patch struct {
	PathA string
	PathB string
	Lines []PatchLine
}

// PatchLine represents a single line of diff output.
// Line was created/deleted when one of "LineA" or "LineB" is zero.
// Hunks starts when both "LineA" and "LineB" are zero.
type PatchLine struct {
	LineA   int
	LineB   int
	Content string
}

// PatchParse parses standard diff output.
func patchParse(out string) ([]*Patch, error) {
	hunkPattern := regexp.MustCompile(`@@ -(\d+)(?:,\d+)? \+(\d+)(?:,\d+)? @@.*`)

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
		if strings.HasPrefix(line, "---") {
			if strings.HasPrefix(line, "--- a/content/") {
				patch.PathA = strings.TrimPrefix(line, "--- a/content/")
			}
			continue
		}

		// Detect file name in "b".
		if strings.HasPrefix(line, "+++") {
			if strings.HasPrefix(line, "+++ b/content/") {
				patch.PathB = strings.TrimPrefix(line, "+++ b/content/")
			}
			continue
		}

		// Detect start of hunk.
		if strings.HasPrefix(line, "@@ ") {
			match := hunkPattern.FindStringSubmatch(line)
			if len(match) < 3 {
				continue
			}

			lineA, _ = strconv.Atoi(match[1])
			lineB, _ = strconv.Atoi(match[2])

			if len(patch.Lines) != 0 {
				patch.Lines = append(patch.Lines, PatchLine{0, 0, line})
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
