package main

import (
	"fmt"
	"log"

	"github.com/Hacker-007/githelper/internal/commit"
	"github.com/atotto/clipboard"
	"github.com/charmbracelet/huh"
)

func main() {
	theme := huh.ThemeBase()

	commitMsg, err := commit.NewCommitMessage(theme.Help.FullDesc)
	if err != nil {
		log.Fatal(err)
	}

	err = clipboard.WriteAll(commitMsg.String())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Commit message copied to clipboard ✔︎")
}
