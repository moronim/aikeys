package cmd

import (
	"fmt"
	"os"

	"github.com/moronim/llmvlt/store"
	"github.com/spf13/cobra"
)

var unsetCmd = &cobra.Command{
	Use:   "unset KEY",
	Short: "Remove a secret from the vault",
	Long: `Remove a secret by key. This is permanent.

Examples:
  llmvlt unset OPENAI_API_KEY`,
	Args: cobra.ExactArgs(1),
	RunE: runUnset,
}

func init() {
	rootCmd.AddCommand(unsetCmd)
}

func runUnset(cmd *cobra.Command, args []string) error {
	key := args[0]

	password, err := getPassword()
	if err != nil {
		return err
	}

	storePath := getStorePath()
	v, err := store.Load(storePath, password)
	if err != nil {
		return fmt.Errorf("could not open vault: %w", err)
	}

	if _, ok := v.Get(key); !ok {
		return fmt.Errorf("secret %q not found in vault", key)
	}

	v.Unset(key)

	if err := store.Save(storePath, password, v); err != nil {
		return fmt.Errorf("could not save vault: %w", err)
	}

	fmt.Fprintf(os.Stderr, "✓ Secret %s removed\n", key)
	return nil
}
