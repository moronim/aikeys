package cmd

import (
	"fmt"

	"github.com/moronim/llmvlt/history"
	"github.com/spf13/cobra"
)

var historyCmd = &cobra.Command{
	Use:   "history",
	Short: "Show experiment run history",
	Long: `Display a log of commands executed via 'llmvlt run', including which
secrets were active and any experiment tags.

Examples:
  llmvlt history
  llmvlt history --last 10`,
	RunE: runHistory,
}

func init() {
	rootCmd.AddCommand(historyCmd)
	historyCmd.Flags().Int("last", 20, "number of entries to show")
}

func runHistory(cmd *cobra.Command, args []string) error {
	last, _ := cmd.Flags().GetInt("last")
	entries, err := history.Read(last)
	if err != nil {
		return fmt.Errorf("could not read history: %w", err)
	}

	if len(entries) == 0 {
		fmt.Println("No history yet. Run a command with: llmvlt run -- <command>")
		return nil
	}

	for _, e := range entries {
		tag := ""
		if e.Tag != "" {
			tag = fmt.Sprintf(" [%s]", e.Tag)
		}
		fmt.Printf("  %s  %s%s\n", e.Timestamp.Format("2006-01-02 15:04:05"), e.Command, tag)
		for _, k := range e.Keys {
			fmt.Printf("           ↳ %s\n", k)
		}
	}

	return nil
}
