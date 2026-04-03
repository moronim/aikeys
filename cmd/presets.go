package cmd

import (
	"fmt"

	"github.com/moronim/llmvlt/preset"
	"github.com/spf13/cobra"
)

var presetsCmd = &cobra.Command{
	Use:   "presets",
	Short: "List available provider presets",
	Long: `Show all built-in presets and their secrets.

Examples:
  llmvlt presets
  llmvlt presets --detail`,
	RunE: runPresets,
}

func init() {
	rootCmd.AddCommand(presetsCmd)
	presetsCmd.Flags().Bool("detail", false, "show individual secrets in each preset")
}

func runPresets(cmd *cobra.Command, args []string) error {
	detail, _ := cmd.Flags().GetBool("detail")
	all := preset.All()

	for _, p := range all {
		fmt.Printf("  %s — %s\n", p.Name, p.Description)
		if detail {
			for _, s := range p.AllSecrets() {
				req := ""
				if !s.Required {
					req = " (optional)"
				}
				fmt.Printf("      %s%s\n", s.Key, req)
			}
			fmt.Println()
		}
	}

	return nil
}
