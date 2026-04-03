package validator

import (
	"regexp"

	"github.com/moronim/llmvlt/preset"
)

// ValidationResult is returned after checking a secret value.
type ValidationResult struct {
	Valid   bool
	Warning string
}

// Validate checks a secret value against known provider patterns.
// Validation is advisory — it warns but never blocks.
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
			Valid:   false,
			Warning: hint,
		}
	}

	return ValidationResult{Valid: true}
}
