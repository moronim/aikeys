package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/moronim/llmvlt/history"
	"github.com/moronim/llmvlt/store"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run -- COMMAND [ARGS...]",
	Short: "Run a command with secrets injected as environment variables",
	Long: `Execute a command with all vault secrets injected into its environment.

Secrets are injected into the child process only — they never appear in the
parent shell, shell history, or ps output.

Examples:
  llmvlt run -- python train.py
  llmvlt run --tag "gpt4-experiment-1" -- python eval.py
  llmvlt run -- jupyter notebook`,
	DisableFlagParsing: false,
	RunE:               runRun,
}

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().String("tag", "", "tag this run for experiment tracking")
}

func runRun(cmd *cobra.Command, args []string) error {
	// Find the command after "--"
	dashIdx := -1
	for i, a := range os.Args {
		if a == "--" {
			dashIdx = i
			break
		}
	}

	var cmdArgs []string
	if dashIdx >= 0 && dashIdx+1 < len(os.Args) {
		cmdArgs = os.Args[dashIdx+1:]
	} else if len(args) > 0 {
		cmdArgs = args
	}

	if len(cmdArgs) == 0 {
		return fmt.Errorf("no command specified — usage: llmvlt run -- <command>")
	}

	password, err := getPassword()
	if err != nil {
		return err
	}

	v, err := store.Load(getStorePath(), password)
	if err != nil {
		return fmt.Errorf("could not open vault: %w", err)
	}

	// Build environment: current env + vault secrets overlaid
	env := os.Environ()
	secrets := v.All()
	usedKeys := make([]string, 0, len(secrets))
	for k, val := range secrets {
		if val != "" { // Only inject non-empty secrets
			env = append(env, fmt.Sprintf("%s=%s", k, val))
			usedKeys = append(usedKeys, k)
		}
	}

	// Log to history
	tag, _ := cmd.Flags().GetString("tag")
	history.Log(cmdArgs, usedKeys, tag)

	// Fork subprocess
	child := exec.Command(cmdArgs[0], cmdArgs[1:]...)
	child.Env = env
	child.Stdin = os.Stdin
	child.Stdout = os.Stdout
	child.Stderr = os.Stderr

	// Forward signals to child
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		for sig := range sigCh {
			if child.Process != nil {
				child.Process.Signal(sig)
			}
		}
	}()

	if err := child.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			os.Exit(exitErr.ExitCode())
		}
		return fmt.Errorf("could not execute command: %w", err)
	}

	return nil
}
