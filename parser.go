package ibt

import (
	"github.com/OJPARKINSON/ibt/headers"
)

// StructParser wraps DirectStructParser for backward compatibility
type StructParser struct {
	*DirectStructParser
}

// NewStructParser creates a parser that outputs TelemetryTick structs
func NewStructParser(reader *IbtReader, header *headers.Header, whitelist ...string) *StructParser {
	directParser := NewDirectStructParser(reader, header, whitelist...)

	return &StructParser{
		DirectStructParser: directParser,
	}
}

// NextStruct returns the next telemetry tick as a struct
// Delegates to DirectStructParser.NextStruct() for optimal performance
func (p *StructParser) NextStruct(tick *TelemetryTick) (*TelemetryTick, bool) {
	return p.DirectStructParser.NextStruct(tick)
}
