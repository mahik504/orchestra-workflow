package contractpack

import "strings"

// Required drops the design pack on a small task.
// The domain decides which files apply. This function decides whether the task needs them.
func Required(task string, documents []string) []string {
	switch strings.ToLower(strings.TrimSpace(task)) {
	case "", "bug", "bugfix", "typo", "glue", "control":
		return nil
	}
	out := make([]string, 0, len(documents))
	seen := map[string]bool{}
	for _, doc := range documents {
		doc = strings.TrimSpace(doc)
		if doc == "" || seen[doc] {
			continue
		}
		seen[doc] = true
		out = append(out, doc)
	}
	if len(out) == 0 {
		return nil
	}
	return out
}
