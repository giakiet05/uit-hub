package localtool

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"
)

func safeSandboxPath(rootDir string, relativePath string) (string, error) {
	if strings.TrimSpace(relativePath) == "" {
		return "", errors.New("path must not be empty")
	}
	if filepath.IsAbs(relativePath) {
		return "", errors.New("path must be relative")
	}

	cleanPath := filepath.Clean(relativePath)
	if cleanPath == "." || strings.HasPrefix(cleanPath, ".."+string(filepath.Separator)) || cleanPath == ".." {
		return "", errors.New("path escapes sandbox")
	}

	if rootDir == "" {
		rootDir = filepath.Join("tmp", "agent-files")
	}
	rootDir, err := filepath.Abs(rootDir)
	if err != nil {
		return "", fmt.Errorf("resolve sandbox root: %w", err)
	}

	fullPath := filepath.Join(rootDir, cleanPath)
	if !strings.HasPrefix(fullPath, rootDir+string(filepath.Separator)) && fullPath != rootDir {
		return "", errors.New("path escapes sandbox")
	}
	return fullPath, nil
}
