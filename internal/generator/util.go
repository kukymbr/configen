package generator

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/kukymbr/configen/internal/logger"
)

const (
	dirsMode  os.FileMode = 0755
	filesMode os.FileMode = 0644
)

var pxIdentifier = regexp.MustCompile(`(?i)^[a-z]+[a-z0-9_]*$`)

func validateIdentifier(name string) error {
	if strings.TrimSpace(name) == "" {
		return errors.New("identifier cannot be empty")
	}

	if !pxIdentifier.MatchString(name) {
		return fmt.Errorf("'%s' is not a valid identifier", name)
	}

	return nil
}

// validateIsDir checks if path exists and is a directory.
func validateIsDir(path string) error {
	stat, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("directory '%s': %w", path, err)
	}

	if !stat.IsDir() {
		return fmt.Errorf("'%s' is not a directory", path)
	}

	return nil
}

// EnsureDir creates dir if not exists.
func EnsureDir(path string) error {
	if err := os.MkdirAll(path, dirsMode); err != nil {
		return fmt.Errorf("dir '%s' does not exist and cannot be created: %w", path, err)
	}

	return nil
}

func writeFile(content []byte, target string) error {
	if err := os.WriteFile(target, content, filesMode); err != nil {
		return fmt.Errorf("failed to write file %s: %w", target, err)
	}

	logger.Successf("Generated %s file", target)

	return nil
}
