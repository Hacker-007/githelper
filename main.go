package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

type CommitType int

var commitTypes = [...]string{"build", "ci", "docs", "feat", "fix", "perf", "refactor", "style", "test"}

const (
	Build CommitType = iota
	CI
	Docs
	Feat
	Fix
	Perf
	Refactor
	Style
	Test
)

type CommitMessage struct {
	ty          CommitType
	description string
	scope       string
	body        string
	footer      string
}

func (commit *CommitMessage) String() string {
	var sb strings.Builder
	sb.WriteString(commitTypes[commit.ty])
	if commit.scope != "" {
		sb.WriteString(fmt.Sprintf("(%s)", commit.scope))
	}

	sb.WriteString(fmt.Sprintf(": %s", commit.description))
	if commit.body != "" {
		sb.WriteString(fmt.Sprintf("\n\n%s", commit.body))
	}

	if commit.footer != "" {
		sb.WriteString(fmt.Sprintf("\n\n%s", commit.footer))
	}

	return sb.String()
}

func createCommitTypeOption(name string, desc string, ty CommitType, theme lipgloss.Style) huh.Option[CommitType] {
	return huh.NewOption(fmt.Sprintf("%8s: %s", name, theme.Render(desc)), ty)
}

func main() {
	commitMsg := &CommitMessage{}

	theme := huh.ThemeBase()
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[CommitType]().
				Title("Commit Type?").
				Options(
					createCommitTypeOption("build", "changes affects build system or external dependencies", Build, theme.Help.FullDesc),
					createCommitTypeOption("ci", "changes to CI configuration or scripts", CI, theme.Help.FullDesc),
					createCommitTypeOption("docs", "documentation-only changes", Docs, theme.Help.FullDesc),
					createCommitTypeOption("feat", "a new feature", Feat, theme.Help.FullDesc),
					createCommitTypeOption("fix", "a bug fix", Fix, theme.Help.FullDesc),
					createCommitTypeOption("perf", "changes that improve performance", Perf, theme.Help.FullDesc),
					createCommitTypeOption("refactor", "changes that neither fixes bug or adds feature", Refactor, theme.Help.FullDesc),
					createCommitTypeOption("style", "changes that do not affect meaning of code", Style, theme.Help.FullDesc),
					createCommitTypeOption("test", "changes add missing test or correct existing tests", Test, theme.Help.FullDesc),
				).
				Value(&commitMsg.ty),
			huh.NewInput().
				Title("Description?").
				Prompt(" ").
				Inline(true).
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("commit description is required")
					}

					return nil
				}).
				Value(&commitMsg.description),
		),
		huh.NewGroup(
			huh.NewInput().
				Title("Scope (optional)?").
				Prompt(" ").
				Inline(true).
				Value(&commitMsg.scope),
			huh.NewText().
				Title("Body (optional)?").
				Value(&commitMsg.body),
			huh.NewText().
				Title("Footer (optional)?").
				Value(&commitMsg.footer),
		),
	)

	err := form.Run()
	if err != nil {
		log.Fatal(err)
	}

	err = clipboard.WriteAll(commitMsg.String())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Commit message copied to clipboard ✔︎")
}
