package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "llmvlt",
	Short: "A secret manager built for AI/ML engineers",
	Long: `llmvlt is a CLI secret manager designed specifically for AI/ML engineers.
It knows your providers (OpenAI, Anthropic, HuggingFace, etc.), validates
key formats, injects secrets into your scripts, and tracks which experiments
used which key versions.`,
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd: true,
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is .llmvlt.yml)")
	rootCmd.PersistentFlags().StringP("password", "p", "", "master password for the vault")
	rootCmd.PersistentFlags().StringP("store", "s", "", "path to the vault store file")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "verbose output")

	viper.BindPFlag("password", rootCmd.PersistentFlags().Lookup("password"))
	viper.BindPFlag("store", rootCmd.PersistentFlags().Lookup("store"))
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		// Search in current directory, then home
		viper.AddConfigPath(".")
		viper.AddConfigPath(filepath.Join(home, ".llmvlt"))
		viper.SetConfigName(".llmvlt")
		viper.SetConfigType("yml")
	}

	viper.SetEnvPrefix("LLMVLT")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_", ".", "_"))

	viper.ReadInConfig()

	// Set defaults
	if viper.GetString("store") == "" {
		viper.Set("store", ".llmvlt.store")
	}
}

func getStorePath() string {
	return viper.GetString("store")
}

func getPassword() (string, error) {
	pw := viper.GetString("password")
	if pw != "" {
		return pw, nil
	}

	// Check environment
	pw = os.Getenv("LLMVLT_PASSWORD")
	if pw != "" {
		return pw, nil
	}

	// Prompt interactively
	return promptPassword("Enter vault password: ")
}
