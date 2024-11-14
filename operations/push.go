package operations

import (
	"fmt"

	"github.com/krzysztofdrys/fer/tools/gh"
	"github.com/krzysztofdrys/fer/tools/git"
)

func Push(p, branch, title, body string, draftPR bool) (string, error) {
	if err := git.PushForceChanges(p, branch); err != nil {
		return "", fmt.Errorf("failed to push changes: %w", err)
	}
	pr, err := gh.EnsureThereIsPR(p, title, body, draftPR)
	if err != nil {
		return "", fmt.Errorf("failed to create PR: %w", err)
	}
	return pr, nil
}
