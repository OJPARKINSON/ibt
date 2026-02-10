package main

import (
	"log"
	"maps"
	"slices"

	"github.com/OJPARKINSON/ibt"
	"github.com/OJPARKINSON/ibt/headers"
)

// trackTempProcessor tracks the track temperature for each lap.
type trackTempProcessor struct {
	tempMap map[int32]float64
}

func newTrackTempProcessor() *trackTempProcessor {
	return &trackTempProcessor{
		tempMap: make(map[int32]float64),
	}
}

func (t *trackTempProcessor) Init(session *headers.Session) error { return nil }

func (t *trackTempProcessor) Fields() any {
	return struct {
		LapID         int32   `ibt:"Lap"`
		TrackTempCrew float64 `ibt:"TrackTempCrew"`
	}{}
}

func (t *trackTempProcessor) ProcessStruct(tick *ibt.TelemetryTick, hasNext bool) error {
	t.tempMap[tick.LapID] = tick.TrackTempCrew
	return nil
}

func (t *trackTempProcessor) FlushPendingData() error { return nil }
func (t *trackTempProcessor) Close() error            { return nil }
func (t *trackTempProcessor) GetMetrics() any         { return nil }

func (t *trackTempProcessor) Print() {
	log.Println("Track Temp:")
	laps := slices.Sorted(maps.Keys(t.tempMap))
	for _, lap := range laps {
		log.Printf("%03d - %.3f°C\n", lap, t.tempMap[lap])
	}
}
