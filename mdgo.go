package mdgo

import (
	"os"
	"path/filepath"

	. "github.com/cdvelop/tinystring"
)

type Mdgo struct {
	rootDir     string
	destination string
	logger      func(...any)
	// input sources (one of these should be set before calling Extract)
	inputPath      string
	inputBytes     []byte
	inputEmbed     ReaderFile
	inputEmbedPath string
}

// ReaderFile abstracts a file reader (e.g. embed.FS or any custom provider)
type ReaderFile interface {
	ReadFile(name string) ([]byte, error)
}

// New creates a new Mdgo instance with the root directory.
// Destination (output directory) and input must be set via methods.
func New(rootDir, destination string) *Mdgo {
	return &Mdgo{
		rootDir:     rootDir,
		destination: destination,
		logger:      nil,
	}
}

// SetLogger sets a custom logger function
func (m *Mdgo) SetLogger(logger func(...any)) *Mdgo {
	m.logger = logger
	return m
}

// InputPath sets the input as a file path (relative to rootDir)
func (m *Mdgo) InputPath(pathFile string) *Mdgo {
	m.inputPath = pathFile
	// clear other inputs
	m.inputBytes = nil
	m.inputEmbed = nil
	m.inputEmbedPath = ""
	return m
}

// InputByte sets the input as a byte slice (markdown content)
func (m *Mdgo) InputByte(content []byte) *Mdgo {
	m.inputBytes = content
	// clear other inputs
	m.inputPath = ""
	m.inputEmbed = nil
	m.inputEmbedPath = ""
	return m
}

// InputEmbed sets the input as any ReaderFile implementation and a relative path inside it
func (m *Mdgo) InputEmbed(r ReaderFile, path string) *Mdgo {
	m.inputEmbed = r
	m.inputEmbedPath = path
	// clear other inputs
	m.inputPath = ""
	m.inputBytes = nil
	return m
}

// Extract extracts code blocks from the configured input and writes to outputFile
// The output file extension determines which code type to extract (.go, .js, .css)
func (m *Mdgo) Extract(outputFile string) error {
	if m.destination == "" {
		return Errf("destination not set; provide destination when calling New(rootDir, destination)")
	}

	// Read markdown from the configured input
	markdown, err := m.readConfiguredSource()
	if err != nil {
		return Errf("reading source: %w", err)
	}

	// Determine code type from output file extension
	codeType := m.getCodeType(outputFile)
	if codeType == "" {
		return Errf("unsupported file extension: %s", filepath.Ext(outputFile))
	}

	// Extract code blocks
	code := m.extractCodeBlocks(markdown, codeType)
	if code == "" {
		return Errf("no %s code blocks found in markdown", codeType)
	}

	// Write to output file
	outputPath := filepath.Join(m.destination, outputFile)
	if err := m.writeIfDifferent(outputPath, code); err != nil {
		return Errf("writing output file: %w", err)
	}

	if m.logger != nil {
		m.logger("Extracted", codeType, "code to", outputPath)
	}

	return nil
}

// readConfiguredSource reads markdown content from the configured input
func (m *Mdgo) readConfiguredSource() (string, error) {
	// prioritize explicit byte input
	if m.inputBytes != nil {
		return string(m.inputBytes), nil
	}
	if m.inputPath != "" {
		fullPath := filepath.Join(m.rootDir, m.inputPath)
		data, err := os.ReadFile(fullPath)
		if err != nil {
			return "", Errf("reading file %s: %w", fullPath, err)
		}
		return string(data), nil
	}
	if m.inputEmbed != nil {
		b, err := m.inputEmbed.ReadFile(m.inputEmbedPath)
		if err != nil {
			return "", Errf("reading embedded file %s: %w", m.inputEmbedPath, err)
		}
		return string(b), nil
	}
	return "", Errf("no input configured; call InputPath, InputByte, or InputEmbed before Extract")
}

// getCodeType determines the code type from file extension
func (m *Mdgo) getCodeType(outputFile string) string {
	ext := Convert(outputFile).PathExt().String()
	switch ext {
	case ".go":
		return "go"
	case ".js":
		return "javascript"
	case ".css":
		return "css"
	default:
		return ""
	}
}
