package operations

import (
	"fmt"
	"log"
	"os"

	"github.com/krzysztofdrys/fer/tools/git"

	"github.com/krzysztofdrys/fer/tools/gh"
)

type Repository struct {
	Owner      string
	Repository string
	Directory  string
}

func (r Repository) URL() string {
	return fmt.Sprintf("https://github.com/%s/%s", r.Owner, r.Repository)
}

type ForeachConfig struct {
	GitCache       string
	CheckoutBranch string
	AuthorEmail    string

	PRTitle string
	PRBody  string
	DraftPR bool

	Repositories []Repository
}

func UpdateAll(cfg ForeachConfig, f func(p string) (bool, error)) ([]string, error) {
	result := []string{}

	if err := os.MkdirAll(cfg.GitCache, 0777); err != nil {
		return nil, fmt.Errorf("failed to create git cache directory %q: %w", cfg.GitCache, err)
	}

	for _, r := range cfg.Repositories {
		pr, err := run(cfg, r, f)
		if err != nil {
			log.Printf("Failed to process %q: %v", r.Repository, err)
		} else if pr != "" {
			result = append(result, pr)
		}
	}

	return result, nil
}

func run(cfg ForeachConfig, r Repository, f func(p string) (bool, error)) (string, error) {
	defaultBranch, err := gh.DefaultBranch(r.Owner, r.Repository)
	if err != nil {
		return "", fmt.Errorf("failed to get default branch for %q: %w", r.Repository, err)
	}

	getCfg := GetConfig{
		GitCache:       cfg.GitCache,
		RepositoryURL:  r.URL(),
		Directory:      r.Directory,
		CheckoutBranch: cfg.CheckoutBranch,
		AuthorEmail:    cfg.AuthorEmail,
		DefaultBranch:  defaultBranch,
	}
	p, err := Get(getCfg)
	if err != nil {
		return "", fmt.Errorf("failed to get %q: %w", r.Repository, err)
	}

	ok, err := f(p)
	if err != nil {
		return "", fmt.Errorf("failed to process %q: %w", r.Repository, err)
	}

	if !ok {
		return "", nil
	}

	ok, err = git.CommitChanges(p, fmt.Sprintf("%s\n\n%s", cfg.PRTitle, cfg.PRBody))
	if err != nil {
		return "", fmt.Errorf("failed to commit changes: %w", err)
	}
	if !ok {
		return "", nil
	}

	pr, err := Push(p, cfg.CheckoutBranch, cfg.PRTitle, cfg.PRBody, cfg.DraftPR)
	if err != nil {
		return "", fmt.Errorf("failed to push %q: %w", r.Repository, err)
	}

	return pr, nil
}
