package cmd

import (
	"fmt"
	"os"

	"github.com/moronim/llmvlt/preset"
	"github.com/moronim/llmvlt/store"
	"github.com/moronim/llmvlt/validator"
	"github.com/spf13/cobra"
)

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Check vault health: missing values, format warnings, rotation reminders",
	Long: `Scan your vault and report issues:
  - Secrets scaffolded but never filled in
  - Values that don't match expected provider formats
  - Keys that haven't been rotated in a while (based on preset rotation_days)

Examples:
  llmvlt check`,
	RunE: runCheck,
}

func init() {
	rootCmd.AddCommand(checkCmd)
}

func runCheck(cmd *cobra.Command, args []string) error {
	password, err := getPassword()
	if err != nil {
		return err
	}

	v, err := store.Load(getStorePath(), password)
	if err != nil {
		return fmt.Errorf("could not open vault: %w", err)
	}

	keys := v.Keys()
	if len(keys) == 0 {
		fmt.Println("Vault is empty.")
		return nil
	}

	issues := 0

	for _, key := range keys {
		value, _ := v.Get(key)

		// Check for empty values
		if value == "" {
			fmt.Fprintf(os.Stderr, "⬚ %s — empty, needs a value\n", key)
			issues++
			continue
		}

		// Validate format
		result := validator.Validate(key, value)
		if result.Warning != "" {
			fmt.Fprintf(os.Stderr, "⚠ %s — %s\n", key, result.Warning)
			issues++
			continue
		}

		// Check rotation recommendation
		secretDef := preset.SecretDefForKey(key)
		if secretDef != nil && secretDef.RotationDays > 0 {
			age := v.SecretAgeDays(key)
			if age > secretDef.RotationDays {
				fmt.Fprintf(os.Stderr, "⏳ %s — last set %d days ago, rotation recommended every %d days\n",
					key, age, secretDef.RotationDays)
				issues++
				continue
			}
		}

		fmt.Fprintf(os.Stderr, "✓ %s\n", key)
	}

	if issues == 0 {
		fmt.Fprintln(os.Stderr, "\nAll secrets look good.")
	} else {
		fmt.Fprintf(os.Stderr, "\n%d issue(s) found.\n", issues)
	}

	return nil
}
