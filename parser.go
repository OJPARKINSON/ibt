package ibt

import (
	"github.com/OJPARKINSON/ibt/headers"
)

// StructParser wraps DirectStructParser for backward compatibility
// This now uses direct byte-to-struct conversion instead of map intermediate
type StructParser struct {
	*DirectStructParser
}

// NewStructParser creates a parser that outputs TelemetryTick structs
// Now uses DirectStructParser for 10-15% performance improvement
func NewStructParser(reader *MmapReader, header *headers.Header, whitelist ...string) *StructParser {
	directParser := NewDirectStructParser(reader, header, whitelist...)

	return &StructParser{
		DirectStructParser: directParser,
	}
}

// NextStruct returns the next telemetry tick as a struct
// Delegates to DirectStructParser.NextStruct() for optimal performance
func (p *StructParser) NextStruct() (*TelemetryTick, bool) {
	return p.DirectStructParser.NextStruct()
}

// Parser is used to iterate and process telemetry variables for a given ibt file and it's headers.
//
// This parser now uses DirectStructParser internally for optimal performance (50-60% faster
// than the old map-based approach), then converts structs to maps for backward compatibility.
type Parser struct {
	// Embedded DirectStructParser for high-performance struct-based parsing
	*DirectStructParser

	// Whitelist is stored for backward compatibility with UpdateWhitelist()
	whitelist []string
}

// NewParser creates a new parser from a given ibt file, it's headers, and a variable whitelist.
//
// reader - Opened ibt file.
//
// header - Parsed headers of ibt file.
//
// whitelist - Variables to process. For example, "gear", "speed", "rpm" etc. If no values or a
// single value of "*" is received, all variables will be processed.
//
// This implementation now uses DirectStructParser internally for 50-60% performance improvement
// while maintaining the same public API.
func NewParser(reader *MmapReader, header *headers.Header, whitelist ...string) *Parser {
	// Create DirectStructParser for high-performance struct-based parsing
	directParser := NewDirectStructParser(reader, header, whitelist...)

	return &Parser{
		DirectStructParser: directParser,
		whitelist:         whitelist,
	}
}

// Next parses and returns the next tick of telemetry variables and whether it can be called again.
//
// A return of false will indicate that the buffer has reached the end. If the buffer has reached the end and Next() is called again,
// a nil and false will be returned. Additionally, a check can be done to check if the returned Tick is nil to determine if the EOF was reached.
//
// Should expected variable values be missing, please ensure that they are added to the Parser whitelist.
//
// Implementation: This method now uses DirectStructParser.NextStruct() internally for 50-60% better
// performance, then converts the struct to a map for backward compatibility.
func (p *Parser) Next() (Tick, bool) {
	// Use struct-based parsing (fast path)
	tick, hasNext := p.DirectStructParser.NextStruct()
	if tick == nil {
		return nil, false
	}

	// Convert struct to map for backward compatibility
	return tick.ToMap(p.whitelist), hasNext
}

// ParseAt the given buffer offset and return a processed tick.
//
// ParseAt is useful if a specific offset is known. An example would be the
// telemetry variable buffers that be provided during live telemetry parsing.
//
// When nil is returned, the buffer has reached EOF.
//
// Implementation: Uses DirectStructParser for parsing, then converts to map.
func (p *Parser) ParseAt(offset int) Tick {
	// Calculate which tick this offset corresponds to
	tickIndex := (offset - p.DirectStructParser.header.TelemetryHeader.BufOffset) / p.DirectStructParser.header.TelemetryHeader.BufLen

	// Save current position
	savedCurrent := p.DirectStructParser.current

	// Seek to the target tick
	p.Seek(tickIndex)

	// Parse using struct parser
	tick, _ := p.DirectStructParser.NextStruct()

	// Restore position
	p.DirectStructParser.current = savedCurrent

	if tick == nil {
		return nil
	}

	// Convert struct to map for backward compatibility
	return tick.ToMap(p.whitelist)
}

// Seek the parser to a specific tick within the ibt file.
// This method delegates to the embedded DirectStructParser.
func (p *Parser) Seek(iter int) {
	p.DirectStructParser.current = iter
}

// UpdateWhitelist replaces the current whitelist with the given fields.
//
// Note: This method updates the Parser's whitelist but does not rebuild the DirectStructParser.
// For best performance with a new whitelist, create a new Parser instance.
func (p *Parser) UpdateWhitelist(whitelist ...string) {
	p.whitelist = whitelist
	p.DirectStructParser.whitelist = whitelist
}
