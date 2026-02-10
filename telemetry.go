package ibt

import (
	"reflect"
	"time"
)

type TelemetryTick struct {
	// Lap & Position
	LapID                  int32   `ibt:"Lap"`
	LapDistPct             float64 `ibt:"LapDistPct"`
	LapDist                float64 `ibt:"LapDist"`
	PlayerCarPosition      float64 `ibt:"PlayerCarPosition"`
	PlayerCarClassPosition int32   `ibt:"PlayerCarClassPosition"`
	PlayerCarIdx           int32   `ibt:"PlayerCarIdx"`

	// Lap Timing
	LapCurrentLapTime float64 `ibt:"LapCurrentLapTime"`
	LapLastLapTime    float64 `ibt:"LapLastLapTime"`
	LapLastNLapTime   float64 `ibt:"LapLastNLapTime"`
	LapBestLap        int32   `ibt:"LapBestLap"`
	LapBestLapTime    float64 `ibt:"LapBestLapTime"`
	LapBestNLapLap    int32   `ibt:"LapBestNLapLap"`
	LapBestNLapTime   float64 `ibt:"LapBestNLapTime"`
	LapLasNLapSeq     int32   `ibt:"LapLasNLapSeq"`

	// Delta Timing
	LapDeltaToBestLap              float64 `ibt:"LapDeltaToBestLap"`
	LapDeltaToBestLap_DD           float64 `ibt:"LapDeltaToBestLap_DD"`
	LapDeltaToBestLap_OK           bool    `ibt:"LapDeltaToBestLap_OK"`
	LapDeltaToOptimalLap           float64 `ibt:"LapDeltaToOptimalLap"`
	LapDeltaToOptimalLap_DD        float64 `ibt:"LapDeltaToOptimalLap_DD"`
	LapDeltaToOptimalLap_OK        bool    `ibt:"LapDeltaToOptimalLap_OK"`
	LapDeltaToSessionBestLap       float64 `ibt:"LapDeltaToSessionBestLap"`
	LapDeltaToSessionBestLap_DD    float64 `ibt:"LapDeltaToSessionBestLap_DD"`
	LapDeltaToSessionBestLap_OK    bool    `ibt:"LapDeltaToSessionBestLap_OK"`
	LapDeltaToSessionLastlLap      float64 `ibt:"LapDeltaToSessionLastlLap"`
	LapDeltaToSessionLastlLap_DD   float64 `ibt:"LapDeltaToSessionLastlLap_DD"`
	LapDeltaToSessionLastlLap_OK   bool    `ibt:"LapDeltaToSessionLastlLap_OK"`
	LapDeltaToSessionOptimalLap    float64 `ibt:"LapDeltaToSessionOptimalLap"`
	LapDeltaToSessionOptimalLap_DD float64 `ibt:"LapDeltaToSessionOptimalLap_DD"`
	LapDeltaToSessionOptimalLap_OK bool    `ibt:"LapDeltaToSessionOptimalLap_OK"`

	// Session Info
	SessionNum        int32   `ibt:"SessionNum"`
	SessionTime       float64 `ibt:"SessionTime"`
	SessionTimeRemain float64 `ibt:"SessionTimeRemain"`
	SessionLapsRemain int32   `ibt:"SessionLapsRemain"`
	SessionState      int32   `ibt:"SessionState"`
	SessionUniqueID   int32   `ibt:"SessionUniqueID"`

	// Driver Inputs
	Throttle                        float64 `ibt:"Throttle"`
	ThrottleRaw                     float64 `ibt:"ThrottleRaw"`
	Brake                           float64 `ibt:"Brake"`
	BrakeRaw                        float64 `ibt:"BrakeRaw"`
	Clutch                          float64 `ibt:"Clutch"`
	Gear                            uint32  `ibt:"Gear"`
	SteeringWheelAngle              float64 `ibt:"SteeringWheelAngle"`
	SteeringWheelAngleMax           float64 `ibt:"SteeringWheelAngleMax"`
	SteeringWheelTorque             float64 `ibt:"SteeringWheelTorque"`
	SteeringWheelPctTorque          float64 `ibt:"SteeringWheelPctTorque"`
	SteeringWheelPctTorqueSign      float64 `ibt:"SteeringWheelPctTorqueSign"`
	SteeringWheelPctTorqueSignStops float64 `ibt:"SteeringWheelPctTorqueSignStops"`
	SteeringWheelPctDamper          float64 `ibt:"SteeringWheelPctDamper"`

	// Speed & Motion
	Speed     float64 `ibt:"Speed"`
	VelocityX float64 `ibt:"VelocityX"`
	VelocityY float64 `ibt:"VelocityY"`
	VelocityZ float64 `ibt:"VelocityZ"`

	// Acceleration
	LatAccel  float64 `ibt:"LatAccel"`
	LongAccel float64 `ibt:"LongAccel"`
	VertAccel float64 `ibt:"VertAccel"`

	// Orientation
	Pitch     float64 `ibt:"pitch"`
	PitchRate float64 `ibt:"PitchRate"`
	Roll      float64 `ibt:"roll"`
	RollRate  float64 `ibt:"RollRate"`
	Yaw       float64 `ibt:"yaw"`
	YawNorth  float64 `ibt:"YawNorth"`
	YawRate   float64 `ibt:"YawRate"`

	// GPS Position
	Lat float64 `ibt:"Lat"`
	Lon float64 `ibt:"Lon"`
	Alt float64 `ibt:"alt"`

	// Engine
	RPM               float64 `ibt:"RPM"`
	ShiftGrindRPM     float64 `ibt:"ShiftGrindRPM"`
	ShiftIndicatorPct float64 `ibt:"ShiftIndicatorPct"`
	ShiftPowerPct     float64 `ibt:"ShiftPowerPct"`
	Voltage           float64 `ibt:"Voltage"`
	WaterTemp         float64 `ibt:"WaterTemp"`
	WaterLevel        float64 `ibt:"WaterLevel"`
	OilTemp           float64 `ibt:"OilTemp"`
	OilPress          float64 `ibt:"OilPress"`
	OilLevel          float64 `ibt:"OilLevel"`
	FuelPress         float64 `ibt:"FuelPress"`
	ManifoldPress     float64 `ibt:"ManifoldPress"`

	// Fuel
	FuelLevel      float64 `ibt:"FuelLevel"`
	FuelLevelPct   float64 `ibt:"FuelLevelPct"`
	FuelUsePerHour float64 `ibt:"FuelUsePerHour"`

	// Environment
	AirDensity       float64 `ibt:"AirDensity"`
	AirPressure      float64 `ibt:"AirPressure"`
	AirTemp          float64 `ibt:"AirTemp"`
	FogLevel         float64 `ibt:"FogLevel"`
	RelativeHumidity float64 `ibt:"RelativeHumidity"`
	Skies            int32   `ibt:"Skies"`
	TrackTemp        float64 `ibt:"TrackTemp"`
	TrackTempCrew    float64 `ibt:"TrackTempCrew"`
	WeatherType      int32   `ibt:"WeatherType"`
	WindDir          float64 `ibt:"WindDir"`
	WindVel          float64 `ibt:"WindVel"`

	// Track Position
	IsOnTrack    bool `ibt:"IsOnTrack"`
	IsOnTrackCar bool `ibt:"IsOnTrackCar"`
	OnPitRoad    bool `ibt:"OnPitRoad"`

	// System
	CpuUsageBG     float64 `ibt:"CpuUsageBG"`
	FrameRate      float64 `ibt:"FrameRate"`
	DriverMarker   bool    `ibt:"DriverMarker"`
	EnterExitReset int32   `ibt:"EnterExitReset"`

	// Pit Service
	PitRepairLeft    float64 `ibt:"PitRepairLeft"`
	PitOptRepairLeft float64 `ibt:"PitOptRepairLeft"`
	PitSvFuel        float64 `ibt:"PitSvFuel"`
	PitSvLFP         float64 `ibt:"PitSvLFP"`
	PitSvRFP         float64 `ibt:"PitSvRFP"`
	PitSvLRP         float64 `ibt:"PitSvLRP"`
	PitSvRRP         float64 `ibt:"PitSvRRP"`

	// Tire Pressure
	LFpressure float64 `ibt:"LFpressure"`
	RFpressure float64 `ibt:"RFpressure"`
	LRpressure float64 `ibt:"LRpressure"`
	RRpressure float64 `ibt:"RRpressure"`

	// Tire Temperature (Middle/Surface)
	LFtempM float64 `ibt:"LFtempM"`
	RFtempM float64 `ibt:"RFtempM"`
	LRtempM float64 `ibt:"LRtempM"`
	RRtempM float64 `ibt:"RRtempM"`

	// Tire Temperature (Carcass - Left, Middle, Right)
	LFtempCL float64 `ibt:"LFtempCL"`
	LFtempCM float64 `ibt:"LFtempCM"`
	LFtempCR float64 `ibt:"LFtempCR"`
	RFtempCL float64 `ibt:"RFtempCL"`
	RFtempCM float64 `ibt:"RFtempCM"`
	RFtempCR float64 `ibt:"RFtempCR"`
	LRtempCL float64 `ibt:"LRtempCL"`
	LRtempCM float64 `ibt:"LRtempCM"`
	LRtempCR float64 `ibt:"LRtempCR"`
	RRtempCL float64 `ibt:"RRtempCL"`
	RRtempCM float64 `ibt:"RRtempCM"`
	RRtempCR float64 `ibt:"RRtempCR"`

	// Tire Wear
	LFwearL float64 `ibt:"LFwearL"`
	LFwearM float64 `ibt:"LFwearM"`
	LFwearR float64 `ibt:"LFwearR"`
	RFwearL float64 `ibt:"RFwearL"`
	RFwearM float64 `ibt:"RFwearM"`
	RFwearR float64 `ibt:"RFwearR"`
	LRwearL float64 `ibt:"LRwearL"`
	LRwearM float64 `ibt:"LRwearM"`
	LRwearR float64 `ibt:"LRwearR"`
	RRwearL float64 `ibt:"RRwearL"`
	RRwearM float64 `ibt:"RRwearM"`
	RRwearR float64 `ibt:"RRwearR"`

	// Tire Cold Pressure
	LFcoldPressure float64 `ibt:"LFcoldPressure"`
	RFcoldPressure float64 `ibt:"RFcoldPressure"`
	LRcoldPressure float64 `ibt:"LRcoldPressure"`
	RRcoldPressure float64 `ibt:"RRcoldPressure"`

	// Suspension - Shock Deflection
	LFshockDefl float64 `ibt:"LFshockDefl"`
	RFshockDefl float64 `ibt:"RFshockDefl"`
	LRshockDefl float64 `ibt:"LRshockDefl"`
	RRshockDefl float64 `ibt:"RRshockDefl"`
	CFshockDefl float64 `ibt:"CFshockDefl"`
	CRshockDefl float64 `ibt:"CRshockDefl"`

	// Suspension - Shock Velocity
	LFshockVel float64 `ibt:"LFshockVel"`
	RFshockVel float64 `ibt:"RFshockVel"`
	LRshockVel float64 `ibt:"LRshockVel"`
	RRshockVel float64 `ibt:"RRshockVel"`
	CFshockVel float64 `ibt:"CFshockVel"`
	CRshockVel float64 `ibt:"CRshockVel"`

	// Wheel Speed
	LFspeed float64 `ibt:"LFspeed"`
	RFspeed float64 `ibt:"RFspeed"`
	LRspeed float64 `ibt:"LRspeed"`
	RRspeed float64 `ibt:"RRspeed"`

	// Brake Line Pressure
	LFbrakeLinePress float64 `ibt:"LFbrakeLinePress"`
	RFbrakeLinePress float64 `ibt:"RFbrakeLinePress"`
	LRbrakeLinePress float64 `ibt:"LRbrakeLinePress"`
	RRbrakeLinePress float64 `ibt:"RRbrakeLinePress"`

	// Dynamic In-Car Adjustments
	DcABS                   float64 `ibt:"dcABS"`
	DcAntiRollFront         float64 `ibt:"dcAntiRollFront"`
	DcAntiRollRear          float64 `ibt:"dcAntiRollRear"`
	DcBoostLevel            float64 `ibt:"dcBoostLevel"`
	DcBrakeBias             float64 `ibt:"dcBrakeBias"`
	DcDiffEntry             float64 `ibt:"dcDiffEntry"`
	DcDiffExit              float64 `ibt:"dcDiffExit"`
	DcDiffMiddle            float64 `ibt:"dcDiffMiddle"`
	DcEngineBraking         float64 `ibt:"dcEngineBraking"`
	DcEnginePower           float64 `ibt:"dcEnginePower"`
	DcFuelMixture           float64 `ibt:"dcFuelMixture"`
	DcRevLimiter            float64 `ibt:"dcRevLimiter"`
	DcThrottleShape         float64 `ibt:"dcThrottleShape"`
	DcTractionControl       float64 `ibt:"dcTractionControl"`
	DcTractionControl2      float64 `ibt:"dcTractionControl2"`
	DcTractionControlToggle bool    `ibt:"dcTractionControlToggle"`
	DcWeightJackerLeft      float64 `ibt:"dcWeightJackerLeft"`
	DcWeightJackerRight     float64 `ibt:"dcWeightJackerRight"`
	DcWingFront             float64 `ibt:"dcWingFront"`
	DcWingRear              float64 `ibt:"dcWingRear"`

	SessionID   string
	SessionType string
	SessionName string
	TrackName   string
	TrackID     int
	WorkerID    int
	GroupNum    int

	TickTime time.Time
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
