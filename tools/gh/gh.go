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

func (cmd CmdPR) Create(msg, body string) (string, error) {
	std, _, err := cmd.WithArgs("create", "--fill", "--title", msg, "--body", body).WithPrintStdErr().Run()
	return std, err
}

func (cmd CmdPR) Edit(msg, body string) (string, error) {
	std, _, err := cmd.WithArgs("edit", "--title", msg, "--body", body).WithPrintStdErr().Run()
	return std, err
}

func CreatePR(p, msg, body string) (string, error) {
	ghSTD, err := createPR(p, msg, body)
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

func createPR(p, msg, body string) (string, error) {
	gh := gh.InDirectory(p)

	stdout, err := gh.PR().Create(msg, body)
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
	return editStdout, nil
}
