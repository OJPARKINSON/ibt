package ibt

import (
	"context"
	"errors"
	"fmt"
	"sort"

	"github.com/teamjorge/ibt/headers"
	"github.com/teamjorge/ibt/utilities"
)

type Processor interface {
	Process(input Tick, hasNext bool, session *headers.Session) error
	Whitelist() []string
	FlushPendingData() error
	Close() error
	GetMetrics() interface{} // Returns ProcessorMetrics from internal/processing
}

type StructProcessor interface {
	Processor
	ProcessStruct(tick *TelemetryTick, hasNext bool, session *headers.Session) error
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

	// Only parse fields that are actually needed by all processors combined
	whitelist := buildWhitelist(header.VarHeader, processors...)

	useStructs := false
	for _, processor := range processors {
		if _, ok := processor.(StructProcessor); ok {
			useStructs = true
			break
		}
	}

	if useStructs {
		whitelist := processors[0].Whitelist()
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

			for _, processor := range processors {
				if sp, ok := processor.(StructProcessor); ok {
					if err := sp.ProcessStruct(tick, hasNext, header.SessionInfo); err != nil {
						return err
					}
				}
			}

			if !hasNext {
				break
			}
		}
	} else {

		// Use optimized parser with all our performance improvements
		parser := NewZeroCopyParser(stub.r, header, whitelist...)
		for {
			select {
			case <-ctx.Done():
				return errors.New("context cancelled")
			default:
			}

			tick, hasNext := parser.NextZeroCopy()

			// Process all processors with the same tick - avoid redundant filtering
			for _, proc := range processors {
				procWhitelist := proc.Whitelist()

				// If processor needs all fields, use original tick
				if len(procWhitelist) >= len(whitelist) {
					if err := proc.Process(tick, hasNext, header.SessionInfo); err != nil {
						return err
					}
				} else {
					// Filter tick for this specific processor
					filteredTick := tick.Filter(procWhitelist...)
					if err := proc.Process(filteredTick, hasNext, header.SessionInfo); err != nil {
						return err
					}
				}
			}

			if !hasNext {
				break
			}
		}
	}

	return nil
}

// getcinoketeWhitelist compiles the whitelists from all processors and removes overlap
func buildWhitelist(vars map[string]headers.VarHeader, processors ...Processor) []string {
	whitelist := make([]string, 0)

	for _, proc := range processors {
		whitelist = append(whitelist, parseAndValidateWhitelist(vars, proc)...)
	}

	return utilities.GetDistinct(whitelist)
}

// parseWhitelist will retrieve vars when * is used and ensure a unique list
//
// Variables that are not found in the VarHeader will automatically be excluded.
func parseAndValidateWhitelist(vars map[string]headers.VarHeader, processor Processor) []string {
	whitelist := processor.Whitelist()

	if len(whitelist) == 0 {
		return headers.AvailableVars(vars)
	}

	for _, col := range whitelist {
		if col == "*" {
			return headers.AvailableVars(vars)
		}
	}

	columns := make([]string, 0)

	// Ensure only valid columns are added
	for _, col := range whitelist {
		if _, ok := vars[col]; ok {
			columns = append(columns, col)
		}
	}

	return columns
}
