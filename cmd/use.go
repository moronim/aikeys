package cmd

import (
	"fmt"
	"os"

	"github.com/moronim/llmvlt/preset"
	"github.com/moronim/llmvlt/store"
	"github.com/spf13/cobra"
)

var useCmd = &cobra.Command{
	Use:   "use PROVIDER",
	Short: "Activate a specific provider's keys (deactivate others)",
	Long: `Switch active provider context. This is useful when benchmarking the same
script across multiple LLM providers.

When you "use" a provider, only that provider's secrets are injected by
'llmvlt run'. Other secrets remain in the vault but are not exported.

Use "all" to re-enable all providers.

Examples:
  llmvlt use anthropic   # only Anthropic keys active
  llmvlt use openai      # switch to OpenAI
  llmvlt use all         # re-enable everything`,
	Args: cobra.ExactArgs(1),
	RunE: runUse,
}

func init() {
	rootCmd.AddCommand(useCmd)
}

func runUse(cmd *cobra.Command, args []string) error {
	provider := args[0]

	if provider != "all" {
		// Verify this is a known provider
		if _, err := preset.Get(provider + "-stack"); err != nil {
			// Try without -stack suffix
			if _, err := preset.Get(provider); err != nil {
				return fmt.Errorf("unknown provider %q — run 'llmvlt presets' to see available providers", provider)
			}
		}
	}

	password, err := getPassword()
	if err != nil {
		return err
	}

	storePath := getStorePath()
	v, err := store.Load(storePath, password)
	if err != nil {
		return fmt.Errorf("could not open vault: %w", err)
	}

	v.SetActiveProvider(provider)

	if err := store.Save(storePath, password, v); err != nil {
		return fmt.Errorf("could not save vault: %w", err)
	}

	if provider == "all" {
		fmt.Fprintln(os.Stderr, "✓ All providers active")
	} else {
		fmt.Fprintf(os.Stderr, "✓ Switched to %s — only %s keys will be injected by 'llmvlt run'\n", provider, provider)
	}

	return nil
}
