package ibt

import (
	"encoding/binary"
	"math"
	"testing"
	"unsafe"
)

// Benchmark current unsafe implementation
func BenchmarkUnsafeFloat32(b *testing.B) {
	buf := []byte{0x00, 0x00, 0x80, 0x3F} // 1.0 as float32
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = *(*float32)(unsafe.Pointer(&buf[0]))
	}
}

// Benchmark safe alternative
func BenchmarkSafeFloat32(b *testing.B) {
	buf := []byte{0x00, 0x00, 0x80, 0x3F} // 1.0 as float32
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		bits := binary.LittleEndian.Uint32(buf)
		_ = math.Float32frombits(bits)
	}
}

// Benchmark current unsafe int32
func BenchmarkUnsafeInt32(b *testing.B) {
	buf := []byte{0x01, 0x00, 0x00, 0x00} // 1 as int32
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = *(*int32)(unsafe.Pointer(&buf[0]))
	}
}

// Benchmark safe alternative
func BenchmarkSafeInt32(b *testing.B) {
	buf := []byte{0x01, 0x00, 0x00, 0x00} // 1 as int32
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = int32(binary.LittleEndian.Uint32(buf))
	}
}

// Benchmark realistic workload: processing 1 complete tick (39 fields)
func BenchmarkFullTickUnsafe(b *testing.B) {
	// Simulated buffer with 39 float32 fields
	buf := make([]byte, 39*4)
	for i := range buf {
		buf[i] = byte(i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for offset := 0; offset < len(buf); offset += 4 {
			_ = *(*float32)(unsafe.Pointer(&buf[offset]))
		}
	}
}

func BenchmarkFullTickSafe(b *testing.B) {
	// Simulated buffer with 39 float32 fields
	buf := make([]byte, 39*4)
	for i := range buf {
		buf[i] = byte(i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for offset := 0; offset < len(buf); offset += 4 {
			bits := binary.LittleEndian.Uint32(buf[offset : offset+4])
			_ = math.Float32frombits(bits)
		}
	}
}
