package mdgo

import (
	. "github.com/cdvelop/tinystring"
	"os"
	"path/filepath"
)

// writeIfDifferent writes content to file only if it doesn't exist or content is different
func (m *Mdgo) writeIfDifferent(filePath, content string) error {
	// Check if file exists
	existingContent, err := os.ReadFile(filePath)
	if err == nil {
		// File exists, check if content is different
		if string(existingContent) == content {
			m.logger("File", filePath, "already up to date, skipping write")
			return nil
		}
	}

	// File doesn't exist or content is different, write it
	// Ensure directory exists
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return Errf("creating directory %s: %w", dir, err)
	}

	// Write the file
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		return Errf("writing file: %w", err)
	}

	m.logger("Written file", filePath)

	return nil
}
