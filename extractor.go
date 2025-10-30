package mdgo

import (
	. "github.com/cdvelop/tinystring"
	"regexp"
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
			blocks = append(blocks, Convert(match[1]).TrimSpace().String())
		}
	}

	// Join all blocks with double newline
	return Convert(blocks).Join("\n\n").String()
}
