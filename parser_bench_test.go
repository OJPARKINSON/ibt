package ibt

import (
	"testing"

	"github.com/OJPARKINSON/ibt/headers"
)

func BenchmarkParserNext(b *testing.B) {
	reader, err := NewMmapReader(".testing/valid_test_file.ibt")
	if err != nil {
		b.Fatalf("failed to open testing file - %v", err)
	}
	defer reader.Close()

	testHeaders, err := headers.ParseHeaders(reader)
	if err != nil {
		b.Fatalf("failed to parse header for testing file - %v", err)
	}

	b.Run("single_field", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			p := NewParser(reader, testHeaders, "LapCurrentLapTime")
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
			p := NewParser(reader, testHeaders, fields...)
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
			p := NewParser(reader, testHeaders, "*")
			for {
				_, hasNext := p.Next()
				if !hasNext {
					break
				}
			}
		}
	})
}

func BenchmarkTickFilter(b *testing.B) {
	tick := Tick{
		"LapCurrentLapTime":  float32(37.5),
		"Speed":              float32(120.5),
		"RPM":                float32(7500),
		"Gear":               int(4),
		"Throttle":           float32(0.85),
		"Brake":              float32(0.0),
		"SteeringWheelAngle": float32(-0.15),
		"LapDist":            float32(1500.5),
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
	reader, err := NewMmapReader(".testing/valid_test_file.ibt")
	if err != nil {
		b.Fatalf("failed to open testing file - %v", err)
	}
	defer reader.Close()

	testHeaders, err := headers.ParseHeaders(reader)
	if err != nil {
		b.Fatalf("failed to parse header for testing file - %v", err)
	}

	b.Run("zero_copy_single_field", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			p := NewZeroCopyParser(reader, testHeaders, "LapCurrentLapTime")
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
			p := NewZeroCopyParser(reader, testHeaders, fields...)
			for {
				_, hasNext := p.NextZeroCopy()
				if !hasNext {
					break
				}
			}
		}
	})
}
