package mdgo

import (
	. "github.com/cdvelop/tinystring"
)

// writeIfDifferent writes content to file only if it doesn't exist or content is different
func (m *Mdgo) writeIfDifferent(filePath, content string) error {
	// Check if file exists
	existingContent, err := m.readFile(filePath)
	if err == nil {
		// File exists, check if content is different
		if string(existingContent) == content {
			m.logger("File", filePath, "already up to date, skipping write")
			return nil
		}
	}

	// Write the file
	if err := m.writeFile(filePath, []byte(content)); err != nil {
		return Errf("writing file: %v", err)
	}

	m.logger("Written file", filePath)

	return nil
}
