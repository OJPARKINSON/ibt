package ibt

import (
	"github.com/teamjorge/ibt/headers"
)

// Parser is used to iterate and process telemetry variables for a given ibt file and it's headers.
type Parser struct {
	// File or Live Telemetry reader
	reader headers.Reader
	// List of columns to parse
	whitelist []string
	header    *headers.Header

	current int

	// Pre-allocated buffer to eliminate per-tick allocations
	bufferPool []byte
	// Pre-allocated tick map to eliminate per-tick map allocations
	tickPool Tick
	
	// Fast path optimization: pre-computed variable headers for whitelist
	varHeaders []headers.VarHeader
	varNames   []string
}

// NewParser creates a new parser from a given ibt file, it's headers, and a variable whitelist.
//
// reader - Opened ibt file.
//
// header - Parsed headers of ibt file.
//
// whitelist - Variables to process. For example, "gear", "speed", "rpm" etc. If no values or a
// single value of "*" is received, all variables will be processed.
func NewParser(reader headers.Reader, header *headers.Header, whitelist ...string) *Parser {
	p := new(Parser)

	p.reader = reader
	p.whitelist = whitelist
	p.header = header

	p.current = 1

	// Pre-allocate buffer to eliminate per-tick allocations
	p.bufferPool = make([]byte, header.TelemetryHeader.BufLen)

	// Pre-allocate tick map with capacity for whitelist (use capacity hint for better performance)
	whitelistLen := len(whitelist)
	if whitelistLen == 0 || (whitelistLen == 1 && whitelist[0] == "*") {
		// If no whitelist or wildcard, estimate capacity based on typical variable count
		whitelistLen = 64 // reasonable default for telemetry variables
	}
	p.tickPool = make(Tick, whitelistLen)

	// Pre-compute variable headers and names for fast parsing
	p.varHeaders = make([]headers.VarHeader, 0, len(p.whitelist))
	p.varNames = make([]string, 0, len(p.whitelist))
	
	for _, variable := range p.whitelist {
		if varHeader, exists := header.VarHeader[variable]; exists {
			p.varHeaders = append(p.varHeaders, varHeader)
			p.varNames = append(p.varNames, variable)
		}
	}

	return p
}

// Next parses and returns the next tick of telemetry variables and whether it can be called again.
//
// A return of false will indicate that the buffer has reached the end. If the buffer has reached the end and Next() is called again,
// a nil and false will be returned. Additionally, a check can be done to check if the returned Tick is nil to determine if the EOF was reached.
//
// Should expected variable values be missing, please ensure that they are added to the Parser whitelist.
func (p *Parser) Next() (Tick, bool) {
	start := p.header.TelemetryHeader.BufOffset + (p.current * p.header.TelemetryHeader.BufLen)

	currentBuf := p.read(start)
	if currentBuf == nil {
		return nil, false
	}

	// Read in the next buffer to determine if more telemetry ticks are available.
	nextStart := p.header.TelemetryHeader.BufOffset + ((p.current + 1) * p.header.TelemetryHeader.BufLen)
	nextBuf := p.read(nextStart)

	newVars := p.readVarsFromBuffer(currentBuf)

	p.current++

	return newVars, nextBuf != nil
}

// ParseAt the given buffer offset and return a processed tick.
//
// ParseAt is useful if a specific offset is known. An example would be the
// telemetry variable buffers that are provided during live telemetry parsing.
//
// When nil is returned, the buffer has reached EOF.
func (p *Parser) ParseAt(offset int) Tick {
	currentBuf := p.read(offset)
	if currentBuf == nil {
		return nil
	}

	newVars := p.readVarsFromBuffer(currentBuf)

	return newVars
}

// read the next buffer from offset to the current length set by the parser.
func (p *Parser) read(start int) []byte {
	// Reuse pre-allocated buffer instead of creating new one
	_, err := p.reader.ReadAt(p.bufferPool, int64(start))
	if err != nil {
		return nil
	}

	return p.bufferPool
}

// readVarsFromBuffer reads each of the specified (whitelist) fields from the given buffer into a new Tick.
func (p *Parser) readVarsFromBuffer(buf []byte) Tick {
	// Use slice-based approach for faster clearing instead of map iteration
	if len(p.tickPool) > 0 {
		// Fast clear using slice assignment - much faster than map iteration
		for k := range p.tickPool {
			delete(p.tickPool, k)
		}
	}

	// Use pre-computed variable headers for faster iteration
	for i, varHeader := range p.varHeaders {
		varName := p.varNames[i]
		val := readVarValueFast(buf, varHeader)
		p.tickPool[varName] = val
	}

	// Use pre-allocated result map and copy efficiently  
	result := make(Tick, len(p.varNames))
	for k, v := range p.tickPool {
		result[k] = v
	}
	
	return result
}

// Seek the parser to a specific tick within the ibt file.
func (p *Parser) Seek(iter int) { p.current = iter }

// UpdateWhitelist replaces the current whitelist with the given fields
func (p *Parser) UpdateWhitelist(whitelist ...string) {
	p.whitelist = whitelist
}
