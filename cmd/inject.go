package cmd

import (
	"fmt"

	"github.com/moronim/llmvlt/injector"
	"github.com/moronim/llmvlt/store"
	"github.com/spf13/cobra"
)

var injectCmd = &cobra.Command{
	Use:   "inject",
	Short: "Output secrets in a specific format for injection",
	Long: `Output all vault secrets in a format suitable for injection into your environment.

Formats:
  shell    — export KEY="value" statements (default)
  dotenv   — KEY=value lines for .env files
  jupyter  — Python os.environ assignments for notebook cells

Examples:
  eval $(llmvlt inject)
  llmvlt inject --format dotenv > .env
  llmvlt inject --format jupyter  # paste into notebook cell`,
	RunE: runInject,
}

func init() {
	rootCmd.AddCommand(injectCmd)
	injectCmd.Flags().String("format", "shell", "output format: shell, dotenv, jupyter")
	injectCmd.Flags().StringP("output", "o", "", "write to file instead of stdout")
}

func runInject(cmd *cobra.Command, args []string) error {
	password, err := getPassword()
	if err != nil {
		return err
	}

	v, err := store.Load(getStorePath(), password)
	if err != nil {
		return fmt.Errorf("could not open vault: %w", err)
	}

	format, _ := cmd.Flags().GetString("format")
	outFile, _ := cmd.Flags().GetString("output")

	secrets := v.All()

	// Filter out empty values
	filtered := make(map[string]string)
	for k, val := range secrets {
		if val != "" {
			filtered[k] = val
		}
	}

	if len(filtered) == 0 {
		return fmt.Errorf("no secrets with values in vault — nothing to inject")
	}

	inj, err := injector.Get(format)
	if err != nil {
		return err
	}

	output, err := inj.Inject(filtered)
	if err != nil {
		return fmt.Errorf("injection failed: %w", err)
	}

	if outFile != "" {
		return injector.WriteToFile(outFile, output)
	}

	fmt.Print(output)
	return nil
}
