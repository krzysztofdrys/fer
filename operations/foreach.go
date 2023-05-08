package operations

import (
	"fmt"
	"log"
	"os"
)

type Repository struct {
	Repository string
	Directory  string
	MainBranch string
}

type ForeachConfig struct {
	GitCache       string
	CheckoutBranch string
	AuthorEmail    string

	PRTitle string
	PRBody  string

	Repositories []Repository
}

type Result struct {
	Error error
	PR    string
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
	getCfg := GetConfig{
		GitCache:       cfg.GitCache,
		MainBranch:     r.MainBranch,
		Repository:     r.Repository,
		Directory:      r.Directory,
		CheckoutBranch: cfg.CheckoutBranch,
		AuthorEmail:    cfg.AuthorEmail,
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

	pr, err := Push(p, cfg.CheckoutBranch, cfg.PRTitle, cfg.PRBody)
	if err != nil {
		return "", fmt.Errorf("failed to push %q: %w", r.Repository, err)
	}

	return pr, nil
}
