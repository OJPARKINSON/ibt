package ibt

import (
	"encoding/binary"
	"math"
	"sync"
	"time"

	"github.com/OJPARKINSON/ibt/headers"
)

var (
	globalFieldSetters     map[string]fieldSetter
	globalFieldSettersOnce sync.Once
)

func getFieldSetters() map[string]fieldSetter {
	globalFieldSettersOnce.Do(func() { globalFieldSetters = buildFieldSetters() })
	return globalFieldSetters
}

type varSetter struct {
	offset int
	setter fieldSetter
}

// DirectStructParser reads telemetry bytes directly into TelemetryTick structs
// without intermediate map allocations. This provides performance improvements over map-based parsing.
//
// The parser uses safe byte conversions via encoding/binary, which modern
// Go compilers optimise to the same assembly as unsafe pointer casts.
//
// This parser is safe for concurrent use when each goroutine has its own instance.
type DirectStructParser struct {
	header     *headers.Header
	current    int
	mmapData   []byte
	varSetters []varSetter
	maxOffset  int // pre-computed maximum byte offset needed by any setter

	sessionStartTime time.Time
}

// fieldSetter defines how to read a field from a buffer into a TelemetryTick.
type fieldSetter func(tick *TelemetryTick, buf []byte, offset int)

func NewDirectStructParser(reader *IbtReader, header *headers.Header, whitelist ...string) *DirectStructParser {
	varHeaderMap := header.VarHeader
	if varHeaderMap == nil {
		varHeaderMap = make(map[string]headers.VarHeader)
	}

	fieldSetters := getFieldSetters()

	// Build whitelist set for filtering
	useWhitelist := len(whitelist) > 0 && (len(whitelist) != 1 || whitelist[0] != "*")
	var whitelistSet map[string]bool
	if useWhitelist {
		whitelistSet = make(map[string]bool, len(whitelist))
		for _, name := range whitelist {
			whitelistSet[name] = true
		}
	}

	// Build pre-resolved varSetters slice: one entry per var that has a setter
	setters := make([]varSetter, 0, len(varHeaderMap))
	maxOffset := 0
	for name, vh := range varHeaderMap {
		if useWhitelist && !whitelistSet[name] {
			continue
		}
		if setter, ok := fieldSetters[name]; ok {
			setters = append(setters, varSetter{offset: vh.Offset, setter: setter})
			// Track the highest byte offset needed.
			// float64 fields read 8 bytes; all others read <=4.
			// Use offset+8 as conservative upper bound for all types.
			end := vh.Offset + 8
			if end > maxOffset {
				maxOffset = end
			}
		}
	}

	var mmapData []byte
	if reader != nil {
		mmapData = reader.Data()
	}

	var sessionStartTime time.Time
	if header.DiskHeader != nil {
		sessionStartTime = time.Unix(header.DiskHeader.StartDate, 0)
	}

	return &DirectStructParser{
		header:           header,
		current:          0,
		mmapData:         mmapData,
		varSetters:       setters,
		maxOffset:        maxOffset,
		sessionStartTime: sessionStartTime,
	}
}

// read returns a buffer at the given offset.
func (p *DirectStructParser) read(start int) []byte {
	bufLen := p.header.TelemetryHeader.BufLen
	end := start + bufLen

	if end > len(p.mmapData) {
		return nil
	}

	return p.mmapData[start:end]
}

// NextStruct reads the next telemetry tick directly into a struct.
// Returns (nil, false) when end of file is reached.
func (p *DirectStructParser) NextStruct(tick *TelemetryTick) (*TelemetryTick, bool) {
	start := p.header.TelemetryHeader.BufOffset + (p.current * p.header.TelemetryHeader.BufLen)

	currentBuf := p.read(start)
	if currentBuf == nil {
		return nil, false
	}

	// Check if more data available
	nextStart := p.header.TelemetryHeader.BufOffset + ((p.current + 1) * p.header.TelemetryHeader.BufLen)
	hasNext := nextStart+p.header.TelemetryHeader.BufLen <= len(p.mmapData)

	// Single bounds check: all offsets were pre-validated at construction,
	// so we only need to verify the buffer is large enough once.
	if len(currentBuf) < p.maxOffset {
		return nil, false
	}

	// Populate fields using pre-resolved setters (no map lookup per tick).
	// Individual setters skip bounds checks since we validated above.
	for _, vs := range p.varSetters {
		vs.setter(tick, currentBuf, vs.offset)
	}

	// Calculate absolute timestamp: session start + elapsed session time
	// SessionTime is in seconds (float64), so convert to Duration and add to start time
	tick.TickTime = p.sessionStartTime.Add(time.Duration(tick.SessionTime * float64(time.Second)))

	p.current++

	return tick, hasNext
}

// buildFieldSetters creates the field mapping table using safe byte conversions.
// Each setter reads bytes from the buffer and assigns to the appropriate struct field.
func buildFieldSetters() map[string]fieldSetter {
	return map[string]fieldSetter{
		// Integer fields - Lap & Position
		"Lap":                    setInt32(func(t *TelemetryTick, v int32) { t.LapID = v }),
		"PlayerCarClassPosition": setInt32(func(t *TelemetryTick, v int32) { t.PlayerCarClassPosition = v }),
		"PlayerCarIdx":           setInt32(func(t *TelemetryTick, v int32) { t.PlayerCarIdx = v }),

		// Integer fields - Lap Timing
		"LapBestLap":     setInt32(func(t *TelemetryTick, v int32) { t.LapBestLap = v }),
		"LapBestNLapLap": setInt32(func(t *TelemetryTick, v int32) { t.LapBestNLapLap = v }),
		"LapLasNLapSeq":  setInt32(func(t *TelemetryTick, v int32) { t.LapLasNLapSeq = v }),

		// Integer fields - Session
		"SessionNum":        setInt32(func(t *TelemetryTick, v int32) { t.SessionNum = v }),
		"SessionLapsRemain": setInt32(func(t *TelemetryTick, v int32) { t.SessionLapsRemain = v }),
		"SessionState":      setInt32(func(t *TelemetryTick, v int32) { t.SessionState = v }),
		"SessionUniqueID":   setInt32(func(t *TelemetryTick, v int32) { t.SessionUniqueID = v }),

		// Integer fields - Environment & System
		"Skies":          setInt32(func(t *TelemetryTick, v int32) { t.Skies = v }),
		"WeatherType":    setInt32(func(t *TelemetryTick, v int32) { t.WeatherType = v }),
		"EnterExitReset": setInt32(func(t *TelemetryTick, v int32) { t.EnterExitReset = v }),

		// Uint32 fields
		"Gear": setUint32(func(t *TelemetryTick, v uint32) { t.Gear = v }),

		// Bool fields - Delta Timing
		"LapDeltaToBestLap_OK":           setBool(func(t *TelemetryTick, v bool) { t.LapDeltaToBestLap_OK = v }),
		"LapDeltaToOptimalLap_OK":        setBool(func(t *TelemetryTick, v bool) { t.LapDeltaToOptimalLap_OK = v }),
		"LapDeltaToSessionBestLap_OK":    setBool(func(t *TelemetryTick, v bool) { t.LapDeltaToSessionBestLap_OK = v }),
		"LapDeltaToSessionLastlLap_OK":   setBool(func(t *TelemetryTick, v bool) { t.LapDeltaToSessionLastlLap_OK = v }),
		"LapDeltaToSessionOptimalLap_OK": setBool(func(t *TelemetryTick, v bool) { t.LapDeltaToSessionOptimalLap_OK = v }),

		// Bool fields - Track Position & System
		"IsOnTrack":               setBool(func(t *TelemetryTick, v bool) { t.IsOnTrack = v }),
		"IsOnTrackCar":            setBool(func(t *TelemetryTick, v bool) { t.IsOnTrackCar = v }),
		"OnPitRoad":               setBool(func(t *TelemetryTick, v bool) { t.OnPitRoad = v }),
		"DriverMarker":            setBool(func(t *TelemetryTick, v bool) { t.DriverMarker = v }),
		"dcTractionControlToggle": setBool(func(t *TelemetryTick, v bool) { t.DcTractionControlToggle = v }),

		// Float32 fields (stored as float64) - Lap & Position
		"LapDistPct":        setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.LapDistPct = v }),
		"LapDist":           setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.LapDist = v }),
		"PlayerCarPosition": setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.PlayerCarPosition = v }),

		// Float32 fields - Lap Timing
		"LapCurrentLapTime": setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.LapCurrentLapTime = v }),
		"LapLastLapTime":    setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.LapLastLapTime = v }),
		"LapLastNLapTime":   setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.LapLastNLapTime = v }),
		"LapBestLapTime":    setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.LapBestLapTime = v }),
		"LapBestNLapTime":   setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.LapBestNLapTime = v }),

		// Float32 fields - Delta Timing
		"LapDeltaToBestLap":              setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.LapDeltaToBestLap = v }),
		"LapDeltaToBestLap_DD":           setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.LapDeltaToBestLap_DD = v }),
		"LapDeltaToOptimalLap":           setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.LapDeltaToOptimalLap = v }),
		"LapDeltaToOptimalLap_DD":        setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.LapDeltaToOptimalLap_DD = v }),
		"LapDeltaToSessionBestLap":       setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.LapDeltaToSessionBestLap = v }),
		"LapDeltaToSessionBestLap_DD":    setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.LapDeltaToSessionBestLap_DD = v }),
		"LapDeltaToSessionLastlLap":      setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.LapDeltaToSessionLastlLap = v }),
		"LapDeltaToSessionLastlLap_DD":   setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.LapDeltaToSessionLastlLap_DD = v }),
		"LapDeltaToSessionOptimalLap":    setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.LapDeltaToSessionOptimalLap = v }),
		"LapDeltaToSessionOptimalLap_DD": setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.LapDeltaToSessionOptimalLap_DD = v }),

		// Float32 fields - Driver Inputs
		"Throttle":                        setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.Throttle = v }),
		"ThrottleRaw":                     setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.ThrottleRaw = v }),
		"Brake":                           setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.Brake = v }),
		"BrakeRaw":                        setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.BrakeRaw = v }),
		"Clutch":                          setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.Clutch = v }),
		"SteeringWheelAngle":              setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.SteeringWheelAngle = v }),
		"SteeringWheelAngleMax":           setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.SteeringWheelAngleMax = v }),
		"SteeringWheelTorque":             setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.SteeringWheelTorque = v }),
		"SteeringWheelPctTorque":          setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.SteeringWheelPctTorque = v }),
		"SteeringWheelPctTorqueSign":      setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.SteeringWheelPctTorqueSign = v }),
		"SteeringWheelPctTorqueSignStops": setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.SteeringWheelPctTorqueSignStops = v }),
		"SteeringWheelPctDamper":          setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.SteeringWheelPctDamper = v }),

		// Float32 fields - Speed & Motion
		"Speed":     setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.Speed = v }),
		"VelocityX": setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.VelocityX = v }),
		"VelocityY": setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.VelocityY = v }),
		"VelocityZ": setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.VelocityZ = v }),

		// Float32 fields - Acceleration
		"LatAccel":  setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.LatAccel = v }),
		"LongAccel": setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.LongAccel = v }),
		"VertAccel": setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.VertAccel = v }),

		// Float32 fields - Orientation
		"pitch":     setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.Pitch = v }),
		"PitchRate": setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.PitchRate = v }),
		"roll":      setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.Roll = v }),
		"RollRate":  setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.RollRate = v }),
		"yaw":       setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.Yaw = v }),
		"YawNorth":  setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.YawNorth = v }),
		"YawRate":   setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.YawRate = v }),

		// Float32 fields - GPS Position (Lat/Lon are actually double precision)
		"alt": setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.Alt = v }),

		// Float32 fields - Engine
		"RPM":               setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.RPM = v }),
		"ShiftGrindRPM":     setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.ShiftGrindRPM = v }),
		"ShiftIndicatorPct": setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.ShiftIndicatorPct = v }),
		"ShiftPowerPct":     setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.ShiftPowerPct = v }),
		"Voltage":           setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.Voltage = v }),
		"WaterTemp":         setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.WaterTemp = v }),
		"WaterLevel":        setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.WaterLevel = v }),
		"OilTemp":           setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.OilTemp = v }),
		"OilPress":          setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.OilPress = v }),
		"OilLevel":          setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.OilLevel = v }),
		"FuelPress":         setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.FuelPress = v }),
		"ManifoldPress":     setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.ManifoldPress = v }),

		// Float32 fields - Fuel
		"FuelLevel":      setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.FuelLevel = v }),
		"FuelLevelPct":   setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.FuelLevelPct = v }),
		"FuelUsePerHour": setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.FuelUsePerHour = v }),

		// Float32 fields - Environment
		"AirDensity":       setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.AirDensity = v }),
		"AirPressure":      setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.AirPressure = v }),
		"AirTemp":          setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.AirTemp = v }),
		"FogLevel":         setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.FogLevel = v }),
		"RelativeHumidity": setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.RelativeHumidity = v }),
		"TrackTemp":        setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.TrackTemp = v }),
		"TrackTempCrew":    setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.TrackTempCrew = v }),
		"WindDir":          setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.WindDir = v }),
		"WindVel":          setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.WindVel = v }),

		// Float32 fields - System
		"CpuUsageBG": setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.CpuUsageBG = v }),
		"FrameRate":  setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.FrameRate = v }),

		// Float32 fields - Pit Service
		"PitRepairLeft":    setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.PitRepairLeft = v }),
		"PitOptRepairLeft": setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.PitOptRepairLeft = v }),
		"PitSvFuel":        setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.PitSvFuel = v }),
		"PitSvLFP":         setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.PitSvLFP = v }),
		"PitSvRFP":         setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.PitSvRFP = v }),
		"PitSvLRP":         setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.PitSvLRP = v }),
		"PitSvRRP":         setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.PitSvRRP = v }),

		// Float32 fields - Tire Pressure
		"LFpressure": setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.LFpressure = v }),
		"RFpressure": setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.RFpressure = v }),
		"LRpressure": setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.LRpressure = v }),
		"RRpressure": setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.RRpressure = v }),

		// Float32 fields - Tire Temperature (Middle/Surface)
		"LFtempM": setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.LFtempM = v }),
		"RFtempM": setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.RFtempM = v }),
		"LRtempM": setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.LRtempM = v }),
		"RRtempM": setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.RRtempM = v }),

		// Float32 fields - Tire Temperature (Carcass)
		"LFtempCL": setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.LFtempCL = v }),
		"LFtempCM": setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.LFtempCM = v }),
		"LFtempCR": setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.LFtempCR = v }),
		"RFtempCL": setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.RFtempCL = v }),
		"RFtempCM": setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.RFtempCM = v }),
		"RFtempCR": setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.RFtempCR = v }),
		"LRtempCL": setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.LRtempCL = v }),
		"LRtempCM": setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.LRtempCM = v }),
		"LRtempCR": setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.LRtempCR = v }),
		"RRtempCL": setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.RRtempCL = v }),
		"RRtempCM": setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.RRtempCM = v }),
		"RRtempCR": setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.RRtempCR = v }),

		// Float32 fields - Tire Wear
		"LFwearL": setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.LFwearL = v }),
		"LFwearM": setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.LFwearM = v }),
		"LFwearR": setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.LFwearR = v }),
		"RFwearL": setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.RFwearL = v }),
		"RFwearM": setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.RFwearM = v }),
		"RFwearR": setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.RFwearR = v }),
		"LRwearL": setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.LRwearL = v }),
		"LRwearM": setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.LRwearM = v }),
		"LRwearR": setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.LRwearR = v }),
		"RRwearL": setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.RRwearL = v }),
		"RRwearM": setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.RRwearM = v }),
		"RRwearR": setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.RRwearR = v }),

		// Float32 fields - Tire Cold Pressure
		"LFcoldPressure": setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.LFcoldPressure = v }),
		"RFcoldPressure": setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.RFcoldPressure = v }),
		"LRcoldPressure": setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.LRcoldPressure = v }),
		"RRcoldPressure": setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.RRcoldPressure = v }),

		// Float32 fields - Suspension - Shock Deflection
		"LFshockDefl": setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.LFshockDefl = v }),
		"RFshockDefl": setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.RFshockDefl = v }),
		"LRshockDefl": setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.LRshockDefl = v }),
		"RRshockDefl": setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.RRshockDefl = v }),
		"CFshockDefl": setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.CFshockDefl = v }),
		"CRshockDefl": setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.CRshockDefl = v }),

		// Float32 fields - Suspension - Shock Velocity
		"LFshockVel": setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.LFshockVel = v }),
		"RFshockVel": setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.RFshockVel = v }),
		"LRshockVel": setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.LRshockVel = v }),
		"RRshockVel": setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.RRshockVel = v }),
		"CFshockVel": setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.CFshockVel = v }),
		"CRshockVel": setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.CRshockVel = v }),

		// Float32 fields - Wheel Speed
		"LFspeed": setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.LFspeed = v }),
		"RFspeed": setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.RFspeed = v }),
		"LRspeed": setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.LRspeed = v }),
		"RRspeed": setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.RRspeed = v }),

		// Float32 fields - Brake Line Pressure
		"LFbrakeLinePress": setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.LFbrakeLinePress = v }),
		"RFbrakeLinePress": setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.RFbrakeLinePress = v }),
		"LRbrakeLinePress": setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.LRbrakeLinePress = v }),
		"RRbrakeLinePress": setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.RRbrakeLinePress = v }),

		// Float32 fields - Dynamic In-Car Adjustments
		"dcABS":               setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.DcABS = v }),
		"dcAntiRollFront":     setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.DcAntiRollFront = v }),
		"dcAntiRollRear":      setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.DcAntiRollRear = v }),
		"dcBoostLevel":        setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.DcBoostLevel = v }),
		"dcBrakeBias":         setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.DcBrakeBias = v }),
		"dcDiffEntry":         setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.DcDiffEntry = v }),
		"dcDiffExit":          setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.DcDiffExit = v }),
		"dcDiffMiddle":        setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.DcDiffMiddle = v }),
		"dcEngineBraking":     setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.DcEngineBraking = v }),
		"dcEnginePower":       setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.DcEnginePower = v }),
		"dcFuelMixture":       setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.DcFuelMixture = v }),
		"dcRevLimiter":        setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.DcRevLimiter = v }),
		"dcThrottleShape":     setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.DcThrottleShape = v }),
		"dcTractionControl":   setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.DcTractionControl = v }),
		"dcTractionControl2":  setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.DcTractionControl2 = v }),
		"dcWeightJackerLeft":  setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.DcWeightJackerLeft = v }),
		"dcWeightJackerRight": setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.DcWeightJackerRight = v }),
		"dcWingFront":         setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.DcWingFront = v }),
		"dcWingRear":          setFloat32AsFloat64(func(t *TelemetryTick, v float64) { t.DcWingRear = v }),

		// Float64 fields (GPS coordinates and session time use double precision)
		"SessionTime":       setFloat64(func(t *TelemetryTick, v float64) { t.SessionTime = v }),
		"SessionTimeRemain": setFloat64(func(t *TelemetryTick, v float64) { t.SessionTimeRemain = v }),
		"Lat":               setFloat64(func(t *TelemetryTick, v float64) { t.Lat = v }),
		"Lon":               setFloat64(func(t *TelemetryTick, v float64) { t.Lon = v }),
	}
}

// Helper functions to create field setters with proper type conversions.
// These use safe byte conversions via encoding/binary, which modern Go compilers
// optimise to the same assembly as unsafe pointer casts with zero performance cost.

// setInt32 creates a setter for int32 fields.
// Bounds checking is done once in NextStruct() before the setter loop.
func setInt32(assign func(*TelemetryTick, int32)) fieldSetter {
	return func(tick *TelemetryTick, buf []byte, offset int) {
		val := int32(binary.LittleEndian.Uint32(buf[offset : offset+4]))
		assign(tick, val)
	}
}

// setUint32 creates a setter for uint32 fields.
func setUint32(assign func(*TelemetryTick, uint32)) fieldSetter {
	return func(tick *TelemetryTick, buf []byte, offset int) {
		val := binary.LittleEndian.Uint32(buf[offset : offset+4])
		assign(tick, val)
	}
}

// setFloat32AsFloat64 creates a setter for float32 fields stored as float64.
func setFloat32AsFloat64(assign func(*TelemetryTick, float64)) fieldSetter {
	return func(tick *TelemetryTick, buf []byte, offset int) {
		bits := binary.LittleEndian.Uint32(buf[offset : offset+4])
		val := float64(math.Float32frombits(bits))
		assign(tick, val)
	}
}

// setFloat64 creates a setter for float64 fields.
func setFloat64(assign func(*TelemetryTick, float64)) fieldSetter {
	return func(tick *TelemetryTick, buf []byte, offset int) {
		bits := binary.LittleEndian.Uint64(buf[offset : offset+8])
		val := math.Float64frombits(bits)
		assign(tick, val)
	}
}

// setBool creates a setter for bool fields.
func setBool(assign func(*TelemetryTick, bool)) fieldSetter {
	return func(tick *TelemetryTick, buf []byte, offset int) {
		val := buf[offset] != 0
		assign(tick, val)
	}
}
