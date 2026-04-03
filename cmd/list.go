package cmd

import (
	"fmt"
	"sort"

	"github.com/moronim/llmvlt/preset"
	"github.com/moronim/llmvlt/store"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List all secret keys in the vault",
	Long: `List all secret key names stored in the vault. Values are never shown.

Keys with empty values (scaffolded but not yet filled) are marked with ⬚.
Keys recognized as part of a known provider preset show the provider name.

Examples:
  llmvlt list
  llmvlt ls`,
	RunE: runList,
}

func init() {
	rootCmd.AddCommand(listCmd)
}

func runList(cmd *cobra.Command, args []string) error {
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
		fmt.Println("Vault is empty. Add secrets with: llmvlt set KEY value")
		return nil
	}

	sort.Strings(keys)

	for _, key := range keys {
		value, _ := v.Get(key)
		status := "✓"
		if value == "" {
			status = "⬚"
		}

		// Check if key belongs to a known provider
		provider := preset.ProviderForKey(key)
		if provider != "" {
			fmt.Printf("  %s %s  (%s)\n", status, key, provider)
		} else {
			fmt.Printf("  %s %s\n", status, key)
		}
	}

	return nil
}
