package utils

import (
	"fmt"
	"os"
)

const filesMode os.FileMode = 0644

func WriteFile(content []byte, target string) error {
	if err := os.WriteFile(target, content, filesMode); err != nil {
		return fmt.Errorf("failed to write file %s: %w", target, err)
	}

	return nil
}
