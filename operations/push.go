package operations

import (
	"fmt"

	"github.com/krzysztofdrys/fer/tools/gh"
	"github.com/krzysztofdrys/fer/tools/git"
)

func Push(p, branch, title, body string) (string, error) {
	if err := git.PushChanges(p, branch); err != nil {
		return "", fmt.Errorf("failed to push changes: %w", err)
	}
	pr, err := gh.CreatePR(p, title, body)
	if err != nil {
		return "", fmt.Errorf("failed to create PR: %w", err)
	}
	return pr, nil
}
