package prettydiff

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type FileDiff struct {
	OldFile string
	NewFile string
	Hunks   []Hunk
}

type Hunk struct {
	OldStart int
	OldCount int
	NewStart int
	NewCount int
	Context  string
	Changes  []Change
}

type ChangeType int

const (
	Unchanged ChangeType = iota
	Added
	Removed
)

type Change struct {
	Type    ChangeType
	Content string
	OldLine int
	NewLine int
}

func ParseUnifiedDiff(diffText string) (*FileDiff, error) {
	diff := &FileDiff{}
	lines := strings.Split(diffText, "\n")

	if len(lines) < 2 {
		return nil, fmt.Errorf("invalid diff format: not enough lines")
	}

	if !strings.HasPrefix(lines[0], "---") || !strings.HasPrefix(lines[1], "+++") {
		return nil, fmt.Errorf("invalid diff format: missing file headers")
	}

	diff.OldFile = strings.TrimPrefix(lines[0], "--- ")
	diff.NewFile = strings.TrimPrefix(lines[1], "+++ ")

	currentHunk := (*Hunk)(nil)
	oldLineNo := 0
	newLineNo := 0

	for i := 2; i < len(lines); i++ {
		line := lines[i]
		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "@@ ") {
			if currentHunk != nil {
				diff.Hunks = append(diff.Hunks, *currentHunk)
			}

			hunk, err := parseHunkHeader(line)
			if err != nil {
				return nil, fmt.Errorf("invalid hunk header: %v", err)
			}

			currentHunk = hunk
			oldLineNo = hunk.OldStart
			newLineNo = hunk.NewStart
			continue
		}

		if currentHunk == nil {
			continue
		}

		change := Change{}
		switch {
		case strings.HasPrefix(line, " "):
			change.Type = Unchanged
			change.Content = line[1:]
			change.OldLine = oldLineNo
			change.NewLine = newLineNo
			oldLineNo++
			newLineNo++
		case strings.HasPrefix(line, "+"):
			if strings.HasPrefix(line, "\\ No newline at end of file") {
				continue
			}
			change.Type = Added
			change.Content = line[1:]
			change.NewLine = newLineNo
			newLineNo++
		case strings.HasPrefix(line, "-"):
			if strings.HasPrefix(line, "\\ No newline at end of file") {
				continue
			}
			change.Type = Removed
			change.Content = line[1:]
			change.OldLine = oldLineNo
			oldLineNo++
		default:
			if strings.HasPrefix(line, "\\ No newline at end of file") {
				continue
			}
		}

		currentHunk.Changes = append(currentHunk.Changes, change)
	}

	if currentHunk != nil {
		diff.Hunks = append(diff.Hunks, *currentHunk)
	}

	return diff, nil
}

func parseHunkHeader(header string) (*Hunk, error) {
	matches := regexp.MustCompile(`@@ -(\d+)(?:,(\d+))? \+(\d+)(?:,(\d+))? @@(.*)$`).FindStringSubmatch(header)
	if matches == nil {
		return nil, fmt.Errorf("invalid hunk header format")
	}

	oldStart, _ := strconv.Atoi(matches[1])
	oldCount := 1
	if matches[2] != "" {
		oldCount, _ = strconv.Atoi(matches[2])
	}

	newStart, _ := strconv.Atoi(matches[3])
	newCount := 1
	if matches[4] != "" {
		newCount, _ = strconv.Atoi(matches[4])
	}

	context := ""
	if len(matches) > 5 {
		context = strings.TrimSpace(matches[5])
	}

	return &Hunk{
		OldStart: oldStart,
		OldCount: oldCount,
		NewStart: newStart,
		NewCount: newCount,
		Context:  context,
		Changes:  make([]Change, 0),
	}, nil
}
