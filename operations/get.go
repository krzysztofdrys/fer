package operations

import (
	"fmt"
	"path/filepath"

	"github.com/krzysztofdrys/fer/tools/git"
)

type GetConfig struct {
	GitCache       string
	DefaultBranch  string
	RepositoryURL  string
	Directory      string
	CheckoutBranch string
	AuthorEmail    string
}

func Get(cfg GetConfig) (string, error) {
	err := git.EnsureFreshRepo(
		cfg.GitCache,
		cfg.Directory,
		cfg.RepositoryURL,
		cfg.DefaultBranch)
	if err != nil {
		return "", fmt.Errorf("failed to reset repository to fresh %q branch: %w", cfg.DefaultBranch, err)
	}

	p := filepath.Join(cfg.GitCache, cfg.Directory)
	if err := git.SwitchToNewBranch(p, cfg.CheckoutBranch); err != nil {
		return "", fmt.Errorf("failed to checkout %q branch: %w", cfg.CheckoutBranch, err)
	}

	if err := git.EnsureGitAuthor(p, cfg.AuthorEmail); err != nil {
		return "", fmt.Errorf("failed to set git author: %w", err)
	}

	return p, nil
}
