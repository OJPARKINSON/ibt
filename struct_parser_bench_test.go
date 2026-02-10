package ibt

import (
	"fmt"
	"testing"

	"github.com/OJPARKINSON/ibt/headers"
)

// setupBenchParser creates a DirectStructParser from the test file with an optional whitelist.
func setupBenchParser(b *testing.B, whitelist ...string) (*DirectStructParser, func()) {
	b.Helper()

	reader, err := NewIbtReader(".testing/valid_test_file.ibt")
	if err != nil {
		b.Fatalf("failed to open test file: %v", err)
	}

	header, err := headers.ParseHeaders(reader)
	if err != nil {
		reader.Close()
		b.Fatalf("failed to parse headers: %v", err)
	}

	parser := NewDirectStructParser(reader, header, whitelist...)
	cleanup := func() { reader.Close() }
	return parser, cleanup
}

// smallWhitelist is a typical processor usage with ~10 fields.
var smallWhitelist = []string{
	"Lap", "LapDistPct", "Speed", "RPM", "Throttle", "Brake",
	"Gear", "SteeringWheelAngle", "LapCurrentLapTime", "SessionTime",
}

// BenchmarkNextStruct benchmarks the hot path: DirectStructParser.NextStruct().
func BenchmarkNextStruct(b *testing.B) {
	b.Run("whitelist_10_fields", func(b *testing.B) {
		parser, cleanup := setupBenchParser(b, smallWhitelist...)
		defer cleanup()

		tick := &TelemetryTick{}
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, hasNext := parser.NextStruct(tick)
			if !hasNext {
				// Reset to beginning for continuous benchmarking
				parser.current = 0
			}
		}
	})

	b.Run("all_fields", func(b *testing.B) {
		parser, cleanup := setupBenchParser(b)
		defer cleanup()

		tick := &TelemetryTick{}
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, hasNext := parser.NextStruct(tick)
			if !hasNext {
				parser.current = 0
			}
		}
	})

	b.Run("tick_reuse_zero_alloc", func(b *testing.B) {
		parser, cleanup := setupBenchParser(b, smallWhitelist...)
		defer cleanup()

		// Pre-allocate a single tick and reuse it every iteration
		tick := &TelemetryTick{}
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			*tick = TelemetryTick{} // zero out without allocating
			_, hasNext := parser.NextStruct(tick)
			if !hasNext {
				parser.current = 0
			}
		}
	})
}

// BenchmarkNewDirectStructParser benchmarks parser construction cost.
func BenchmarkNewDirectStructParser(b *testing.B) {
	reader, err := NewIbtReader(".testing/valid_test_file.ibt")
	if err != nil {
		b.Fatalf("failed to open test file: %v", err)
	}
	defer reader.Close()

	header, err := headers.ParseHeaders(reader)
	if err != nil {
		b.Fatalf("failed to parse headers: %v", err)
	}

	// Warm up the global field setters (sync.Once)
	_ = getFieldSetters()

	b.Run("no_whitelist", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = NewDirectStructParser(reader, header)
		}
	})

	b.Run("with_whitelist_10", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = NewDirectStructParser(reader, header, smallWhitelist...)
		}
	})
}

// BenchmarkEndToEnd benchmarks a full parse loop: create parser → iterate all ticks.
func BenchmarkEndToEnd(b *testing.B) {
	reader, err := NewIbtReader(".testing/valid_test_file.ibt")
	if err != nil {
		b.Fatalf("failed to open test file: %v", err)
	}
	defer reader.Close()

	header, err := headers.ParseHeaders(reader)
	if err != nil {
		b.Fatalf("failed to parse headers: %v", err)
	}

	// Count ticks for throughput reporting
	countParser := NewDirectStructParser(reader, header)
	tick := &TelemetryTick{}
	numTicks := 0
	for {
		_, hasNext := countParser.NextStruct(tick)
		numTicks++
		if !hasNext {
			break
		}
	}

	dataSize := int64(header.TelemetryHeader.BufLen * numTicks)

	b.Run(fmt.Sprintf("all_fields_%d_ticks", numTicks), func(b *testing.B) {
		b.SetBytes(dataSize)
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			parser := NewDirectStructParser(reader, header)
			t := &TelemetryTick{}
			for {
				_, hasNext := parser.NextStruct(t)
				if !hasNext {
					break
				}
			}
		}
	})

	b.Run(fmt.Sprintf("whitelist_10_%d_ticks", numTicks), func(b *testing.B) {
		b.SetBytes(dataSize)
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			parser := NewDirectStructParser(reader, header, smallWhitelist...)
			t := &TelemetryTick{}
			for {
				_, hasNext := parser.NextStruct(t)
				if !hasNext {
					break
				}
			}
		}
	})
}
