package commit

import (
	"fmt"
	"strings"

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

func NewCommitMessage(theme *huh.Theme) (*CommitMessage, error) {
	commitMsg := &CommitMessage{}
	descriptionStyle := theme.Help.FullDesc

	// TODOs:
	// * use file picker to add files to commit
	// * create full pipeline for Git functionality
	// * integrate local OLlama LLM to automatically generate the
	//   commit messages based on textual Git diff.
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[CommitType]().
				Title("Commit Type?").
				Options(
					createCommitTypeOption("build", "changes affects build system or external dependencies", Build, descriptionStyle),
					createCommitTypeOption("ci", "changes to CI configuration or scripts", CI, descriptionStyle),
					createCommitTypeOption("docs", "documentation-only changes", Docs, descriptionStyle),
					createCommitTypeOption("feat", "a new feature", Feat, descriptionStyle),
					createCommitTypeOption("fix", "a bug fix", Fix, descriptionStyle),
					createCommitTypeOption("perf", "changes that improve performance", Perf, descriptionStyle),
					createCommitTypeOption("refactor", "changes that neither fixes bug or adds feature", Refactor, descriptionStyle),
					createCommitTypeOption("style", "changes that do not affect meaning of code", Style, descriptionStyle),
					createCommitTypeOption("test", "changes add missing test or correct existing tests", Test, descriptionStyle),
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
	).WithTheme(theme)

	err := form.Run()
	return commitMsg, err
}
