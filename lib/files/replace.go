package files

import (
	"bytes"
	"fmt"
	"os"
	"strings"
)

func StrictReplaceString(from, to string) func(string) (string, bool) {
	return func(s string) (string, bool) {
		if !strings.Contains(s, from) {
			return s, false
		}
		return strings.ReplaceAll(s, from, to), true
	}
}

func JoinOperations(fs ...func(string) (string, bool)) func(string) (string, bool) {
	return func(s string) (string, bool) {
		anyOk := false
		line := s
		for _, f := range fs {
			replaced, ok := f(line)
			if ok {
				line = replaced
				anyOk = true
			}
		}
		return line, anyOk
	}
}

func ReplaceLines(p string, f func(string) (string, bool)) (bool, error) {
	bb, err := os.ReadFile(p)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, fmt.Errorf("failed to read file: %w", err)
	}
	lines := strings.Split(string(bb), "\n")
	result := bytes.Buffer{}

	hadAtLeastOneOK := false

	for i, l := range lines {
		if i != 0 {
			result.WriteString("\n")
		}

		r, ok := f(l)
		if ok {
			result.WriteString(r)
		} else {
			result.WriteString(l)
		}
		hadAtLeastOneOK = hadAtLeastOneOK || ok
	}

	if !hadAtLeastOneOK {
		return false, nil
	}

	if err := os.WriteFile(p, result.Bytes(), 0777); err != nil {
		return false, fmt.Errorf("failed to write file %q: %w", p, err)
	}

	return true, nil
}
