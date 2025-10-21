package ibt

import (
	"reflect"
	"testing"
)

// Test struct templates
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
	Speed    float64 `ibt:"Speed"`    // Has tag
	Gear     uint32  // No tag
	Throttle float64 `ibt:"Throttle"` // Has tag
}

func TestBuildWhitelistFromStruct(t *testing.T) {
	tests := []struct {
		name     string
		template interface{}
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
			template: reflect.TypeOf(BasicTelemetry{}),
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

			// Check that all expected fields are present
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

func TestToMapFromStruct(t *testing.T) {
	// Create a sample TelemetryTick with known values
	tick := &TelemetryTick{
		Speed:    float64(123.45),
		Gear:     uint32(4),
		RPM:      float64(5500.0),
		Throttle: float64(0.85),
		Brake:    float64(0.0),
		LapID:    int32(5),
	}

	t.Run("basic template extraction", func(t *testing.T) {
		result := tick.ToMapFromStruct(BasicTelemetry{})

		// Check map size
		if len(result) != 3 {
			t.Errorf("Expected 3 fields, got %d: %v", len(result), result)
		}

		// Check values
		if speed, ok := result["Speed"].(float32); !ok || speed != float32(123.45) {
			t.Errorf("Speed: got %v, expected 123.45", result["Speed"])
		}

		if gear, ok := result["Gear"].(uint32); !ok || gear != 4 {
			t.Errorf("Gear: got %v, expected 4", result["Gear"])
		}

		if rpm, ok := result["RPM"].(float32); !ok || rpm != float32(5500.0) {
			t.Errorf("RPM: got %v, expected 5500.0", result["RPM"])
		}
	})

	t.Run("extended template extraction", func(t *testing.T) {
		result := tick.ToMapFromStruct(ExtendedTelemetry{})

		// Check map size
		if len(result) != 6 {
			t.Errorf("Expected 6 fields, got %d: %v", len(result), result)
		}

		// Check that all fields are present
		expectedFields := []string{"Speed", "Gear", "RPM", "Throttle", "Brake", "Lap"}
		for _, field := range expectedFields {
			if _, exists := result[field]; !exists {
				t.Errorf("Missing field: %s", field)
			}
		}
	})

	t.Run("empty template returns empty map", func(t *testing.T) {
		result := tick.ToMapFromStruct(EmptyTemplate{})

		if len(result) != 0 {
			t.Errorf("Expected empty map, got %d fields: %v", len(result), result)
		}
	})

	t.Run("mixed template only extracts tagged fields", func(t *testing.T) {
		result := tick.ToMapFromStruct(MixedTemplate{})

		if len(result) != 2 {
			t.Errorf("Expected 2 fields, got %d: %v", len(result), result)
		}

		if _, exists := result["Speed"]; !exists {
			t.Error("Expected Speed field")
		}

		if _, exists := result["Throttle"]; !exists {
			t.Error("Expected Throttle field")
		}

		if _, exists := result["Gear"]; exists {
			t.Error("Gear should not be present (no tag)")
		}
	})
}

func TestToMapWithTemplate(t *testing.T) {
	// This test requires a parser, so we'll create a minimal mock setup
	// For now, test that the method exists and can be called

	// Create a sample TelemetryTick
	tick := &TelemetryTick{
		Speed:    float64(123.45),
		Gear:     uint32(4),
		RPM:      float64(5500.0),
		Throttle: float64(0.85),
	}

	// Create a parser (we'll use nil reader/header for this test)
	parser := &DirectStructParser{
		templateCache: make(map[reflect.Type][]string),
	}

	t.Run("first call caches template", func(t *testing.T) {
		template := BasicTelemetry{}
		result := parser.ToMapWithTemplate(tick, template)

		// Check result
		if len(result) != 3 {
			t.Errorf("Expected 3 fields, got %d: %v", len(result), result)
		}

		// Check cache
		typ := reflect.TypeOf(template)
		parser.cacheMu.RLock()
		_, cached := parser.templateCache[typ]
		parser.cacheMu.RUnlock()

		if !cached {
			t.Error("Template should be cached after first call")
		}
	})

	t.Run("second call uses cache", func(t *testing.T) {
		template := BasicTelemetry{}

		// First call
		_ = parser.ToMapWithTemplate(tick, template)

		// Second call should use cache
		result := parser.ToMapWithTemplate(tick, template)

		if len(result) != 3 {
			t.Errorf("Expected 3 fields, got %d: %v", len(result), result)
		}

		// Verify cache still has the entry
		typ := reflect.TypeOf(template)
		parser.cacheMu.RLock()
		whitelist, cached := parser.templateCache[typ]
		parser.cacheMu.RUnlock()

		if !cached {
			t.Error("Template should still be cached")
		}

		if len(whitelist) != 3 {
			t.Errorf("Cached whitelist should have 3 entries, got %d", len(whitelist))
		}
	})

	t.Run("different templates get separate cache entries", func(t *testing.T) {
		parser := &DirectStructParser{
			templateCache: make(map[reflect.Type][]string),
		}

		template1 := BasicTelemetry{}
		template2 := ExtendedTelemetry{}

		_ = parser.ToMapWithTemplate(tick, template1)
		_ = parser.ToMapWithTemplate(tick, template2)

		parser.cacheMu.RLock()
		cacheSize := len(parser.templateCache)
		parser.cacheMu.RUnlock()

		if cacheSize != 2 {
			t.Errorf("Expected 2 cache entries, got %d", cacheSize)
		}
	})

	t.Run("ClearTemplateCache removes all entries", func(t *testing.T) {
		parser := &DirectStructParser{
			templateCache: make(map[reflect.Type][]string),
		}

		// Add some entries
		_ = parser.ToMapWithTemplate(tick, BasicTelemetry{})
		_ = parser.ToMapWithTemplate(tick, ExtendedTelemetry{})

		// Clear cache
		parser.ClearTemplateCache()

		parser.cacheMu.RLock()
		cacheSize := len(parser.templateCache)
		parser.cacheMu.RUnlock()

		if cacheSize != 0 {
			t.Errorf("Expected cache to be empty, got %d entries", cacheSize)
		}
	})

	t.Run("pointer to template works correctly", func(t *testing.T) {
		parser := &DirectStructParser{
			templateCache: make(map[reflect.Type][]string),
		}

		template := &BasicTelemetry{}
		result := parser.ToMapWithTemplate(tick, template)

		if len(result) != 3 {
			t.Errorf("Expected 3 fields, got %d: %v", len(result), result)
		}
	})
}

// Benchmark tests
func BenchmarkToMap_StringWhitelist(b *testing.B) {
	tick := &TelemetryTick{
		Speed:    float64(123.45),
		Gear:     uint32(4),
		RPM:      float64(5500.0),
		Throttle: float64(0.85),
		Brake:    float64(0.0),
	}

	whitelist := []string{"Speed", "Gear", "RPM"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = tick.ToMap(whitelist)
	}
}

func BenchmarkToMapFromStruct(b *testing.B) {
	tick := &TelemetryTick{
		Speed:    float64(123.45),
		Gear:     uint32(4),
		RPM:      float64(5500.0),
		Throttle: float64(0.85),
		Brake:    float64(0.0),
	}

	template := BasicTelemetry{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = tick.ToMapFromStruct(template)
	}
}

func BenchmarkToMapWithTemplate_FirstCall(b *testing.B) {
	tick := &TelemetryTick{
		Speed:    float64(123.45),
		Gear:     uint32(4),
		RPM:      float64(5500.0),
		Throttle: float64(0.85),
		Brake:    float64(0.0),
	}

	template := BasicTelemetry{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		parser := &DirectStructParser{
			templateCache: make(map[reflect.Type][]string),
		}
		_ = parser.ToMapWithTemplate(tick, template)
	}
}

func BenchmarkToMapWithTemplate_CachedCall(b *testing.B) {
	tick := &TelemetryTick{
		Speed:    float64(123.45),
		Gear:     uint32(4),
		RPM:      float64(5500.0),
		Throttle: float64(0.85),
		Brake:    float64(0.0),
	}

	template := BasicTelemetry{}

	parser := &DirectStructParser{
		templateCache: make(map[reflect.Type][]string),
	}

	// Prime the cache
	_ = parser.ToMapWithTemplate(tick, template)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = parser.ToMapWithTemplate(tick, template)
	}
}

func BenchmarkBuildWhitelistFromStruct(b *testing.B) {
	template := BasicTelemetry{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = BuildWhitelistFromStruct(template)
	}
}

func BenchmarkBuildWhitelistFromStruct_Extended(b *testing.B) {
	template := ExtendedTelemetry{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = BuildWhitelistFromStruct(template)
	}
}
