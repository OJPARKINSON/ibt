package utilities

import (
	"unsafe"
)

// FastByte4ToFloat uses unsafe pointer arithmetic for zero-copy conversion
// This is significantly faster than binary.LittleEndian + math.Float32frombits
func FastByte4ToFloat(in []byte) float32 {
	if len(in) < 4 {
		return 0
	}
	// Cast the byte slice directly to float32 using unsafe
	return *(*float32)(unsafe.Pointer(&in[0]))
}

// FastByte4ToInt uses unsafe pointer arithmetic for zero-copy conversion
func FastByte4ToInt(in []byte) int {
	if len(in) < 4 {
		return 0
	}
	converted := *(*uint32)(unsafe.Pointer(&in[0]))
	if converted == 4294967295 {
		return -1
	}
	return int(converted)
}

// FastByte8ToFloat uses unsafe pointer arithmetic for zero-copy conversion  
func FastByte8ToFloat(in []byte) float64 {
	if len(in) < 8 {
		return 0
	}
	return *(*float64)(unsafe.Pointer(&in[0]))
}

// FastByte8ToInt64 uses unsafe pointer arithmetic for zero-copy conversion
func FastByte8ToInt64(in []byte) int64 {
	if len(in) < 8 {
		return 0
	}
	return *(*int64)(unsafe.Pointer(&in[0]))
}