package server

import (
	"bytes"
	"strings"
)

// extractUnifiedDiff tries to pull a clean unified diff out of a noisy LLM output.
// It only returns a patch if it finds at least one valid `diff --git` header.
func extractUnifiedDiff(out []byte) []byte {
	i := bytes.Index(out, []byte("diff --git"))
	if i < 0 {
		return nil
	}
	lines := strings.Split(string(out[i:]), "\n")
	kept := make([]string, 0, len(lines))
	started := false
	for _, raw := range lines {
		l := strings.TrimSuffix(raw, "\r")
		if l == "" {
			if started {
				kept = append(kept, l)
			}
			continue
		}

		// Ignore common tool chatter/noise.
		if strings.HasPrefix(l, "tokens used") || strings.HasPrefix(l, "codex") || strings.HasPrefix(l, "thinking") || strings.HasPrefix(l, "exec") || strings.HasPrefix(l, "file update") || strings.HasPrefix(l, "succeeded in") {
			continue
		}

		// Normalize common quoting artifacts.
		n := strings.TrimLeft(l, "`\"'")
		n = strings.TrimRight(n, "`\"'.")

		// Headers
		if strings.HasPrefix(n, "diff --git ") {
			started = true
			kept = append(kept, n)
			continue
		}

		if !started {
			// Don't accept hunks without a header.
			continue
		}

		// Unified diff + git-style headers.
		switch {
		case strings.HasPrefix(n, "index "):
		case strings.HasPrefix(n, "--- "):
		case strings.HasPrefix(n, "+++ "):
		case strings.HasPrefix(n, "@@ "):
		case strings.HasPrefix(n, "new file mode "):
		case strings.HasPrefix(n, "deleted file mode "):
		case strings.HasPrefix(n, "similarity index "):
		case strings.HasPrefix(n, "rename from "):
		case strings.HasPrefix(n, "rename to "):
		case strings.HasPrefix(n, "\\ No newline at end of file"):
		case strings.HasPrefix(n, "+"):
		case strings.HasPrefix(n, "-"):
		case strings.HasPrefix(n, " "):
		default:
			continue
		}
		kept = append(kept, n)
	}
	if !started || len(kept) == 0 {
		return nil
	}
	return []byte(strings.Join(kept, "\n") + "\n")
}
