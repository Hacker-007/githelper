package cmd

import (
	"fmt"

	"github.com/Hacker-007/githelper/internal/commit"
	"github.com/atotto/clipboard"
	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

var genCommitCmd = &cobra.Command{
	Use:     "gen-commit",
	Short:   "Generate a Conventional Commit message",
	Long:    `Runs an interactive CLI session to generate a Conventional Commit message.`,
	Aliases: []string{"gc"},
	RunE: func(cmd *cobra.Command, args []string) error {
		theme := huh.ThemeCatppuccin()
		commitMsg, err := commit.NewCommitMessage(theme)
		if err != nil {
			return err
		}

		err = clipboard.WriteAll(commitMsg.String())
		if err != nil {
			return err
		}

		fmt.Println("Commit message copied to clipboard ✔︎")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(genCommitCmd)
}
