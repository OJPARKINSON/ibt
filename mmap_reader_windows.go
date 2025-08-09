//go:build windows

package ibt

import (
	"fmt"
	"io"
	"os"
	"syscall"
	"unsafe"
)

var (
	kernel32 = syscall.NewLazyDLL("kernel32.dll")
	procCreateFileMappingW = kernel32.NewProc("CreateFileMappingW")
	procMapViewOfFile = kernel32.NewProc("MapViewOfFile")
	procUnmapViewOfFile = kernel32.NewProc("UnmapViewOfFile")
	procCloseHandle = kernel32.NewProc("CloseHandle")
)

const (
	FILE_MAP_READ = 0x4
	PAGE_READONLY = 0x02
)

// MmapReader provides memory-mapped file access for Windows
type MmapReader struct {
	data []byte
	file *os.File
	mapping syscall.Handle
	view uintptr
}

// NewMmapReader creates a memory-mapped reader for the given file (Windows implementation)
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

	if stat.Size() == 0 {
		file.Close()
		return &MmapReader{
			data: []byte{},
			file: nil,
		}, nil
	}

	// Create file mapping object
	mapping, _, err := procCreateFileMappingW.Call(
		uintptr(syscall.Handle(file.Fd())),
		0,
		PAGE_READONLY,
		0,
		0,
		0,
	)
	if mapping == 0 {
		file.Close()
		return nil, fmt.Errorf("CreateFileMapping failed: %v", err)
	}

	// Map view of file
	view, _, err := procMapViewOfFile.Call(
		mapping,
		FILE_MAP_READ,
		0,
		0,
		0,
	)
	if view == 0 {
		procCloseHandle.Call(mapping)
		file.Close()
		return nil, fmt.Errorf("MapViewOfFile failed: %v", err)
	}

	// Create Go slice from mapped memory
	data := (*[1 << 30]byte)(unsafe.Pointer(view))[:stat.Size():stat.Size()]

	return &MmapReader{
		data: data,
		file: file,
		mapping: syscall.Handle(mapping),
		view: view,
	}, nil
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

// Close unmaps the file and closes handles
func (m *MmapReader) Close() error {
	var err error

	// Unmap view of file
	if m.view != 0 {
		if ret, _, winErr := procUnmapViewOfFile.Call(m.view); ret == 0 {
			err = fmt.Errorf("UnmapViewOfFile failed: %v", winErr)
		}
		m.view = 0
	}

	// Close mapping handle
	if m.mapping != 0 {
		if ret, _, winErr := procCloseHandle.Call(uintptr(m.mapping)); ret == 0 && err == nil {
			err = fmt.Errorf("CloseHandle (mapping) failed: %v", winErr)
		}
		m.mapping = 0
	}

	// Close file handle
	if m.file != nil {
		if closeErr := m.file.Close(); closeErr != nil && err == nil {
			err = closeErr
		}
		m.file = nil
	}

	m.data = nil
	return err
}

// ReadFrom implements io.WriterTo (required by headers.Reader interface)
func (m *MmapReader) ReadFrom(r interface{}) (int64, error) {
	// Not used in telemetry parsing, but required for interface compliance
	return 0, nil
}