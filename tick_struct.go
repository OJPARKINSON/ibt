package ibt

import "time"

type TelemetryTick struct {
	LapID              int32   `ibt:"Lap"`
	Speed              float64 `ibt:"Speed"`
	LapDistPct         float64 `ibt:"LapDistPct"`
	Throttle           float64 `ibt:"Throttle"`
	Brake              float64 `ibt:"Brake"`
	Gear               uint32  `ibt:"Gear"`
	RPM                float64 `ibt:"RPM"`
	SteeringWheelAngle float64 `ibt:"SteeringWheelAngle"`
	VelocityX          float64 `ibt:"VelocityX"`
	VelocityY          float64 `ibt:"VelocityY"`
	VelocityZ          float64 `ibt:"VelocityZ"`
	Lat                float64 `ibt:"Lat"`
	Lon                float64 `ibt:"Lon"`
	SessionTime        float64 `ibt:"SessionTime"`
	PlayerCarPosition  float64 `ibt:"PlayerCarPosition"`
	FuelLevel          float64 `ibt:"FuelLevel"`
	PlayerCarIdx       int32   `ibt:"PlayerCarIdx"`
	SessionNum         int32   `ibt:"SessionNum"`
	Alt                float64 `ibt:"alt"`
	LatAccel           float64 `ibt:"LatAccel"`
	LongAccel          float64 `ibt:"LongAccel"`
	VertAccel          float64 `ibt:"VertAccel"`
	Pitch              float64 `ibt:"pitch"`
	Roll               float64 `ibt:"roll"`
	Yaw                float64 `ibt:"yaw"`
	YawNorth           float64 `ibt:"YawNorth"`
	Voltage            float64 `ibt:"Voltage"`
	LapLastLapTime     float64 `ibt:"LapLastLapTime"`
	WaterTemp          float64 `ibt:"WaterTemp"`
	LapDeltaToBestLap  float64 `ibt:"LapDeltaToBestLap"`
	LapCurrentLapTime  float64 `ibt:"LapCurrentLapTime"`
	LFpressure         float64 `ibt:"LFpressure"`
	RFpressure         float64 `ibt:"RFpressure"`
	LRpressure         float64 `ibt:"LRpressure"`
	RRpressure         float64 `ibt:"RRpressure"`
	LFtempM            float64 `ibt:"LFtempM"`
	RFtempM            float64 `ibt:"RFtempM"`
	LRtempM            float64 `ibt:"LRtempM"`
	RRtempM            float64 `ibt:"RRtempM"`

	SessionID   string
	SessionType string
	SessionName string
	TrackName   string
	TrackID     int
	WorkerID    int
	GroupNum    int

	TickTime time.Time
}

func buildFieldMap() map[string]func(*TelemetryTick, interface{}) {
	// Helper to convert both float32 and float64 to float64
	toFloat64 := func(v interface{}) float64 {
		switch val := v.(type) {
		case float64:
			return val
		case float32:
			return float64(val)
		default:
			return 0
		}
	}

	// Helper to convert to int32
	toInt32 := func(v interface{}) int32 {
		switch val := v.(type) {
		case int32:
			return val
		case int:
			return int32(val)
		default:
			return 0
		}
	}

	// Helper to convert to uint32
	toUint32 := func(v interface{}) uint32 {
		switch val := v.(type) {
		case uint32:
			return val
		case int:
			return uint32(val)
		case int32:
			return uint32(val)
		default:
			return 0
		}
	}

	return map[string]func(*TelemetryTick, interface{}){
		"Lap":                func(t *TelemetryTick, v interface{}) { t.LapID = toInt32(v) },
		"Speed":              func(t *TelemetryTick, v interface{}) { t.Speed = toFloat64(v) },
		"LapDistPct":         func(t *TelemetryTick, v interface{}) { t.LapDistPct = toFloat64(v) },
		"Throttle":           func(t *TelemetryTick, v interface{}) { t.Throttle = toFloat64(v) },
		"Brake":              func(t *TelemetryTick, v interface{}) { t.Brake = toFloat64(v) },
		"Gear":               func(t *TelemetryTick, v interface{}) { t.Gear = toUint32(v) },
		"RPM":                func(t *TelemetryTick, v interface{}) { t.RPM = toFloat64(v) },
		"SteeringWheelAngle": func(t *TelemetryTick, v interface{}) { t.SteeringWheelAngle = toFloat64(v) },
		"VelocityX":          func(t *TelemetryTick, v interface{}) { t.VelocityX = toFloat64(v) },
		"VelocityY":          func(t *TelemetryTick, v interface{}) { t.VelocityY = toFloat64(v) },
		"VelocityZ":          func(t *TelemetryTick, v interface{}) { t.VelocityZ = toFloat64(v) },
		"Lat":                func(t *TelemetryTick, v interface{}) { t.Lat = toFloat64(v) },
		"Lon":                func(t *TelemetryTick, v interface{}) { t.Lon = toFloat64(v) },
		"SessionTime":        func(t *TelemetryTick, v interface{}) { t.SessionTime = toFloat64(v) },
		"PlayerCarPosition":  func(t *TelemetryTick, v interface{}) { t.PlayerCarPosition = toFloat64(v) },
		"FuelLevel":          func(t *TelemetryTick, v interface{}) { t.FuelLevel = toFloat64(v) },
		"PlayerCarIdx":       func(t *TelemetryTick, v interface{}) { t.PlayerCarIdx = toInt32(v) },
		"SessionNum":         func(t *TelemetryTick, v interface{}) { t.SessionNum = toInt32(v) },
		"alt":                func(t *TelemetryTick, v interface{}) { t.Alt = toFloat64(v) },
		"LatAccel":           func(t *TelemetryTick, v interface{}) { t.LatAccel = toFloat64(v) },
		"LongAccel":          func(t *TelemetryTick, v interface{}) { t.LongAccel = toFloat64(v) },
		"VertAccel":          func(t *TelemetryTick, v interface{}) { t.VertAccel = toFloat64(v) },
		"pitch":              func(t *TelemetryTick, v interface{}) { t.Pitch = toFloat64(v) },
		"roll":               func(t *TelemetryTick, v interface{}) { t.Roll = toFloat64(v) },
		"yaw":                func(t *TelemetryTick, v interface{}) { t.Yaw = toFloat64(v) },
		"YawNorth":           func(t *TelemetryTick, v interface{}) { t.YawNorth = toFloat64(v) },
		"Voltage":            func(t *TelemetryTick, v interface{}) { t.Voltage = toFloat64(v) },
		"LapLastLapTime":     func(t *TelemetryTick, v interface{}) { t.LapLastLapTime = toFloat64(v) },
		"WaterTemp":          func(t *TelemetryTick, v interface{}) { t.WaterTemp = toFloat64(v) },
		"LapDeltaToBestLap":  func(t *TelemetryTick, v interface{}) { t.LapDeltaToBestLap = toFloat64(v) },
		"LapCurrentLapTime":  func(t *TelemetryTick, v interface{}) { t.LapCurrentLapTime = toFloat64(v) },
		"LFpressure":         func(t *TelemetryTick, v interface{}) { t.LFpressure = toFloat64(v) },
		"RFpressure":         func(t *TelemetryTick, v interface{}) { t.RFpressure = toFloat64(v) },
		"LRpressure":         func(t *TelemetryTick, v interface{}) { t.LRpressure = toFloat64(v) },
		"RRpressure":         func(t *TelemetryTick, v interface{}) { t.RRpressure = toFloat64(v) },
		"LFtempM":            func(t *TelemetryTick, v interface{}) { t.LFtempM = toFloat64(v) },
		"RFtempM":            func(t *TelemetryTick, v interface{}) { t.RFtempM = toFloat64(v) },
		"LRtempM":            func(t *TelemetryTick, v interface{}) { t.LRtempM = toFloat64(v) },
		"RRtempM":            func(t *TelemetryTick, v interface{}) { t.RRtempM = toFloat64(v) },
	}
}
