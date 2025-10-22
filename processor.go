package ibt

import (
	"context"
	"fmt"
	"sort"

	"github.com/OJPARKINSON/ibt/headers"
	"github.com/OJPARKINSON/ibt/utilities"
)

// Processor processes telemetry data tick by tick.
// Whitelists are automatically extracted from the Fields() struct tags.
//
// Example:
//
//	type MyProcessor struct {
//		cache []*ibt.TelemetryTick
//	}
//
//	func (p *MyProcessor) Fields() interface{} {
//		return struct {
//			Speed float64 `ibt:"Speed"`
//			RPM   float64 `ibt:"RPM"`
//		}{}
//	}
//
//	func (p *MyProcessor) ProcessStruct(tick *ibt.TelemetryTick, ...) error {
//		p.cache = append(p.cache, tick)
//		return nil
//	}
type Processor interface {
	// ProcessStruct is called for each tick of telemetry data
	ProcessStruct(tick *TelemetryTick, hasNext bool, session *headers.Session) error

	// Fields returns a struct with ibt tags defining required telemetry fields
	// The whitelist is automatically extracted from these tags
	Fields() interface{}

	// FlushPendingData flushes any cached data
	FlushPendingData() error

	// Close finalizes the processor
	Close() error

	// GetMetrics returns processor metrics for monitoring
	GetMetrics() interface{}
}

func Process(ctx context.Context, stubs StubGroup, processors ...Processor) error {
	sort.Sort(stubs)

	for _, stub := range stubs {
		// Open the stub's reader before processing
		if err := stub.Open(); err != nil {
			return fmt.Errorf("failed to open stub: %w", err)
		}

		if err := process(ctx, stub, processors...); err != nil {
			return err
		}
	}

	return nil
}

func process(ctx context.Context, stub Stub, processors ...Processor) error {
	header := stub.header

	// Auto-extract whitelist from all processors using their Fields() definitions
	whitelist := buildWhitelist(header.VarHeader, processors...)

	// Create struct parser (always uses struct-based parsing)
	parser := NewStructParser(stub.r, header, whitelist...)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		tick, hasNext := parser.NextStruct()
		if tick == nil {
			break
		}

		// Process all processors (single simple path - no legacy branching)
		for _, processor := range processors {
			if err := processor.ProcessStruct(tick, hasNext, header.SessionInfo); err != nil {
				return err
			}
		}

		if !hasNext {
			break
		}
	}

	return nil
}

// buildWhitelist compiles the whitelists from all processors by auto-extracting
// field names from their Fields() struct tags, and removes overlap.
func buildWhitelist(vars map[string]headers.VarHeader, processors ...Processor) []string {
	whitelist := make([]string, 0)

	for _, proc := range processors {
		// Auto-extract whitelist from Fields() struct tags
		fields := proc.Fields()
		autoWhitelist := BuildWhitelistFromStruct(fields)

		// Validate against available vars in the file
		for _, fieldName := range autoWhitelist {
			if _, ok := vars[fieldName]; ok {
				whitelist = append(whitelist, fieldName)
			}
		}
	}

	return utilities.GetDistinct(whitelist)
}
