package ibt

import (
	"os"
	"testing"

	"github.com/teamjorge/ibt/headers"
)

func BenchmarkParserNext(b *testing.B) {
	f, err := os.Open(".testing/valid_test_file.ibt")
	if err != nil {
		b.Fatalf("failed to open testing file - %v", err)
	}
	defer f.Close()

	testHeaders, err := headers.ParseHeaders(f)
	if err != nil {
		b.Fatalf("failed to parse header for testing file - %v", err)
	}

	b.Run("single_field", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			p := NewParser(f, testHeaders, "LapCurrentLapTime")
			for {
				_, hasNext := p.Next()
				if !hasNext {
					break
				}
			}
		}
	})

	b.Run("multiple_fields", func(b *testing.B) {
		fields := []string{"LapCurrentLapTime", "Speed", "RPM", "Gear", "Throttle", "Brake"}
		for i := 0; i < b.N; i++ {
			p := NewParser(f, testHeaders, fields...)
			for {
				_, hasNext := p.Next()
				if !hasNext {
					break
				}
			}
		}
	})

	b.Run("all_fields", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			p := NewParser(f, testHeaders, "*")
			for {
				_, hasNext := p.Next()
				if !hasNext {
					break
				}
			}
		}
	})
}

func BenchmarkReadVarValue(b *testing.B) {
	testData := make([]byte, 1024)
	for i := range testData {
		testData[i] = byte(i % 256)
	}

	vh := headers.VarHeader{
		Offset: 0,
		Count:  1,
		Rtype:  4, // float32
	}

	b.Run("single_float32_original", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = readVarValue(testData, vh)
		}
	})

	b.Run("single_float32_fast", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = readVarValueFast(testData, vh)
		}
	})

	vh.Count = 16
	b.Run("array_float32_original", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = readVarValue(testData, vh)
		}
	})

	b.Run("array_float32_fast", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = readVarValueFast(testData, vh)
		}
	})
}

func BenchmarkTickFilter(b *testing.B) {
	tick := Tick{
		"LapCurrentLapTime": float32(37.5),
		"Speed":             float32(120.5),
		"RPM":               float32(7500),
		"Gear":              int(4),
		"Throttle":          float32(0.85),
		"Brake":             float32(0.0),
		"SteeringWheelAngle": float32(-0.15),
		"LapDist":           float32(1500.5),
	}

	whitelist := []string{"LapCurrentLapTime", "Speed", "RPM"}

	b.Run("filter_subset", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = tick.Filter(whitelist...)
		}
	})

	b.Run("filter_all", func(b *testing.B) {
		allFields := make([]string, 0, len(tick))
		for k := range tick {
			allFields = append(allFields, k)
		}
		for i := 0; i < b.N; i++ {
			_ = tick.Filter(allFields...)
		}
	})
}

func BenchmarkZeroCopyParser(b *testing.B) {
	f, err := os.Open(".testing/valid_test_file.ibt")
	if err != nil {
		b.Fatalf("failed to open testing file - %v", err)
	}
	defer f.Close()

	testHeaders, err := headers.ParseHeaders(f)
	if err != nil {
		b.Fatalf("failed to parse header for testing file - %v", err)
	}

	b.Run("zero_copy_single_field", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			p := NewZeroCopyParser(f, testHeaders, "LapCurrentLapTime")
			for {
				_, hasNext := p.NextZeroCopy()
				if !hasNext {
					break
				}
			}
		}
	})

	b.Run("zero_copy_multiple_fields", func(b *testing.B) {
		fields := []string{"LapCurrentLapTime", "Speed", "RPM", "Gear", "Throttle", "Brake"}
		for i := 0; i < b.N; i++ {
			p := NewZeroCopyParser(f, testHeaders, fields...)
			for {
				_, hasNext := p.NextZeroCopy()
				if !hasNext {
					break
				}
			}
		}
	})
}