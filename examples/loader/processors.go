package main

import (
	"fmt"

	"github.com/OJPARKINSON/ibt"
	"github.com/OJPARKINSON/ibt/headers"
)

func (l *loaderProcessor) Init(session *headers.Session) error { return nil }

type loaderProcessor struct {
	// Our storage client
	*storage
	// Cache for holding telemetry ticks
	cache []*ibt.TelemetryTick
	// Group number for identifying batches
	groupNumber int
	threshold   int
}

// Simple Constructor for creating our processor
func newLoaderProcessor(storage *storage, groupNumber int, threshold int) *loaderProcessor {
	return &loaderProcessor{
		storage:     storage,
		cache:       make([]*ibt.TelemetryTick, 0, threshold),
		groupNumber: groupNumber,
		threshold:   threshold,
	}
}

// Fields defines the telemetry fields this processor needs.
// Whitelist is automatically extracted from ibt tags.
func (l *loaderProcessor) Fields() interface{} {
	return struct {
		LapID       int32   `ibt:"Lap"`
		ThrottleRaw float64 `ibt:"ThrottleRaw"`
		BrakeRaw    float64 `ibt:"BrakeRaw"`
		Clutch      float64 `ibt:"Clutch"`
		LapDistPct  float64 `ibt:"LapDistPct"`
		Lat         float64 `ibt:"Lat"`
		Lon         float64 `ibt:"Lon"`
	}{}
}

// ProcessStruct processes a single tick of telemetry
func (l *loaderProcessor) ProcessStruct(tick *ibt.TelemetryTick, hasNext bool) error {
	// Set group number
	tick.GroupNum = l.groupNumber

	// Add to cache
	l.cache = append(l.cache, tick)

	// If cache is past threshold, bulk load
	if len(l.cache) >= l.threshold {
		if err := l.loadBatch(); err != nil {
			return fmt.Errorf("failed to load batch - %v", err)
		}
		// Reset cache
		l.cache = l.cache[:0]
	}

	return nil
}

// FlushPendingData flushes any remaining cached data
func (l *loaderProcessor) FlushPendingData() error {
	if len(l.cache) > 0 {
		return l.loadBatch()
	}
	return nil
}

// Close finalizes the processor
func (l *loaderProcessor) Close() error {
	return l.FlushPendingData()
}

// GetMetrics returns processor metrics
func (l *loaderProcessor) GetMetrics() interface{} {
	return map[string]interface{}{
		"batches_loaded": l.storage.Loaded(),
		"cache_size":     len(l.cache),
	}
}

func (l *loaderProcessor) loadBatch() error {
	// Bulk load our batch to storage.
	return l.ExecStructs(l.cache)
}
