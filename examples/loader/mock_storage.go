package main

import "github.com/OJPARKINSON/ibt"

// This is a mock external storage client.
//
// Think of it as a database, API, or external file.
type storage struct {
	batchesLoaded int
}

// Simple constructor.
func newStorage() *storage { return new(storage) }

func (s *storage) Connect() error { return nil }

// ExecStructs processes struct-based telemetry data.
func (s *storage) ExecStructs(data []*ibt.TelemetryTick) error {
	s.batchesLoaded += len(data)
	return nil
}

func (s *storage) Close() error { return nil }

func (s *storage) Loaded() int { return s.batchesLoaded }
