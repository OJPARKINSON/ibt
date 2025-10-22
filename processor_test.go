package ibt

import (
	"context"
	"errors"
	"testing"

	"github.com/OJPARKINSON/ibt/headers"
)

type testProcessor struct {
	results   []*TelemetryTick
	session   *headers.Session
	whitelist []string
}

func (t *testProcessor) Fields() interface{} {
	return struct {
		LapCurrentLapTime float64 `ibt:"LapCurrentLapTime"`
	}{}
}

func (t *testProcessor) ProcessStruct(tick *TelemetryTick, hasNext bool, session *headers.Session) error {
	t.results = append(t.results, tick)
	t.session = session
	return nil
}

func (t *testProcessor) FlushPendingData() error { return nil }

func (t *testProcessor) Close() error { return nil }

func (t *testProcessor) GetMetrics() interface{} { return nil }

type testErrorProcessor struct{}

func (t *testErrorProcessor) Fields() interface{} {
	return struct {
		LapCurrentLapTime float64 `ibt:"LapCurrentLapTime"`
	}{}
}

func (t *testErrorProcessor) ProcessStruct(tick *TelemetryTick, hasNext bool, session *headers.Session) error {
	return errors.New("unit test error")
}

func (t *testErrorProcessor) FlushPendingData() error { return nil }

func (t *testErrorProcessor) Close() error { return nil }

func (t *testErrorProcessor) GetMetrics() interface{} { return nil }

func TestProcess(t *testing.T) {
	reader, err := NewMmapReader(".testing/valid_test_file.ibt")
	if err != nil {
		t.Errorf("failed to open testing file - %v", err)
		return
	}
	defer reader.Close()

	testHeaders, err := headers.ParseHeaders(reader)
	if err != nil {
		t.Errorf("failed to parse header for testing file - %v", err)
		return
	}

	stubs := StubGroup{
		{filepath: ".testing/valid_test_file.ibt", header: testHeaders, r: reader},
	}

	t.Run("test Process() normal processor", func(t *testing.T) {
		proc := testProcessor{}

		if err := Process(context.Background(), stubs, &proc); err != nil {
			t.Errorf("expected Process() to run without err. received error: %v", err)
		}

		t.Logf("Total results received: %d", len(proc.results))
		if len(proc.results) > 0 {
			t.Logf("First result: %v", proc.results[0].LapCurrentLapTime)
		}
		if len(proc.results) > 69 {
			t.Logf("Result at index 69: %v", proc.results[69].LapCurrentLapTime)
		}

		// Test that we're getting the expected values from parser position 0
		if len(proc.results) == 0 {
			t.Errorf("expected results but got none")
			return
		}

		actualFirst := float32(proc.results[0].LapCurrentLapTime)
		expectedFirst := float32(37.6619)
		if actualFirst != expectedFirst {
			t.Errorf("expected value to check to be %f. got %f", expectedFirst, actualFirst)
		}

		// Test a later value to ensure progression
		if len(proc.results) > 69 {
			actualLater := float32(proc.results[69].LapCurrentLapTime)
			expectedLater := float32(38.8119)
			if actualLater != expectedLater {
				t.Errorf("expected later value to be %f. got %f", expectedLater, actualLater)
			}
		}
	})

	t.Run("test Process() err processor", func(t *testing.T) {
		proc := testErrorProcessor{}

		if err := Process(context.Background(), stubs, &proc); err == nil {
			t.Error("expected Process() to return an error")
		}
	})

	t.Run("test process() invalid file", func(t *testing.T) {
		proc := testProcessor{}

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		if err := process(ctx, stubs[0], &proc); err == nil {
			t.Errorf("expected process() to exit with a context done error")
		}
	})
}

func TestFieldsExtraction(t *testing.T) {
	reader, err := NewMmapReader(".testing/valid_test_file.ibt")
	if err != nil {
		t.Errorf("failed to open testing file - %v", err)
		return
	}
	defer reader.Close()

	testHeaders, err := headers.ParseHeaders(reader)
	if err != nil {
		t.Errorf("failed to parse header for testing file - %v", err)
		return
	}

	varHeader := testHeaders.VarHeader

	t.Run("test Fields() auto-extraction", func(t *testing.T) {
		proc := testProcessor{}
		fields := proc.Fields()
		whitelist := BuildWhitelistFromStruct(fields)

		if len(whitelist) != 1 {
			t.Errorf("expected 1 field in whitelist. found %d", len(whitelist))
		}

		if whitelist[0] != "LapCurrentLapTime" {
			t.Errorf("expected field to be LapCurrentLapTime, got %s", whitelist[0])
		}
	})

	t.Run("test buildWhitelist extracts valid fields", func(t *testing.T) {
		proc := testProcessor{}

		// Test that buildWhitelist can extract from processor
		cols := buildWhitelist(varHeader, &proc)

		if len(cols) != 1 {
			t.Errorf("expected 1 column. found %d", len(cols))
		}

		if cols[0] != "LapCurrentLapTime" {
			t.Errorf("expected LapCurrentLapTime, got %s", cols[0])
		}

		// Verify the extracted field is valid
		if _, exists := varHeader[cols[0]]; !exists {
			t.Errorf("extracted field %s does not exist in varHeader", cols[0])
		}
	})
}
