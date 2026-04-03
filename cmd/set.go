package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/moronim/llmvlt/store"
	"github.com/moronim/llmvlt/validator"
	"github.com/spf13/cobra"
)

var setCmd = &cobra.Command{
	Use:   "set KEY [VALUE]",
	Short: "Set a secret value",
	Long: `Set a secret in the vault. If VALUE is omitted, reads from stdin.

The tool validates key formats for known providers (OpenAI, Anthropic, etc.)
and warns if the format looks wrong — but never blocks the operation.

Examples:
  llmvlt set OPENAI_API_KEY sk-abc123...
  echo "sk-abc123..." | llmvlt set OPENAI_API_KEY`,
	Args: cobra.RangeArgs(1, 2),
	RunE: runSet,
}

func init() {
	rootCmd.AddCommand(setCmd)
}

func runSet(cmd *cobra.Command, args []string) error {
	key := args[0]
	var value string

	if len(args) == 2 {
		value = args[1]
	} else {
		// Read from stdin
		scanner := bufio.NewScanner(os.Stdin)
		if scanner.Scan() {
			value = strings.TrimSpace(scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			return fmt.Errorf("could not read from stdin: %w", err)
		}
	}

	if value == "" {
		return fmt.Errorf("value cannot be empty")
	}

	// Validate the key format if we know the provider
	result := validator.Validate(key, value)
	if result.Warning != "" {
		fmt.Fprintf(os.Stderr, "⚠ %s\n", result.Warning)
	}
	if result.Valid && result.Warning == "" {
		fmt.Fprintf(os.Stderr, "✓ %s format looks valid\n", key)
	}

	// Load vault
	password, err := getPassword()
	if err != nil {
		return err
	}

	storePath := getStorePath()
	v, err := store.Load(storePath, password)
	if err != nil {
		return fmt.Errorf("could not open vault: %w", err)
	}

	// Set and save
	v.Set(key, value)
	if err := store.Save(storePath, password, v); err != nil {
		return fmt.Errorf("could not save vault: %w", err)
	}

	fmt.Fprintf(os.Stderr, "✓ Secret %s saved\n", key)
	return nil
}
