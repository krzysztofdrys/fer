package circleci

import (
	"fmt"
	"path/filepath"

	"github.com/krzysztofdrys/fer/call"
)

func Reformat(p string) error {
	configDir := filepath.Join(p, ".circleci")

	err := call.
		Command("yq").
		WithArgs("-i", filepath.Join(configDir, "config.yml")).
		SimpleRun()
	if err != nil {
		return fmt.Errorf("failed to run yq to format config: %w", err)
	}

	if err := CheckConfig(configDir); err != nil {
		return err
	}
	return nil
}

func CheckConfig(p string) error {
	err := call.
		Command("cirleci").
		WithArgs("config", "check").
		InDirectory(p).
		SimpleRun()

	if err != nil {
		return fmt.Errorf("checking if circleci config is correct failed: %w", err)
	}
	return nil
}
