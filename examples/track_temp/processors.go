package main

import (
	"fmt"
	"sort"

	"github.com/OJPARKINSON/ibt"
	"github.com/OJPARKINSON/ibt/headers"
	"golang.org/x/exp/maps"
)

// TrackTempProcessor tracks the track temperature for each lap using struct-based processing
type trackTempProcessor struct {
	tempMap map[int32]float64
}

// NewTrackTempProcessor creates and initialises a new trackTempProcessor
func newTrackTempProcessor() *trackTempProcessor {
	t := new(trackTempProcessor)

	// tempMap will store a temperature value against a lap number
	t.tempMap = make(map[int32]float64)

	return t
}

// Display name of the processor
func (t *trackTempProcessor) Name() string { return "Track Temp" }

// Fields defines the telemetry fields this processor needs.
// Whitelist is automatically extracted from ibt tags.
func (t *trackTempProcessor) Fields() interface{} {
	return struct {
		LapID         int32   `ibt:"Lap"`
		TrackTempCrew float64 `ibt:"TrackTempCrew"`
	}{}
}

// ProcessStruct processes a single tick of telemetry
func (t *trackTempProcessor) ProcessStruct(tick *ibt.TelemetryTick, hasNext bool, session *headers.Session) error {
	// Store track temperature for this lap
	t.tempMap[tick.LapID] = tick.TrackTempCrew
	return nil
}

// FlushPendingData flushes any cached data
func (t *trackTempProcessor) FlushPendingData() error {
	return nil
}

// Close finalizes the processor
func (t *trackTempProcessor) Close() error {
	return nil
}

// GetMetrics returns processor metrics
func (t *trackTempProcessor) GetMetrics() interface{} {
	return map[string]interface{}{
		"total_laps": len(t.tempMap),
	}
}

// Print the summarised Track Temperature
func (t *trackTempProcessor) Print() {
	fmt.Println("Track Temp:")
	laps := maps.Keys(t.tempMap)
	sort.Slice(laps, func(i, j int) bool { return laps[i] < laps[j] })

	for _, lap := range laps {
		fmt.Printf("%03d - %.3f°C\n", lap, t.tempMap[lap])
	}
}
