package ibt

import (
	"reflect"
	"time"
)

type TelemetryTick struct {
	// Lap & Position
	LapID              int32   `ibt:"Lap"`
	LapDistPct         float64 `ibt:"LapDistPct"`
	LapDist            float64 `ibt:"LapDist"`
	PlayerCarPosition  float64 `ibt:"PlayerCarPosition"`
	PlayerCarClassPosition int32 `ibt:"PlayerCarClassPosition"`
	PlayerCarIdx       int32   `ibt:"PlayerCarIdx"`

	// Lap Timing
	LapCurrentLapTime  float64 `ibt:"LapCurrentLapTime"`
	LapLastLapTime     float64 `ibt:"LapLastLapTime"`
	LapLastNLapTime    float64 `ibt:"LapLastNLapTime"`
	LapBestLap         int32   `ibt:"LapBestLap"`
	LapBestLapTime     float64 `ibt:"LapBestLapTime"`
	LapBestNLapLap     int32   `ibt:"LapBestNLapLap"`
	LapBestNLapTime    float64 `ibt:"LapBestNLapTime"`
	LapLasNLapSeq      int32   `ibt:"LapLasNLapSeq"`

	// Delta Timing
	LapDeltaToBestLap  float64 `ibt:"LapDeltaToBestLap"`
	LapDeltaToBestLap_DD float64 `ibt:"LapDeltaToBestLap_DD"`
	LapDeltaToBestLap_OK bool   `ibt:"LapDeltaToBestLap_OK"`
	LapDeltaToOptimalLap float64 `ibt:"LapDeltaToOptimalLap"`
	LapDeltaToOptimalLap_DD float64 `ibt:"LapDeltaToOptimalLap_DD"`
	LapDeltaToOptimalLap_OK bool `ibt:"LapDeltaToOptimalLap_OK"`
	LapDeltaToSessionBestLap float64 `ibt:"LapDeltaToSessionBestLap"`
	LapDeltaToSessionBestLap_DD float64 `ibt:"LapDeltaToSessionBestLap_DD"`
	LapDeltaToSessionBestLap_OK bool `ibt:"LapDeltaToSessionBestLap_OK"`
	LapDeltaToSessionLastlLap float64 `ibt:"LapDeltaToSessionLastlLap"`
	LapDeltaToSessionLastlLap_DD float64 `ibt:"LapDeltaToSessionLastlLap_DD"`
	LapDeltaToSessionLastlLap_OK bool `ibt:"LapDeltaToSessionLastlLap_OK"`
	LapDeltaToSessionOptimalLap float64 `ibt:"LapDeltaToSessionOptimalLap"`
	LapDeltaToSessionOptimalLap_DD float64 `ibt:"LapDeltaToSessionOptimalLap_DD"`
	LapDeltaToSessionOptimalLap_OK bool `ibt:"LapDeltaToSessionOptimalLap_OK"`

	// Session Info
	SessionNum         int32   `ibt:"SessionNum"`
	SessionTime        float64 `ibt:"SessionTime"`
	SessionTimeRemain  float64 `ibt:"SessionTimeRemain"`
	SessionLapsRemain  int32   `ibt:"SessionLapsRemain"`
	SessionState       int32   `ibt:"SessionState"`
	SessionUniqueID    int32   `ibt:"SessionUniqueID"`

	// Driver Inputs
	Throttle           float64 `ibt:"Throttle"`
	ThrottleRaw        float64 `ibt:"ThrottleRaw"`
	Brake              float64 `ibt:"Brake"`
	BrakeRaw           float64 `ibt:"BrakeRaw"`
	Clutch             float64 `ibt:"Clutch"`
	Gear               uint32  `ibt:"Gear"`
	SteeringWheelAngle float64 `ibt:"SteeringWheelAngle"`
	SteeringWheelAngleMax float64 `ibt:"SteeringWheelAngleMax"`
	SteeringWheelTorque float64 `ibt:"SteeringWheelTorque"`
	SteeringWheelPctTorque float64 `ibt:"SteeringWheelPctTorque"`
	SteeringWheelPctTorqueSign float64 `ibt:"SteeringWheelPctTorqueSign"`
	SteeringWheelPctTorqueSignStops float64 `ibt:"SteeringWheelPctTorqueSignStops"`
	SteeringWheelPctDamper float64 `ibt:"SteeringWheelPctDamper"`

	// Speed & Motion
	Speed              float64 `ibt:"Speed"`
	VelocityX          float64 `ibt:"VelocityX"`
	VelocityY          float64 `ibt:"VelocityY"`
	VelocityZ          float64 `ibt:"VelocityZ"`

	// Acceleration
	LatAccel           float64 `ibt:"LatAccel"`
	LongAccel          float64 `ibt:"LongAccel"`
	VertAccel          float64 `ibt:"VertAccel"`

	// Orientation
	Pitch              float64 `ibt:"pitch"`
	PitchRate          float64 `ibt:"PitchRate"`
	Roll               float64 `ibt:"roll"`
	RollRate           float64 `ibt:"RollRate"`
	Yaw                float64 `ibt:"yaw"`
	YawNorth           float64 `ibt:"YawNorth"`
	YawRate            float64 `ibt:"YawRate"`

	// GPS Position
	Lat                float64 `ibt:"Lat"`
	Lon                float64 `ibt:"Lon"`
	Alt                float64 `ibt:"alt"`

	// Engine
	RPM                float64 `ibt:"RPM"`
	ShiftGrindRPM      float64 `ibt:"ShiftGrindRPM"`
	ShiftIndicatorPct  float64 `ibt:"ShiftIndicatorPct"`
	ShiftPowerPct      float64 `ibt:"ShiftPowerPct"`
	Voltage            float64 `ibt:"Voltage"`
	WaterTemp          float64 `ibt:"WaterTemp"`
	WaterLevel         float64 `ibt:"WaterLevel"`
	OilTemp            float64 `ibt:"OilTemp"`
	OilPress           float64 `ibt:"OilPress"`
	OilLevel           float64 `ibt:"OilLevel"`
	FuelPress          float64 `ibt:"FuelPress"`
	ManifoldPress      float64 `ibt:"ManifoldPress"`

	// Fuel
	FuelLevel          float64 `ibt:"FuelLevel"`
	FuelLevelPct       float64 `ibt:"FuelLevelPct"`
	FuelUsePerHour     float64 `ibt:"FuelUsePerHour"`

	// Environment
	AirDensity         float64 `ibt:"AirDensity"`
	AirPressure        float64 `ibt:"AirPressure"`
	AirTemp            float64 `ibt:"AirTemp"`
	FogLevel           float64 `ibt:"FogLevel"`
	RelativeHumidity   float64 `ibt:"RelativeHumidity"`
	Skies              int32   `ibt:"Skies"`
	TrackTemp          float64 `ibt:"TrackTemp"`
	TrackTempCrew      float64 `ibt:"TrackTempCrew"`
	WeatherType        int32   `ibt:"WeatherType"`
	WindDir            float64 `ibt:"WindDir"`
	WindVel            float64 `ibt:"WindVel"`

	// Track Position
	IsOnTrack          bool    `ibt:"IsOnTrack"`
	IsOnTrackCar       bool    `ibt:"IsOnTrackCar"`
	OnPitRoad          bool    `ibt:"OnPitRoad"`

	// System
	CpuUsageBG         float64 `ibt:"CpuUsageBG"`
	FrameRate          float64 `ibt:"FrameRate"`
	DriverMarker       bool    `ibt:"DriverMarker"`
	EnterExitReset     int32   `ibt:"EnterExitReset"`

	// Pit Service
	PitRepairLeft      float64 `ibt:"PitRepairLeft"`
	PitOptRepairLeft   float64 `ibt:"PitOptRepairLeft"`
	PitSvFuel          float64 `ibt:"PitSvFuel"`
	PitSvLFP           float64 `ibt:"PitSvLFP"`
	PitSvRFP           float64 `ibt:"PitSvRFP"`
	PitSvLRP           float64 `ibt:"PitSvLRP"`
	PitSvRRP           float64 `ibt:"PitSvRRP"`

	// Tire Pressure
	LFpressure         float64 `ibt:"LFpressure"`
	RFpressure         float64 `ibt:"RFpressure"`
	LRpressure         float64 `ibt:"LRpressure"`
	RRpressure         float64 `ibt:"RRpressure"`

	// Tire Temperature (Middle/Surface)
	LFtempM            float64 `ibt:"LFtempM"`
	RFtempM            float64 `ibt:"RFtempM"`
	LRtempM            float64 `ibt:"LRtempM"`
	RRtempM            float64 `ibt:"RRtempM"`

	// Tire Temperature (Carcass - Left, Middle, Right)
	LFtempCL           float64 `ibt:"LFtempCL"`
	LFtempCM           float64 `ibt:"LFtempCM"`
	LFtempCR           float64 `ibt:"LFtempCR"`
	RFtempCL           float64 `ibt:"RFtempCL"`
	RFtempCM           float64 `ibt:"RFtempCM"`
	RFtempCR           float64 `ibt:"RFtempCR"`
	LRtempCL           float64 `ibt:"LRtempCL"`
	LRtempCM           float64 `ibt:"LRtempCM"`
	LRtempCR           float64 `ibt:"LRtempCR"`
	RRtempCL           float64 `ibt:"RRtempCL"`
	RRtempCM           float64 `ibt:"RRtempCM"`
	RRtempCR           float64 `ibt:"RRtempCR"`

	// Tire Wear
	LFwearL            float64 `ibt:"LFwearL"`
	LFwearM            float64 `ibt:"LFwearM"`
	LFwearR            float64 `ibt:"LFwearR"`
	RFwearL            float64 `ibt:"RFwearL"`
	RFwearM            float64 `ibt:"RFwearM"`
	RFwearR            float64 `ibt:"RFwearR"`
	LRwearL            float64 `ibt:"LRwearL"`
	LRwearM            float64 `ibt:"LRwearM"`
	LRwearR            float64 `ibt:"LRwearR"`
	RRwearL            float64 `ibt:"RRwearL"`
	RRwearM            float64 `ibt:"RRwearM"`
	RRwearR            float64 `ibt:"RRwearR"`

	// Tire Cold Pressure
	LFcoldPressure     float64 `ibt:"LFcoldPressure"`
	RFcoldPressure     float64 `ibt:"RFcoldPressure"`
	LRcoldPressure     float64 `ibt:"LRcoldPressure"`
	RRcoldPressure     float64 `ibt:"RRcoldPressure"`

	// Suspension - Shock Deflection
	LFshockDefl        float64 `ibt:"LFshockDefl"`
	RFshockDefl        float64 `ibt:"RFshockDefl"`
	LRshockDefl        float64 `ibt:"LRshockDefl"`
	RRshockDefl        float64 `ibt:"RRshockDefl"`
	CFshockDefl        float64 `ibt:"CFshockDefl"`
	CRshockDefl        float64 `ibt:"CRshockDefl"`

	// Suspension - Shock Velocity
	LFshockVel         float64 `ibt:"LFshockVel"`
	RFshockVel         float64 `ibt:"RFshockVel"`
	LRshockVel         float64 `ibt:"LRshockVel"`
	RRshockVel         float64 `ibt:"RRshockVel"`
	CFshockVel         float64 `ibt:"CFshockVel"`
	CRshockVel         float64 `ibt:"CRshockVel"`

	// Wheel Speed
	LFspeed            float64 `ibt:"LFspeed"`
	RFspeed            float64 `ibt:"RFspeed"`
	LRspeed            float64 `ibt:"LRspeed"`
	RRspeed            float64 `ibt:"RRspeed"`

	// Brake Line Pressure
	LFbrakeLinePress   float64 `ibt:"LFbrakeLinePress"`
	RFbrakeLinePress   float64 `ibt:"RFbrakeLinePress"`
	LRbrakeLinePress   float64 `ibt:"LRbrakeLinePress"`
	RRbrakeLinePress   float64 `ibt:"RRbrakeLinePress"`

	// Dynamic In-Car Adjustments
	DcABS              float64 `ibt:"dcABS"`
	DcAntiRollFront    float64 `ibt:"dcAntiRollFront"`
	DcAntiRollRear     float64 `ibt:"dcAntiRollRear"`
	DcBoostLevel       float64 `ibt:"dcBoostLevel"`
	DcBrakeBias        float64 `ibt:"dcBrakeBias"`
	DcDiffEntry        float64 `ibt:"dcDiffEntry"`
	DcDiffExit         float64 `ibt:"dcDiffExit"`
	DcDiffMiddle       float64 `ibt:"dcDiffMiddle"`
	DcEngineBraking    float64 `ibt:"dcEngineBraking"`
	DcEnginePower      float64 `ibt:"dcEnginePower"`
	DcFuelMixture      float64 `ibt:"dcFuelMixture"`
	DcRevLimiter       float64 `ibt:"dcRevLimiter"`
	DcThrottleShape    float64 `ibt:"dcThrottleShape"`
	DcTractionControl  float64 `ibt:"dcTractionControl"`
	DcTractionControl2 float64 `ibt:"dcTractionControl2"`
	DcTractionControlToggle bool `ibt:"dcTractionControlToggle"`
	DcWeightJackerLeft  float64 `ibt:"dcWeightJackerLeft"`
	DcWeightJackerRight float64 `ibt:"dcWeightJackerRight"`
	DcWingFront        float64 `ibt:"dcWingFront"`
	DcWingRear         float64 `ibt:"dcWingRear"`

	SessionID   string
	SessionType string
	SessionName string
	TrackName   string
	TrackID     int
	WorkerID    int
	GroupNum    int

	TickTime time.Time
}

// ToMap converts a TelemetryTick struct to a Tick map using reflection.
// This allows struct-based parsers to work with legacy map-based processors.
// Only fields with `ibt` tags and non-zero values are included.
func (t *TelemetryTick) ToMap(whitelist []string) Tick {
	tick := make(Tick, len(whitelist))

	// Create a whitelist lookup map for O(1) checks
	wantFields := make(map[string]bool, len(whitelist))
	for _, field := range whitelist {
		wantFields[field] = true
	}

	// Map iRacing field names to struct values
	// This is more efficient than reflection and matches the actual data layout
	fieldMap := map[string]interface{}{
		"Lap":                             t.LapID,
		"LapDistPct":                      float32(t.LapDistPct),
		"LapDist":                         float32(t.LapDist),
		"PlayerCarPosition":               float32(t.PlayerCarPosition),
		"PlayerCarClassPosition":          t.PlayerCarClassPosition,
		"PlayerCarIdx":                    t.PlayerCarIdx,
		"LapCurrentLapTime":               float32(t.LapCurrentLapTime),
		"LapLastLapTime":                  float32(t.LapLastLapTime),
		"LapLastNLapTime":                 float32(t.LapLastNLapTime),
		"LapBestLap":                      t.LapBestLap,
		"LapBestLapTime":                  float32(t.LapBestLapTime),
		"LapBestNLapLap":                  t.LapBestNLapLap,
		"LapBestNLapTime":                 float32(t.LapBestNLapTime),
		"LapLasNLapSeq":                   t.LapLasNLapSeq,
		"LapDeltaToBestLap":               float32(t.LapDeltaToBestLap),
		"LapDeltaToBestLap_DD":            float32(t.LapDeltaToBestLap_DD),
		"LapDeltaToBestLap_OK":            t.LapDeltaToBestLap_OK,
		"LapDeltaToOptimalLap":            float32(t.LapDeltaToOptimalLap),
		"LapDeltaToOptimalLap_DD":         float32(t.LapDeltaToOptimalLap_DD),
		"LapDeltaToOptimalLap_OK":         t.LapDeltaToOptimalLap_OK,
		"LapDeltaToSessionBestLap":        float32(t.LapDeltaToSessionBestLap),
		"LapDeltaToSessionBestLap_DD":     float32(t.LapDeltaToSessionBestLap_DD),
		"LapDeltaToSessionBestLap_OK":     t.LapDeltaToSessionBestLap_OK,
		"LapDeltaToSessionLastlLap":       float32(t.LapDeltaToSessionLastlLap),
		"LapDeltaToSessionLastlLap_DD":    float32(t.LapDeltaToSessionLastlLap_DD),
		"LapDeltaToSessionLastlLap_OK":    t.LapDeltaToSessionLastlLap_OK,
		"LapDeltaToSessionOptimalLap":     float32(t.LapDeltaToSessionOptimalLap),
		"LapDeltaToSessionOptimalLap_DD":  float32(t.LapDeltaToSessionOptimalLap_DD),
		"LapDeltaToSessionOptimalLap_OK":  t.LapDeltaToSessionOptimalLap_OK,
		"SessionNum":                      t.SessionNum,
		"SessionTime":                     t.SessionTime,
		"SessionTimeRemain":               t.SessionTimeRemain,
		"SessionLapsRemain":               t.SessionLapsRemain,
		"SessionState":                    t.SessionState,
		"SessionUniqueID":                 t.SessionUniqueID,
		"Throttle":                        float32(t.Throttle),
		"ThrottleRaw":                     float32(t.ThrottleRaw),
		"Brake":                           float32(t.Brake),
		"BrakeRaw":                        float32(t.BrakeRaw),
		"Clutch":                          float32(t.Clutch),
		"SteeringWheelAngle":              float32(t.SteeringWheelAngle),
		"SteeringWheelAngleMax":           float32(t.SteeringWheelAngleMax),
		"SteeringWheelTorque":             float32(t.SteeringWheelTorque),
		"SteeringWheelPctTorque":          float32(t.SteeringWheelPctTorque),
		"SteeringWheelPctTorqueSign":      float32(t.SteeringWheelPctTorqueSign),
		"SteeringWheelPctTorqueSignStops": float32(t.SteeringWheelPctTorqueSignStops),
		"SteeringWheelPctDamper":          float32(t.SteeringWheelPctDamper),
		"Speed":                           float32(t.Speed),
		"VelocityX":                       float32(t.VelocityX),
		"VelocityY":                       float32(t.VelocityY),
		"VelocityZ":                       float32(t.VelocityZ),
		"LatAccel":                        float32(t.LatAccel),
		"LongAccel":                       float32(t.LongAccel),
		"VertAccel":                       float32(t.VertAccel),
		"pitch":                           float32(t.Pitch),
		"PitchRate":                       float32(t.PitchRate),
		"roll":                            float32(t.Roll),
		"RollRate":                        float32(t.RollRate),
		"yaw":                             float32(t.Yaw),
		"YawNorth":                        float32(t.YawNorth),
		"YawRate":                         float32(t.YawRate),
		"Lat":                             t.Lat,
		"Lon":                             t.Lon,
		"alt":                             float32(t.Alt),
		"RPM":                             float32(t.RPM),
		"ShiftGrindRPM":                   float32(t.ShiftGrindRPM),
		"ShiftIndicatorPct":               float32(t.ShiftIndicatorPct),
		"ShiftPowerPct":                   float32(t.ShiftPowerPct),
		"Gear":                            t.Gear,
		"Voltage":                         float32(t.Voltage),
		"WaterTemp":                       float32(t.WaterTemp),
		"WaterLevel":                      float32(t.WaterLevel),
		"OilTemp":                         float32(t.OilTemp),
		"OilPress":                        float32(t.OilPress),
		"OilLevel":                        float32(t.OilLevel),
		"FuelPress":                       float32(t.FuelPress),
		"ManifoldPress":                   float32(t.ManifoldPress),
		"FuelLevel":                       float32(t.FuelLevel),
		"FuelLevelPct":                    float32(t.FuelLevelPct),
		"FuelUsePerHour":                  float32(t.FuelUsePerHour),
		"AirDensity":                      float32(t.AirDensity),
		"AirPressure":                     float32(t.AirPressure),
		"AirTemp":                         float32(t.AirTemp),
		"FogLevel":                        float32(t.FogLevel),
		"Skies":                           t.Skies,
		"WeatherType":                     t.WeatherType,
		"RelativeHumidity":                float32(t.RelativeHumidity),
		"TrackTemp":                       float32(t.TrackTemp),
		"TrackTempCrew":                   float32(t.TrackTempCrew),
		"WindDir":                         float32(t.WindDir),
		"WindVel":                         float32(t.WindVel),
		"EnterExitReset":                  t.EnterExitReset,
		"IsOnTrack":                       t.IsOnTrack,
		"IsOnTrackCar":                    t.IsOnTrackCar,
		"OnPitRoad":                       t.OnPitRoad,
		"CpuUsageBG":                      float32(t.CpuUsageBG),
		"FrameRate":                       float32(t.FrameRate),
		"PitRepairLeft":                   float32(t.PitRepairLeft),
		"PitOptRepairLeft":                float32(t.PitOptRepairLeft),
		"PitSvFuel":                       float32(t.PitSvFuel),
		"PitSvLFP":                        float32(t.PitSvLFP),
		"PitSvRFP":                        float32(t.PitSvRFP),
		"PitSvLRP":                        float32(t.PitSvLRP),
		"PitSvRRP":                        float32(t.PitSvRRP),
		"DriverMarker":                    t.DriverMarker,
		"LFpressure":                      float32(t.LFpressure),
		"RFpressure":                      float32(t.RFpressure),
		"LRpressure":                      float32(t.LRpressure),
		"RRpressure":                      float32(t.RRpressure),
		"LFtempM":                         float32(t.LFtempM),
		"RFtempM":                         float32(t.RFtempM),
		"LRtempM":                         float32(t.LRtempM),
		"RRtempM":                         float32(t.RRtempM),
		"LFtempCL":                        float32(t.LFtempCL),
		"LFtempCM":                        float32(t.LFtempCM),
		"LFtempCR":                        float32(t.LFtempCR),
		"RFtempCL":                        float32(t.RFtempCL),
		"RFtempCM":                        float32(t.RFtempCM),
		"RFtempCR":                        float32(t.RFtempCR),
		"LRtempCL":                        float32(t.LRtempCL),
		"LRtempCM":                        float32(t.LRtempCM),
		"LRtempCR":                        float32(t.LRtempCR),
		"RRtempCL":                        float32(t.RRtempCL),
		"RRtempCM":                        float32(t.RRtempCM),
		"RRtempCR":                        float32(t.RRtempCR),
		"LFwearL":                         float32(t.LFwearL),
		"LFwearM":                         float32(t.LFwearM),
		"LFwearR":                         float32(t.LFwearR),
		"RFwearL":                         float32(t.RFwearL),
		"RFwearM":                         float32(t.RFwearM),
		"RFwearR":                         float32(t.RFwearR),
		"LRwearL":                         float32(t.LRwearL),
		"LRwearM":                         float32(t.LRwearM),
		"LRwearR":                         float32(t.LRwearR),
		"RRwearL":                         float32(t.RRwearL),
		"RRwearM":                         float32(t.RRwearM),
		"RRwearR":                         float32(t.RRwearR),
		"LFcoldPressure":                  float32(t.LFcoldPressure),
		"RFcoldPressure":                  float32(t.RFcoldPressure),
		"LRcoldPressure":                  float32(t.LRcoldPressure),
		"RRcoldPressure":                  float32(t.RRcoldPressure),
		"LFshockDefl":                     float32(t.LFshockDefl),
		"RFshockDefl":                     float32(t.RFshockDefl),
		"LRshockDefl":                     float32(t.LRshockDefl),
		"RRshockDefl":                     float32(t.RRshockDefl),
		"CFshockDefl":                     float32(t.CFshockDefl),
		"CRshockDefl":                     float32(t.CRshockDefl),
		"LFshockVel":                      float32(t.LFshockVel),
		"RFshockVel":                      float32(t.RFshockVel),
		"LRshockVel":                      float32(t.LRshockVel),
		"RRshockVel":                      float32(t.RRshockVel),
		"CFshockVel":                      float32(t.CFshockVel),
		"CRshockVel":                      float32(t.CRshockVel),
		"LFspeed":                         float32(t.LFspeed),
		"RFspeed":                         float32(t.RFspeed),
		"LRspeed":                         float32(t.LRspeed),
		"RRspeed":                         float32(t.RRspeed),
		"LFbrakeLinePress":                float32(t.LFbrakeLinePress),
		"RFbrakeLinePress":                float32(t.RFbrakeLinePress),
		"LRbrakeLinePress":                float32(t.LRbrakeLinePress),
		"RRbrakeLinePress":                float32(t.RRbrakeLinePress),
		"dcABS":                           float32(t.DcABS),
		"dcAntiRollFront":                 float32(t.DcAntiRollFront),
		"dcAntiRollRear":                  float32(t.DcAntiRollRear),
		"dcBoostLevel":                    float32(t.DcBoostLevel),
		"dcBrakeBias":                     float32(t.DcBrakeBias),
		"dcDiffEntry":                     float32(t.DcDiffEntry),
		"dcDiffExit":                      float32(t.DcDiffExit),
		"dcDiffMiddle":                    float32(t.DcDiffMiddle),
		"dcEngineBraking":                 float32(t.DcEngineBraking),
		"dcEnginePower":                   float32(t.DcEnginePower),
		"dcFuelMixture":                   float32(t.DcFuelMixture),
		"dcRevLimiter":                    float32(t.DcRevLimiter),
		"dcThrottleShape":                 float32(t.DcThrottleShape),
		"dcTractionControl":               float32(t.DcTractionControl),
		"dcTractionControl2":              float32(t.DcTractionControl2),
		"dcTractionControlToggle":         t.DcTractionControlToggle,
		"dcWeightJackerLeft":              float32(t.DcWeightJackerLeft),
		"dcWeightJackerRight":             float32(t.DcWeightJackerRight),
		"dcWingFront":                     float32(t.DcWingFront),
		"dcWingRear":                      float32(t.DcWingRear),
	}

	// Only include requested fields
	for field := range wantFields {
		if value, exists := fieldMap[field]; exists {
			tick[field] = value
		}
	}

	return tick
}

// BuildWhitelistFromStruct extracts field names from ibt tags in a struct.
// Returns a whitelist suitable for use with ToMap(), ToMapFromStruct(), or parser constructors.
//
// The template parameter can be a struct value, pointer to struct, or reflect.Type.
// Only fields with `ibt` tags are included in the whitelist.
//
// Example:
//
//	type MyFields struct {
//	    Speed float64 `ibt:"Speed"`
//	    Gear  uint32  `ibt:"Gear"`
//	    RPM   float64 `ibt:"RPM"`
//	}
//	whitelist := ibt.BuildWhitelistFromStruct(MyFields{})
//	// Returns: []string{"Speed", "Gear", "RPM"}
//
// This is useful for:
//  1. Building whitelists for parsers: NewDirectStructParser(reader, header, whitelist...)
//  2. Creating type-safe processor whitelists
//  3. Self-documenting field selection
func BuildWhitelistFromStruct(template interface{}) []string {
	var typ reflect.Type

	// Handle different input types
	switch v := template.(type) {
	case reflect.Type:
		typ = v
	default:
		typ = reflect.TypeOf(template)
	}

	// Dereference pointer if needed
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	// Validate it's a struct
	if typ.Kind() != reflect.Struct {
		return []string{}
	}

	whitelist := make([]string, 0, typ.NumField())

	// Iterate through struct fields and extract ibt tags
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)

		// Get the ibt tag value
		if tag, ok := field.Tag.Lookup("ibt"); ok && tag != "" {
			whitelist = append(whitelist, tag)
		}
	}

	return whitelist
}

// ToMapFromStruct converts TelemetryTick to a map using a struct template.
// The template struct's ibt tags define which fields to extract.
//
// This method uses reflection to extract field names from the template struct,
// then delegates to ToMap() for the actual conversion. For high-frequency conversions
// with the same template, consider using DirectStructParser.ToMapWithTemplate()
// which caches the reflection results.
//
// Example:
//
//	type MyTelemetry struct {
//	    Speed    float64 `ibt:"Speed"`
//	    Gear     uint32  `ibt:"Gear"`
//	    Throttle float64 `ibt:"Throttle"`
//	    Brake    float64 `ibt:"Brake"`
//	}
//
//	tick, _ := parser.NextStruct()
//	data := tick.ToMapFromStruct(MyTelemetry{})
//	// Returns: map[string]interface{}{
//	//   "Speed": 123.45,
//	//   "Gear": uint32(4),
//	//   "Throttle": 0.85,
//	//   "Brake": 0.0,
//	// }
//
// Performance:
//   - Reflection overhead: ~500ns per call
//   - Map building: ~200ns per field
//   - Total: ~200ns + (500ns * num_fields) per call
//
// For repeated conversions with the same template, use DirectStructParser.ToMapWithTemplate()
// which caches reflection results and reduces overhead to ~50ns.
func (t *TelemetryTick) ToMapFromStruct(template interface{}) Tick {
	whitelist := BuildWhitelistFromStruct(template)
	return t.ToMap(whitelist)
}
