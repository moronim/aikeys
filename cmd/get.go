package cmd

import (
	"fmt"
	"os"

	"github.com/moronim/llmvlt/store"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var getCmd = &cobra.Command{
	Use:   "get KEY",
	Short: "Retrieve a secret value",
	Long: `Retrieve and print a secret from the vault.

If stdout is a terminal, the value is printed with a trailing newline.
If piped, the raw value is output (no newline) for safe scripting.

Examples:
  llmvlt get OPENAI_API_KEY
  export OPENAI_API_KEY=$(llmvlt get OPENAI_API_KEY)`,
	Args: cobra.ExactArgs(1),
	RunE: runGet,
}

func init() {
	rootCmd.AddCommand(getCmd)
}

func runGet(cmd *cobra.Command, args []string) error {
	key := args[0]

	password, err := getPassword()
	if err != nil {
		return err
	}

	v, err := store.Load(getStorePath(), password)
	if err != nil {
		return fmt.Errorf("could not open vault: %w", err)
	}

	value, ok := v.Get(key)
	if !ok {
		return fmt.Errorf("secret %q not found in vault", key)
	}

	if value == "" {
		return fmt.Errorf("secret %q exists but has no value — set it with: llmvlt set %s <value>", key, key)
	}

	// If stdout is a TTY, add newline. If piped, output raw.
	if term.IsTerminal(int(os.Stdout.Fd())) {
		fmt.Println(value)
	} else {
		fmt.Print(value)
	}

	return nil
}
