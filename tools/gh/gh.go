package gh

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"github.com/krzysztofdrys/fer/call"
)

var gh = Cmd{call.Command("gh")}

type Cmd struct {
	call.Cmd
}

func (cmd Cmd) InDirectory(p string) Cmd {
	return Cmd{cmd.Cmd.InDirectory(p)}
}

func (cmd Cmd) PR() CmdPR {
	return CmdPR{cmd.WithArgs("pr")}
}

type CmdPR struct {
	call.Cmd
}

func (cmd CmdPR) Create(msg, body string, draftPR bool) (string, error) {
	args := []string{"create", "--fill", "--title", msg, "--body", body}
	if draftPR {
		args = append(args, "--draft")
	}
	std, _, err := cmd.WithArgs(args...).WithPrintStdErr().Run()
	return std, err
}

func (cmd CmdPR) Edit(msg, body string) (string, error) {
	args := []string{"edit", "--title", msg, "--body", body}

	std, _, err := cmd.WithArgs(args...).WithPrintStdErr().Run()
	return std, err
}

func (cmd CmdPR) Ready() error {
	args := []string{"ready"}

	_, _, err := cmd.WithArgs(args...).WithPrintStdErr().Run()
	return err
}

func (cmd CmdPR) Draft() error {
	args := []string{"ready", "--undo"}

	_, _, err := cmd.WithArgs(args...).WithPrintStdErr().Run()
	return err
}

func EnsureThereIsPR(p, msg, body string, draftPR bool) (string, error) {
	ghSTD, err := ensureThereIsPR(p, msg, body, draftPR)
	if err != nil {
		return "", fmt.Errorf("failed to push changes: %w", err)
	}
	lines := strings.Split(ghSTD, "\n")
	for _, l := range lines {
		if strings.HasPrefix(l, "https") {
			return l, nil
		}
	}

	return "", fmt.Errorf("couldn't find PR link in gh output:\n%s", ghSTD)
}

func DefaultBranch(owner, repository string) (string, error) {
	out, _, err := gh.WithArgs("api", fmt.Sprintf("/repos/%s/%s", owner, repository), "--jq", ".default_branch").WithPrintStdErr().Run()
	if err != nil {
		return "", fmt.Errorf("failed to call github API %w", err)
	}
	return strings.TrimSpace(out), nil
}

func ensureThereIsPR(p, msg, body string, draftPR bool) (string, error) {
	gh := gh.InDirectory(p)

	stdout, err := gh.PR().Create(msg, body, draftPR)
	if err == nil {
		return stdout, nil
	}
	var exitError *exec.ExitError
	if !errors.As(err, &exitError) {
		return "", fmt.Errorf("failed to create PR: %w", err)
	}
	if exitError.ExitCode() != 1 {
		return "", fmt.Errorf("failed to create PR: %w", err)
	}

	editStdout, err := gh.PR().Edit(msg, body)
	if err != nil {
		return "", fmt.Errorf("failed to edit PR: %w", err)
	}
	if draftPR {
		if err := gh.PR().Draft(); err != nil {
			return "", fmt.Errorf("failed to mark PR as draft: %w", err)
		}
	} else {
		if err := gh.PR().Ready(); err != nil {
			return "", fmt.Errorf("failed to mark PR as ready: %w", err)
		}
	}
	return editStdout, nil
}
