package ibt

import (
	"reflect"
	"testing"
)

// Test struct templates.
type BasicTelemetry struct {
	Speed float64 `ibt:"Speed"`
	Gear  uint32  `ibt:"Gear"`
	RPM   float64 `ibt:"RPM"`
}

type ExtendedTelemetry struct {
	Speed    float64 `ibt:"Speed"`
	Gear     uint32  `ibt:"Gear"`
	RPM      float64 `ibt:"RPM"`
	Throttle float64 `ibt:"Throttle"`
	Brake    float64 `ibt:"Brake"`
	Lap      int32   `ibt:"Lap"`
}

type EmptyTemplate struct {
	// No ibt tags
	Speed float64
	Gear  uint32
}

type MixedTemplate struct {
	Speed    float64 `ibt:"Speed"` // Has tag
	Gear     uint32  // No tag
	Throttle float64 `ibt:"Throttle"` // Has tag
}

func TestBuildWhitelistFromStruct(t *testing.T) {
	tests := []struct {
		name     string
		template any
		expected []string
	}{
		{
			name:     "basic struct with 3 fields",
			template: BasicTelemetry{},
			expected: []string{"Speed", "Gear", "RPM"},
		},
		{
			name:     "extended struct with 6 fields",
			template: ExtendedTelemetry{},
			expected: []string{"Speed", "Gear", "RPM", "Throttle", "Brake", "Lap"},
		},
		{
			name:     "empty struct with no tags",
			template: EmptyTemplate{},
			expected: []string{},
		},
		{
			name:     "mixed struct (some fields with tags)",
			template: MixedTemplate{},
			expected: []string{"Speed", "Throttle"},
		},
		{
			name:     "pointer to struct",
			template: &BasicTelemetry{},
			expected: []string{"Speed", "Gear", "RPM"},
		},
		{
			name:     "reflect.Type",
			template: reflect.TypeFor[BasicTelemetry](),
			expected: []string{"Speed", "Gear", "RPM"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := BuildWhitelistFromStruct(tt.template)

			if len(result) != len(tt.expected) {
				t.Errorf("BuildWhitelistFromStruct() returned %d fields, expected %d", len(result), len(tt.expected))
				t.Errorf("Got: %v", result)
				t.Errorf("Expected: %v", tt.expected)
				return
			}

			for i, expected := range tt.expected {
				if result[i] != expected {
					t.Errorf("Field %d: got %q, expected %q", i, result[i], expected)
				}
			}
		})
	}
}

func TestBuildWhitelistFromStruct_EdgeCases(t *testing.T) {
	t.Run("non-struct type returns empty", func(t *testing.T) {
		result := BuildWhitelistFromStruct(42)
		if len(result) != 0 {
			t.Errorf("Expected empty whitelist for non-struct, got %v", result)
		}
	})

	t.Run("string type returns empty", func(t *testing.T) {
		result := BuildWhitelistFromStruct("hello")
		if len(result) != 0 {
			t.Errorf("Expected empty whitelist for string, got %v", result)
		}
	})
}

func BenchmarkBuildWhitelistFromStruct(b *testing.B) {
	template := BasicTelemetry{}

	for b.Loop() {
		_ = BuildWhitelistFromStruct(template)
	}
}

func BenchmarkBuildWhitelistFromStruct_Extended(b *testing.B) {
	template := ExtendedTelemetry{}

	for b.Loop() {
		_ = BuildWhitelistFromStruct(template)
	}
}
