package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "githelper",
	Short: "Utilities for command-line git",
	Long: `A tool to simplify the complexities of command-line git. Currently, 
it provides a tool to generate commit messages in the style of Conventional
Commits, and more are planned. Eventually, this tool should become a wrapper
for git, providing higher level functionality without needing to expose the inner
details of the git commands needed to run these tasks.
	`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(loadConfiguration)
}

func loadConfiguration() {
	home, err := os.UserHomeDir()
	if err != nil {
		os.Exit(1)
	}

	viper.AddConfigPath(home)
	viper.SetConfigName(".githelper")
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()
	err = viper.ReadInConfig()
	if err != nil {
		os.Exit(1)
	}
}
