// Package ibt provides high-performance parsing of iRacing telemetry (.ibt) files.
//
// The package uses direct byte-to-struct parsing for optimal performance,
// with safe byte conversions via encoding/binary that modern Go compilers
// optimize to the same assembly as unsafe pointer casts.
//
// # Basic Usage
//
// Using the Processor interface for tick-by-tick processing:
//
//	type MyProcessor struct {
//	    data []*ibt.TelemetryTick
//	}
//
//	func (p *MyProcessor) Fields() interface{} {
//	    return struct {
//	        Speed float64 `ibt:"Speed"`
//	        RPM   float64 `ibt:"RPM"`
//	    }{}
//	}
//
//	func (p *MyProcessor) Init(session *headers.Session) error { return nil }
//
//	func (p *MyProcessor) ProcessStruct(tick *ibt.TelemetryTick, hasNext bool) error {
//	    p.data = append(p.data, tick)
//	    return nil
//	}
//
// Then use ibt.Process() with stubs:
//
//	stubs, _ := ibt.ParseStubs("telemetry.ibt")
//	groups := stubs.Group()
//	for _, group := range groups {
//	    ibt.Process(ctx, group, processor)
//	}
//
// # Direct Parser Usage
//
//	reader, _ := ibt.NewIbtReader("telemetry.ibt")
//	defer reader.Close()
//
//	header, _ := headers.ParseHeaders(reader)
//	parser := ibt.NewStructParser(reader, header, "Speed", "RPM", "Gear")
//
//	tick := &ibt.TelemetryTick{}
//	for {
//	    tick, hasNext := parser.NextStruct(tick)
//	    if tick == nil {
//	        break
//	    }
//	    fmt.Printf("Speed: %.2f, RPM: %.2f\n", tick.Speed, tick.RPM)
//	    if !hasNext {
//	        break
//	    }
//	}
//
// # Thread Safety
//
// Parser instances are not safe for concurrent use. Each goroutine should
// have its own parser instance. The underlying IbtReader can be shared.
//
// # Telemetry Fields
//
// The TelemetryTick struct contains 160+ telemetry fields including:
//
//   - Lap data: LapID, LapDistPct, LapCurrentLapTime, LapDeltaToBestLap
//   - Motion: Speed, VelocityX/Y/Z, LatAccel, LongAccel, VertAccel
//   - Controls: Throttle, Brake, Gear, SteeringWheelAngle, RPM
//   - Position: Lat, Lon, Alt, Pitch, Roll, Yaw
//   - Car state: FuelLevel, WaterTemp, Voltage, tire temps/pressures/wear
//   - Session: SessionTime, SessionNum, PlayerCarIdx
//
// See TelemetryTick struct for complete field list with iRacing field name mappings.
package ibt
