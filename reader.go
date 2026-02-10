package ibt

import (
	"io"
	"os"
)

type IbtReader struct {
	data []byte
}

// NewIbtReader reads the entire file into memory for parsing.
func NewIbtReader(filename string) (*IbtReader, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return &IbtReader{data: data}, nil
}

// Read implements the io.Reader interface
func (m *IbtReader) Read(p []byte) (int, error) {
	// For header parsing compatibility - reads from beginning
	if len(m.data) == 0 {
		return 0, io.EOF
	}
	n := copy(p, m.data)
	if n < len(p) {
		return n, io.EOF
	}
	return n, nil
}

// ReadAt implements the io.ReaderAt interface with zero-copy reads
func (m *IbtReader) ReadAt(p []byte, off int64) (int, error) {
	if off < 0 || off >= int64(len(m.data)) {
		return 0, io.EOF
	}

	n := copy(p, m.data[off:])
	return n, nil
}

// Close releases the data held by the reader.
func (m *IbtReader) Close() error {
	if m == nil {
		return nil
	}

	// Check if already closed (both data and file are nil)
	if m.data == nil {
		return os.ErrClosed
	}

	m.data = nil
	return nil
}

// Data returns the underlying byte slice for zero-copy parsing.
func (r *IbtReader) Data() []byte {
	return r.data
}
