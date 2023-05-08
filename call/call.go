package call

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type Cmd struct {
	Name        string
	Args        []string
	Dir         string
	PrintStderr bool
}

func Command(name string) Cmd {
	return Cmd{
		Name: name,
	}
}

func (cmd Cmd) WithArgs(args ...string) Cmd {
	oldArgs := make([]string, len(cmd.Args))
	copy(oldArgs, cmd.Args)
	cmd.Args = append(oldArgs, args...)
	return cmd
}

func (cmd Cmd) InDirectory(directory string) Cmd {
	cmd.Dir = directory
	return cmd
}

func (cmd Cmd) String() string {
	return fmt.Sprintf("%s %s", cmd.Name, strings.Join(cmd.Args, " "))
}

func (cmd Cmd) WithPrintStdErr() Cmd {
	cmd.PrintStderr = true
	return cmd
}

func (cmd Cmd) Run() (string, string, error) {
	c := exec.Command(cmd.Name, cmd.Args...)
	if cmd.Dir != "" {
		c.Dir = cmd.Dir
	}
	bbErr := bytes.Buffer{}
	bbStd := bytes.Buffer{}

	c.Stderr = &bbErr
	c.Stdout = &bbStd

	err := c.Run()
	if cmd.PrintStderr {
		os.Stderr.WriteString(bbErr.String())
	}
	if err != nil {
		return bbStd.String(), bbErr.String(), fmt.Errorf("failed to run %s: %w", cmd, err)
	}
	return bbStd.String(), bbErr.String(), nil
}

func (cmd Cmd) SimpleRun() error {
	_, _, err := cmd.WithPrintStdErr().Run()
	return err
}
