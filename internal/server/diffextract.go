package server

import (
	"bytes"
	"strings"
)

// extractUnifiedDiff tries to pull a clean unified diff out of a noisy LLM output.
// It keeps only patch-looking lines after the first "diff --git".
func extractUnifiedDiff(out []byte) []byte {
	i := bytes.Index(out, []byte("diff --git"))
	if i < 0 {
		return nil
	}
	lines := strings.Split(string(out[i:]), "\n")
	kept := make([]string, 0, len(lines))
	for _, l := range lines {
		// Stop if we clearly left patch mode (common noisy footers)
		if strings.HasPrefix(l, "tokens used") {
			break
		}
		if strings.HasPrefix(l, "codex") || strings.HasPrefix(l, "thinking") || strings.HasPrefix(l, "exec") {
			break
		}
		if strings.HasPrefix(l, "file update") {
			// ignore tool chatter
			continue
		}

		// Unified diff + git-style headers.
		switch {
		case strings.HasPrefix(l, "diff --git "):
		case strings.HasPrefix(l, "index "):
		case strings.HasPrefix(l, "--- "):
		case strings.HasPrefix(l, "+++ "):
		case strings.HasPrefix(l, "@@ "):
		case strings.HasPrefix(l, "new file mode "):
		case strings.HasPrefix(l, "deleted file mode "):
		case strings.HasPrefix(l, "similarity index "):
		case strings.HasPrefix(l, "rename from "):
		case strings.HasPrefix(l, "rename to "):
		case strings.HasPrefix(l, "\\ No newline at end of file"):
		case strings.HasPrefix(l, "+"):
		case strings.HasPrefix(l, "-"):
		case strings.HasPrefix(l, " "):
		default:
			// If we hit a line that doesn't look like patch content, stop.
			return []byte(strings.Join(kept, "\n") + "\n")
		}
		kept = append(kept, l)
	}
	if len(kept) == 0 {
		return nil
	}
	return []byte(strings.Join(kept, "\n") + "\n")
}
