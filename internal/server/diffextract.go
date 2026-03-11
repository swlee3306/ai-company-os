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
		l = strings.TrimSuffix(l, "\r")
		if l == "" {
			// allow blank lines in diffs (rare), but ignore if we haven't started
			if len(kept) > 0 {
				kept = append(kept, l)
			}
			continue
		}
		// Ignore common tool chatter/noise instead of stopping extraction.
		if strings.HasPrefix(l, "tokens used") || strings.HasPrefix(l, "codex") || strings.HasPrefix(l, "thinking") || strings.HasPrefix(l, "exec") || strings.HasPrefix(l, "file update") {
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
			// skip non-patch lines; do not terminate early
			continue
		}
		kept = append(kept, l)
	}
	if len(kept) == 0 {
		return nil
	}
	return []byte(strings.Join(kept, "\n") + "\n")
}
