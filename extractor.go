package mdgo

import (
	"regexp"
	"strings"
)

// extractCodeBlocks extracts code blocks of a specific type from markdown
func (m *Mdgo) extractCodeBlocks(markdown, codeType string) string {
	// Pattern to capture code blocks: ```<codeType>\n...\n```
	// Using DOTALL mode (?s) to match across newlines
	pattern := "(?s)```" + codeType + "\\n(.*?)```"
	re := regexp.MustCompile(pattern)
	matches := re.FindAllStringSubmatch(markdown, -1)

	var blocks []string
	for _, match := range matches {
		if len(match) > 1 {
			blocks = append(blocks, strings.TrimSpace(match[1]))
		}
	}

	// Join all blocks with double newline
	return strings.Join(blocks, "\n\n")
}
