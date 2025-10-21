// Package ibt provides high-performance parsing of iRacing telemetry (.ibt) files.
//
// The package supports multiple parser implementations optimized for different use cases,
// with performance improvements ranging from 17% to 57% over baseline implementations.
//
// # Parser Options
//
// The package provides three parser implementations:
//
//   - DirectStructParser: Zero-allocation direct byte parsing (~57% faster than baseline)
//   - StructParser: High-level wrapper around DirectStructParser (recommended for most use cases)
//   - Parser: Base parser providing map-based access (baseline performance)
//
// # Performance Characteristics
//
// Benchmarks on 1.85M records across 11 test files (AMD Ryzen 7 7600X):
//
//   - DirectStructParser: ~4.7s (397K records/sec) - 57% improvement
//   - StructParser: ~4.7s (397K records/sec) - 57% improvement
//   - Map-based baseline: ~10.8s (171K records/sec)
//
// The DirectStructParser achieves this performance through:
//
//   - Zero-copy memory-mapped I/O (eliminates buffer copying)
//   - Direct byte-to-struct field assignment (skips map intermediate)
//   - Safe byte conversions via encoding/binary (compiler-optimized)
//
// # Basic Usage
//
// Using the recommended StructParser:
//
//	reader, err := ibt.NewMmapReader("telemetry.ibt")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer reader.Close()
//
//	header, err := headers.ParseHeader(reader.Data())
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// Create parser with optional field whitelist for filtering
//	parser := ibt.NewStructParser(reader, header, "Speed", "RPM", "Gear")
//
//	for {
//	    tick, hasNext := parser.NextStruct()
//	    if tick == nil {
//	        break
//	    }
//	    fmt.Printf("Speed: %.2f, RPM: %.2f, Gear: %d\n", tick.Speed, tick.RPM, tick.Gear)
//	    if !hasNext {
//	        break
//	    }
//	}
//
// # Advanced Usage
//
// Using DirectStructParser directly for maximum performance:
//
//	parser := ibt.NewDirectStructParser(reader, header, "Speed", "RPM")
//	for {
//	    tick, hasNext := parser.NextStruct()
//	    if tick == nil {
//	        break
//	    }
//	    // Process tick...
//	    if !hasNext {
//	        break
//	    }
//	}
//
// # Design Decisions
//
// Safe Byte Conversions:
//
// This package uses encoding/binary.LittleEndian and math.Float32frombits for type
// conversions instead of unsafe pointer casts. Modern Go compilers (1.19+) optimize
// these safe operations to the same assembly as unsafe alternatives, providing both
// safety and performance.
//
// Zero-Copy Memory Mapping:
//
// The parser uses memory-mapped I/O to read .ibt files, returning slices directly
// from the mmap region without intermediate buffer allocations. This eliminates
// ~8.69s of runtime.memmove overhead (98.7% reduction).
//
// Field Mapping:
//
// DirectStructParser uses a pre-built map of field setter functions that directly
// assign bytes to struct fields, eliminating the map[string]interface{} intermediate
// layer and reducing allocations by ~715MB per run.
//
// # Thread Safety
//
// Parser instances are not safe for concurrent use. Each goroutine should have its
// own parser instance. The underlying MmapReader can be shared across goroutines
// as long as each has its own Parser.
//
// # Telemetry Fields
//
// The TelemetryTick struct contains 44 telemetry fields including:
//
//   - Lap data: LapID, LapDistPct, LapCurrentLapTime, LapLastLapTime, LapDeltaToBestLap
//   - Motion: Speed, VelocityX/Y/Z, LatAccel, LongAccel, VertAccel
//   - Controls: Throttle, Brake, Gear, SteeringWheelAngle, RPM
//   - Position: Lat, Lon, Alt, Pitch, Roll, Yaw, YawNorth
//   - Car state: FuelLevel, WaterTemp, Voltage
//   - Tires: LF/RF/LR/RRpressure, LF/RF/LR/RRtempM
//   - Session: SessionTime, SessionNum, PlayerCarIdx, PlayerCarPosition
//
// See TelemetryTick struct for complete field list with iRacing field name mappings.
package ibt
