package cmd

import (
	"fmt"
	"os"

	"github.com/moronim/llmvlt/preset"
	"github.com/moronim/llmvlt/store"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new vault",
	Long: `Initialize a new encrypted vault, optionally with a provider preset.

Examples:
  llmvlt init
  llmvlt init --preset openai-stack
  llmvlt init --preset full-llm-stack`,
	RunE: runInit,
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().String("preset", "", "provider preset to scaffold (e.g. openai-stack, full-llm-stack)")
}

func runInit(cmd *cobra.Command, args []string) error {
	storePath := getStorePath()

	// Check if store already exists
	if _, err := os.Stat(storePath); err == nil {
		return fmt.Errorf("vault already exists at %s — use 'llmvlt set' to add secrets", storePath)
	}

	// Get password — use -p flag or env var if provided, otherwise prompt with confirmation
	password := viper.GetString("password")
	if password == "" {
		password = os.Getenv("LLMVLT_PASSWORD")
	}
	if password == "" {
		var err error
		password, err = promptPasswordConfirm()
		if err != nil {
			return err
		}
	}

	// Create empty vault
	v := store.NewVault()

	// If preset specified, scaffold empty keys
	presetName, _ := cmd.Flags().GetString("preset")
	if presetName != "" {
		p, err := preset.Get(presetName)
		if err != nil {
			return fmt.Errorf("unknown preset %q — run 'llmvlt presets' to see available presets", presetName)
		}

		for _, s := range p.AllSecrets() {
			v.Set(s.Key, "")
		}

		fmt.Fprintf(os.Stderr, "✓ Initialized vault with preset: %s\n", p.Name)
		fmt.Fprintf(os.Stderr, "  %d secrets scaffolded. Fill them in with:\n", len(p.AllSecrets()))
		for _, s := range p.AllSecrets() {
			marker := "(required)"
			if !s.Required {
				marker = "(optional)"
			}
			fmt.Fprintf(os.Stderr, "    llmvlt set %s  %s\n", s.Key, marker)
		}
	} else {
		fmt.Fprintln(os.Stderr, "✓ Initialized empty vault")
		fmt.Fprintln(os.Stderr, "  Add secrets with: llmvlt set KEY value")
		fmt.Fprintln(os.Stderr, "  Or init with a preset: llmvlt init --preset openai-stack")
	}

	// Save
	return store.Save(storePath, password, v)
}
