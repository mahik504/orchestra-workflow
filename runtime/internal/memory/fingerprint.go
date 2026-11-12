package memory

import (
	"crypto/sha256"
	"encoding/hex"
	"regexp"
	"strings"
)

const MaxEvaluations = 500

var (
	alreadyHashed = regexp.MustCompile(`^[a-f0-9]{64}$`)
	secretPattern = regexp.MustCompile(`(?i)(?:sk-[A-Za-z0-9_-]{8,}|AKIA[0-9A-Z]{16}|ghp_[A-Za-z0-9]{20,}|api[_-]?key\s*[:=]\s*\S+|bearer\s+[A-Za-z0-9._-]{12,})`)
)

// Fingerprint stores a sha256 of task text. A value that is already a hex digest is kept.
func Fingerprint(text string) string {
	text = strings.TrimSpace(text)
	if text == "" || alreadyHashed.MatchString(text) {
		return text
	}
	sum := sha256.Sum256([]byte(text))
	return hex.EncodeToString(sum[:])
}

// Redact removes secret-shaped strings from notes that are not the task fingerprint.
func Redact(text string) string {
	return secretPattern.ReplaceAllString(text, "[redacted]")
}
