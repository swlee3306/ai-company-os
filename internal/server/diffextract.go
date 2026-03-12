package server

import (
	"bytes"
	"strings"
)

// extractUnifiedDiff tries to pull a clean unified diff out of a noisy LLM output.
// Strategy:
// - Find all occurrences of "diff --git".
// - For each candidate block, keep only patch-looking lines (skip tool chatter).
// - Prefer a candidate that contains file headers (---/+++), and at least one hunk (@@).
func extractUnifiedDiff(out []byte) []byte {
	idxs := []int{}
	needle := []byte("diff --git")
	for i := 0; ; {
		j := bytes.Index(out[i:], needle)
		if j < 0 {
			break
		}
		idxs = append(idxs, i+j)
		i = i + j + len(needle)
	}
	if len(idxs) == 0 {
		return nil
	}

	type cand struct {
		patch      []byte
		hasHeaders bool
		hasHunks   bool
		lines      int
	}
	cands := []cand{}

	for _, start := range idxs {
		lines := strings.Split(string(out[start:]), "\n")
		kept := make([]string, 0, len(lines))
		hasHeaders := false
		hasHunks := false
		started := false
		seenMinus := false
		seenPlus := false

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

			if strings.HasPrefix(n, "diff --git ") {
				// If we already captured a complete diff (headers+hunks), stop at the next diff.
				if started && (seenMinus && seenPlus) && hasHunks {
					break
				}
				started = true
				kept = append(kept, n)
				continue
			}
			if !started {
				continue
			}

			switch {
			case strings.HasPrefix(n, "index "):
			case strings.HasPrefix(n, "--- "):
				hasHeaders = true
				seenMinus = true
			case strings.HasPrefix(n, "+++ "):
				hasHeaders = true
				seenPlus = true
			case strings.HasPrefix(n, "@@ "):
				// Do not accept hunks before file headers are present.
				if !(seenMinus && seenPlus) {
					continue
				}
				hasHunks = true
			case strings.HasPrefix(n, "new file mode "):
			case strings.HasPrefix(n, "deleted file mode "):
			case strings.HasPrefix(n, "similarity index "):
			case strings.HasPrefix(n, "rename from "):
			case strings.HasPrefix(n, "rename to "):
			case strings.HasPrefix(n, "\\ No newline at end of file"):
			case strings.HasPrefix(n, "+"):
				if !(seenMinus && seenPlus) {
					continue
				}
			case strings.HasPrefix(n, "-"):
				if !(seenMinus && seenPlus) {
					continue
				}
			case strings.HasPrefix(n, " "):
				if !(seenMinus && seenPlus) {
					continue
				}
			default:
				// keep scanning; this candidate may contain multiple diffs.
				continue
			}
			kept = append(kept, n)

			// Heuristic end: once we have headers+hunks and then we see a new diff later, we can stop.
			// (We don't implement explicit stop here; selection logic prefers the best match anyway.)
		}

		if !started || len(kept) == 0 {
			continue
		}
		p := []byte(strings.Join(kept, "\n") + "\n")
		cands = append(cands, cand{patch: p, hasHeaders: hasHeaders, hasHunks: hasHunks, lines: len(kept)})
	}

	if len(cands) == 0 {
		return nil
	}

	// Prefer a candidate with headers + hunks, then largest line count.
	best := cands[0]
	for _, c := range cands[1:] {
		bScore := 0
		cScore := 0
		if best.hasHeaders {
			bScore += 2
		}
		if best.hasHunks {
			bScore += 1
		}
		if c.hasHeaders {
			cScore += 2
		}
		if c.hasHunks {
			cScore += 1
		}
		if cScore > bScore || (cScore == bScore && c.lines > best.lines) {
			best = c
		}
	}

	return best.patch
}
