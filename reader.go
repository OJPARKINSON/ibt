//go:build !windows

package ibt

import (
	"io"
	"os"
	"syscall"
	"unsafe"
)

// MmapReader provides memory-mapped file access for faster reading
type MmapReader struct {
	data []byte
	file *os.File
}

// NewMmapReader creates a memory-mapped reader for the given file
func NewMmapReader(filename string) (*MmapReader, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	stat, err := file.Stat()
	if err != nil {
		file.Close()
		return nil, err
	}

	// Handle empty files - can't mmap a 0-byte file
	if stat.Size() == 0 {
		file.Close()
		return &MmapReader{
			data: []byte{},
			file: nil,
		}, nil
	}

	data, err := syscall.Mmap(int(file.Fd()), 0, int(stat.Size()), syscall.PROT_READ, syscall.MAP_PRIVATE)
	if err != nil {
		file.Close()
		return nil, err
	}

	return &MmapReader{
		data: data,
		file: file,
	}, nil
}

// Read implements the io.Reader interface
func (m *MmapReader) Read(p []byte) (int, error) {
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
func (m *MmapReader) ReadAt(p []byte, off int64) (int, error) {
	if off < 0 || off >= int64(len(m.data)) {
		return 0, io.EOF
	}

	n := copy(p, m.data[off:])
	return n, nil
}

// ReadAtUnsafe provides unsafe direct access to memory-mapped data
// WARNING: The returned slice is only valid until the MmapReader is closed
func (m *MmapReader) ReadAtUnsafe(off int64, size int) []byte {
	if off < 0 || off+int64(size) > int64(len(m.data)) {
		return nil
	}

	// Return slice directly from mmap'd memory - zero copy
	return (*[1 << 30]byte)(unsafe.Pointer(&m.data[off]))[:size:size]
}

// Close unmaps the file and closes the file descriptor
func (m *MmapReader) Close() error {
	if m == nil {
		return nil
	}

	// Check if already closed (both data and file are nil)
	if m.data == nil && m.file == nil {
		return os.ErrClosed
	}

	var err error
	// Only Munmap if we actually mapped something (non-empty file)
	if m.data != nil && len(m.data) > 0 {
		err = syscall.Munmap(m.data)
	}
	m.data = nil

	if m.file != nil {
		if closeErr := m.file.Close(); closeErr != nil && err == nil {
			err = closeErr
		}
		m.file = nil
	}
	return err
}

// ReadFrom implements io.WriterTo (required by headers.Reader interface)
func (m *MmapReader) ReadFrom(r interface{}) (int64, error) {
	// Not used in telemetry parsing, but required for interface compliance
	return 0, nil
}

// Data returns the entire memory-mapped data slice
// This enables zero-copy parsing by allowing direct access to the map
func (r *MmapReader) Data() []byte {
	return r.data
}
