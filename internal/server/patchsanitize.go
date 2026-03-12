package server

import (
	"strings"
)

// sanitizePatch removes leading/trailing noise and drops stray diff headers that are not followed by
// file headers. This protects git apply from "fragment without header" and duplicate header-only lines.
func sanitizePatch(patch []byte) []byte {
	lines := strings.Split(string(patch), "\n")
	out := make([]string, 0, len(lines))

	for i := 0; i < len(lines); i++ {
		l := lines[i]
		if strings.HasPrefix(l, "diff --git ") {
			// Look ahead a few lines to confirm this is a real diff block.
			hasFileHeaders := false
			for j := i + 1; j < len(lines) && j <= i+8; j++ {
				if strings.HasPrefix(lines[j], "--- ") || strings.HasPrefix(lines[j], "+++ ") {
					hasFileHeaders = true
					break
				}
				if strings.HasPrefix(lines[j], "diff --git ") {
					break
				}
			}
			if !hasFileHeaders {
				continue
			}
		}
		out = append(out, l)
	}

	// trim leading empty lines
	for len(out) > 0 && strings.TrimSpace(out[0]) == "" {
		out = out[1:]
	}
	// trim trailing empty lines
	for len(out) > 0 && strings.TrimSpace(out[len(out)-1]) == "" {
		out = out[:len(out)-1]
	}

	return []byte(strings.Join(out, "\n") + "\n")
}
