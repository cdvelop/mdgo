package mdgo

import (
	. "github.com/cdvelop/tinystring"
)

type Mdgo struct {
	rootDir     string
	destination string
	// input sources (one of these should be set before calling Extract)
	inputPath string

	readFile  func(name string) ([]byte, error)
	writeFile func(name string, data []byte) error
	logger    func(...any)
}

// New creates a new Mdgo instance with the root directory.
// Destination (output directory) and input must be set via methods.
func New(rootDir, destination string, writerFile func(name string, data []byte) error) *Mdgo {
	return &Mdgo{
		rootDir:     rootDir,
		destination: destination,
		readFile:    func(name string) ([]byte, error) { return nil, Err("not configure reader func") },
		writeFile:   writerFile,
		logger:      func(...any) {},
	}
}

// SetLogger sets a custom logger function
func (m *Mdgo) SetLogger(logger func(...any)) *Mdgo {
	m.logger = logger
	return m
}

// InputPath sets the input as a file path (relative to rootDir)
func (m *Mdgo) InputPath(pathFile string, readerFile func(name string) ([]byte, error)) *Mdgo {
	m.inputPath = pathFile
	m.readFile = readerFile
	return m
}

// InputByte sets the input as a byte slice (markdown content)
func (m *Mdgo) InputByte(content []byte) *Mdgo {
	// clear other inputs
	m.readFile = func(name string) ([]byte, error) {
		return content, nil
	}

	return m
}

// InputEmbed sets the input as any ReaderFile implementation and a relative path inside it
func (m *Mdgo) InputEmbed(path string, readerFile func(name string) ([]byte, error)) *Mdgo {
	m.readFile = readerFile
	// clear other inputs
	m.inputPath = path
	return m
}

// Extract extracts code blocks from the configured input and writes to outputFile
// The output file extension determines which code type to extract (.go, .js, .css)
func (m *Mdgo) Extract(outputFile string) error {
	if m.destination == "" {
		return Errf("destination not set; provide destination when calling New(rootDir, destination)")
	}

	// Read markdown from the configured input
	markdown, err := m.readFile(m.inputPath)
	if err != nil {
		return Errf("reading file %s: %v", m.inputPath, err)
	}

	// Determine code type from output file extension
	codeType := m.getCodeType(outputFile)
	if codeType == "" {
		return Errf("unsupported file extension: %s", Convert(outputFile).PathExt().String())
	}

	// Extract code blocks
	code := m.extractCodeBlocks(string(markdown), codeType)
	if code == "" {
		return Errf("no %s code blocks found in markdown", codeType)
	}

	// Write to output file
	outputPath := PathJoin(m.destination, outputFile).String()
	if err := m.writeIfDifferent(outputPath, code); err != nil {
		return Errf("writing output file: %v", err)
	}

	if m.logger != nil {
		m.logger("Extracted", codeType, "code to", outputPath)
	}

	return nil
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
