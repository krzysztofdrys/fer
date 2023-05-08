package git

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/krzysztofdrys/fer/call"
)

var git = call.Command("git")

func EnsureFreshRepo(cache, directory, repository, branch string) error {
	p := filepath.Join(cache, directory)
	_, err := os.Stat(p)
	if err == nil {
		return reset(p, branch)
	}
	if !os.IsNotExist(err) {
		return fmt.Errorf("failed to stat %q: %w", p, err)
	}
	if err := clone(cache, directory, repository); err != nil {
		return fmt.Errorf("failed to clone %q: %w", repository, err)
	}
	return reset(p, branch)
}

func SwitchToNewBranch(p, branch string) error {
	err := call.
		Command("git").
		WithArgs("checkout", "-B", branch).
		InDirectory(p).
		SimpleRun()
	if err != nil {
		return fmt.Errorf("failed to switch branch to %q: %w", branch, err)
	}
	return nil
}

func CommitChanges(p, msg string) (bool, error) {
	if err := git.WithArgs("add", ".").InDirectory(p).SimpleRun(); err != nil {
		return false, fmt.Errorf("failed to add files to git: %v", err)
	}

	std, _, err := git.WithArgs("status", "--porcelain").InDirectory(p).WithPrintStdErr().Run()
	if err != nil {
		return false, fmt.Errorf("failed to check if there are any git files to commit: %v", err)
	}
	if len(std) == 0 {
		return false, nil
	}

	if err := git.WithArgs("commit", "-m", msg).InDirectory(p).SimpleRun(); err != nil {
		return false, fmt.Errorf("failed to commit changes: %w", err)
	}

	return true, nil
}

func EnsureGitAuthor(p, email string) error {
	err := git.WithArgs("config", "user.email", email).InDirectory(p).SimpleRun()
	if err != nil {
		return fmt.Errorf("failed to set user.email: %w", err)
	}
	return nil
}

func PushChanges(p, branch string) error {
	if err := git.
		WithArgs("push", "--force-with-lease", "origin", branch).
		InDirectory(p).
		SimpleRun(); err != nil {
		return fmt.Errorf("failed to push changes: %w", err)
	}
	return nil
}

func clone(cache, directory, repository string) error {
	err := git.
		WithArgs("clone", repository, directory).
		InDirectory(cache).
		SimpleRun()
	if err != nil {
		return fmt.Errorf("failed to clone repository: %w", err)
	}
	return nil
}

func reset(p, branch string) error {
	err := git.
		WithArgs("reset", "--hard").
		InDirectory(p).
		SimpleRun()
	if err != nil {
		return fmt.Errorf("failed to reset git directory: %v", err)
	}

	err = git.
		WithArgs("checkout", branch).
		InDirectory(p).
		SimpleRun()
	if err != nil {
		return fmt.Errorf("failed to switch branch to %q: %w", branch, err)
	}

	err = call.
		Command("git").
		WithArgs("pull").
		InDirectory(p).
		SimpleRun()
	if err != nil {
		return fmt.Errorf("failed pull newest changes to branch %q: %w", branch, err)
	}
	return nil
}
