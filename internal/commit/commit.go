package commit

import (
	"encoding/json"
	"fmt"
	"strings"

	_ "embed"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

type CommitType int

//go:embed prompts/git_diff_summary.txt
var gitDiffSummaryPrompt string

//go:embed prompts/git_commit_generation.txt
var gitCommitGenerationPrompt string

var commitTypes = []string{"build", "ci", "docs", "feat", "fix", "perf", "refactor", "style", "test"}
var commitMapping = map[string]CommitType{"build": Build, "ci": CI, "docs": Docs, "feat": Feat, "fix": Fix, "perf": Perf, "refactor": Refactor, "style": Style, "test": Test}

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
	ty             CommitType
	description    string
	scope          string
	body           string
	breakingChange string
}

func (commit *CommitMessage) UnmarshalJSON(data []byte) error {
	if string(data) == "null" || string(data) == "" {
		return nil
	}

	var realCommit struct {
		Type           string
		Description    string
		Scope          string
		Body           string
		BreakingChange string
	}

	if err := json.Unmarshal(data, &realCommit); err != nil {
		return err
	}

	*commit = CommitMessage{
		ty:             commitMapping[realCommit.Type],
		description:    realCommit.Description,
		scope:          realCommit.Scope,
		body:           realCommit.Body,
		breakingChange: realCommit.BreakingChange,
	}

	return nil
}

func (commit *CommitMessage) String() string {
	var sb strings.Builder
	sb.WriteString(commitTypes[commit.ty])
	if commit.scope != "" {
		sb.WriteString(fmt.Sprintf("(%s)", commit.scope))
	}

	if commit.breakingChange != "" {
		sb.WriteString("!")
	}

	sb.WriteString(fmt.Sprintf(": %s", commit.description))
	if commit.body != "" {
		sb.WriteString(fmt.Sprintf("\n\n%s", commit.body))
	}

	hasFooter := commit.breakingChange != ""
	if hasFooter {
		sb.WriteString("\n\n")
	}

	if commit.breakingChange != "" {
		sb.WriteString(fmt.Sprintf("BREAKING CHANGES: %s", commit.breakingChange))
	}

	return sb.String()
}

func getDiffSummaries() ([]string, error) {
	changedFiles, err := getChangedFiles()
	if err != nil {
		return []string{}, err
	}

	if len(changedFiles) == 0 {
		return []string{}, nil
	}

	llmClient := &LLMClient{}
	summaries := []string{}
	for _, file := range changedFiles {
		diff, err := getDiff(file)
		if err != nil {
			return []string{}, err
		}

		prompt := fmt.Sprintf("\"\"\"\n%s\n\n%s\n\"\"\"", gitDiffSummaryPrompt, diff)
		summary, err := llmClient.SendStandaloneMessage(prompt)
		if err != nil {
			return []string{}, err
		}

		summaries = append(summaries, summary)
	}

	return summaries, nil
}

func getCommitDetails(summaries []string) (*CommitMessage, error) {
	llmClient := &LLMClient{}
	_, err := llmClient.SendChatMessage(LLMMessage{
		Role:    "system",
		Content: gitCommitGenerationPrompt,
	})

	if err != nil {
		return nil, err
	}

	for idx, summary := range summaries {
		_, err := llmClient.SendChatMessage(LLMMessage{
			Role:    "user",
			Content: fmt.Sprintf("\"\"\"\nSummary %d:\n%s\n\"\"\"", idx+1, summary),
		})

		if err != nil {
			return nil, err
		}
	}

	commit, err := llmClient.SendChatMessage(LLMMessage{
		Role:    "user",
		Content: "done",
	})

	if err != nil {
		return nil, err
	}

	fmt.Println(commit.Content)
	var commitMessage CommitMessage
	json.NewDecoder(strings.NewReader(commit.Content)).Decode(&commitMessage)
	return &commitMessage, nil
}

func createCommitTypeOption(name string, desc string, ty CommitType, theme lipgloss.Style) huh.Option[CommitType] {
	return huh.NewOption(fmt.Sprintf("%8s: %s", name, theme.Render(desc)), ty)
}

func NewCommitMessage(theme *huh.Theme, useLLM bool) (*CommitMessage, error) {
	commitMsg := &CommitMessage{}
	if useLLM {
		summaries, err := getDiffSummaries()
		if err != nil {
			return nil, err
		}

		if len(summaries) != 0 {
			commitMsg, err = getCommitDetails(summaries)
			if err != nil {
				return nil, err
			}
		}
	}

	descriptionStyle := theme.Help.FullDesc

	// TODOs:
	// * create full pipeline for Git functionality
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
			huh.NewInput().
				Title("Breaking Change (optional)?").
				Prompt(" ").
				Inline(true).
				Value(&commitMsg.breakingChange),
		),
	).WithTheme(theme)

	err := form.Run()
	return commitMsg, err
}
