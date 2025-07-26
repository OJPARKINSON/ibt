package ibt

import (
	"sync"
	"github.com/teamjorge/ibt/headers"
)

// ZeroCopyParser is an ultra-fast parser that minimizes allocations
type ZeroCopyParser struct {
	*Parser
	
	// Reusable result tick to avoid allocations
	resultTick Tick
	
	// Pool for result ticks
	tickResultPool *sync.Pool
}

// NewZeroCopyParser creates a parser optimized for minimal allocations
func NewZeroCopyParser(reader headers.Reader, header *headers.Header, whitelist ...string) *ZeroCopyParser {
	baseParser := NewParser(reader, header, whitelist...)
	
	return &ZeroCopyParser{
		Parser: baseParser,
		resultTick: make(Tick, len(whitelist)),
		tickResultPool: &sync.Pool{
			New: func() interface{} {
				return make(Tick, len(whitelist))
			},
		},
	}
}

// NextZeroCopy returns the next tick with minimal allocations
// WARNING: The returned Tick may be modified on the next call to NextZeroCopy
// If you need to retain the data, make a copy
func (p *ZeroCopyParser) NextZeroCopy() (Tick, bool) {
	start := p.header.TelemetryHeader.BufOffset + (p.current * p.header.TelemetryHeader.BufLen)

	currentBuf := p.read(start)
	if currentBuf == nil {
		return nil, false
	}

	// Read in the next buffer to determine if more telemetry ticks are available.
	nextStart := p.header.TelemetryHeader.BufOffset + ((p.current + 1) * p.header.TelemetryHeader.BufLen)
	nextBuf := p.read(nextStart)

	// Reuse the same tick map to avoid allocations
	p.readVarsFromBufferZeroCopy(currentBuf)

	p.current++

	return p.resultTick, nextBuf != nil
}

// readVarsFromBufferZeroCopy reads variables into the reused tick map
func (p *ZeroCopyParser) readVarsFromBufferZeroCopy(buf []byte) {
	// Clear the reused map efficiently
	for k := range p.resultTick {
		delete(p.resultTick, k)
	}

	// Use pre-computed variable headers for faster iteration
	for i, varHeader := range p.varHeaders {
		varName := p.varNames[i]
		val := readVarValueFast(buf, varHeader)
		p.resultTick[varName] = val
	}
}

// GetTickCopy returns a copy of the current tick that is safe to retain
func (p *ZeroCopyParser) GetTickCopy(tick Tick) Tick {
	result := p.tickResultPool.Get().(Tick)
	
	// Clear the pooled tick
	for k := range result {
		delete(result, k)
	}
	
	// Copy values
	for k, v := range tick {
		result[k] = v
	}
	
	return result
}

// ReturnTickCopy returns a tick copy back to the pool
func (p *ZeroCopyParser) ReturnTickCopy(tick Tick) {
	if tick != nil && len(tick) > 0 {
		p.tickResultPool.Put(tick)
	}
}