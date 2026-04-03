package validator

import (
	"fmt"
	"regexp"

	"github.com/moronim/llmvlt/preset"
)

// ValidationResult is returned after checking a secret value.
type ValidationResult struct {
	Valid     bool
	Checked   bool   // true if a pattern was actually tested against the value
	Error     string // non-empty if validation failed (blocks without --force)
}

// Validate checks a secret value against known provider patterns.
// Returns Valid=true for unknown keys or keys without patterns.
// Returns Valid=false with an Error message for known keys that fail format checks.
func Validate(key, value string) ValidationResult {
	def := preset.SecretDefForKey(key)
	if def == nil {
		// Unknown key — no validation available
		return ValidationResult{Valid: true}
	}

	if def.Pattern == "" {
		// Known key but no pattern defined
		return ValidationResult{Valid: true}
	}

	matched, err := regexp.MatchString(def.Pattern, value)
	if err != nil {
		// Regex error — don't punish the user
		return ValidationResult{Valid: true}
	}

	if !matched {
		hint := def.PatternHint
		if hint == "" {
			hint = "value doesn't match expected format"
		}
		return ValidationResult{
			Valid: false,
			Error: fmt.Sprintf("%s: %s", key, hint),
		}
	}

	return ValidationResult{Valid: true, Checked: true}
}
