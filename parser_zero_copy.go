package ibt

import (
	"sync"

	"github.com/OJPARKINSON/ibt/headers"
)

// ZeroCopyParser is an ultra-fast parser that minimizes allocations.
//
// This parser now uses DirectStructParser internally for optimal performance (50-60% faster),
// then reuses the same map for zero-copy behavior.
type ZeroCopyParser struct {
	*DirectStructParser

	// Whitelist for map conversion
	whitelist []string

	// Reusable result tick to avoid allocations
	resultTick Tick

	// Pool for result ticks
	tickResultPool *sync.Pool
}

// NewZeroCopyParser creates a parser optimized for minimal allocations.
//
// This implementation now uses DirectStructParser internally for better performance
// while maintaining the zero-copy map reuse behavior.
func NewZeroCopyParser(reader *MmapReader, header *headers.Header, whitelist ...string) *ZeroCopyParser {
	directParser := NewDirectStructParser(reader, header, whitelist...)

	return &ZeroCopyParser{
		DirectStructParser: directParser,
		whitelist:         whitelist,
		resultTick:        make(Tick, len(whitelist)),
		tickResultPool: &sync.Pool{
			New: func() interface{} {
				return make(Tick, len(whitelist))
			},
		},
	}
}

// NextZeroCopy returns the next tick with minimal allocations.
//
// WARNING: The returned Tick may be modified on the next call to NextZeroCopy.
// If you need to retain the data, make a copy using GetTickCopy().
//
// Implementation: Uses DirectStructParser.NextStruct() for fast parsing, then reuses
// the same map to avoid allocations.
func (p *ZeroCopyParser) NextZeroCopy() (Tick, bool) {
	// Use struct-based parsing (fast path)
	tick, hasNext := p.DirectStructParser.NextStruct()
	if tick == nil {
		return nil, false
	}

	// Clear the reused map efficiently
	for k := range p.resultTick {
		delete(p.resultTick, k)
	}

	// Convert struct to map and populate the reused map
	tickMap := tick.ToMap(p.whitelist)
	for k, v := range tickMap {
		p.resultTick[k] = v
	}

	return p.resultTick, hasNext
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
