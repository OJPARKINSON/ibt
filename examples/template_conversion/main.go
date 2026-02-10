package main

import (
	"log"

	"github.com/OJPARKINSON/ibt"
)

// Define a custom struct with only the fields you want.
// The ibt tags specify which telemetry fields to extract.
type DashboardData struct {
	Speed    float64 `ibt:"Speed"`
	RPM      float64 `ibt:"RPM"`
	Gear     uint32  `ibt:"Gear"`
	Throttle float64 `ibt:"Throttle"`
	Brake    float64 `ibt:"Brake"`
	Lap      int32   `ibt:"Lap"`
}

func main() {
	// Build a whitelist from struct tags (type-safe!)
	log.Println("=== BuildWhitelistFromStruct ===")
	whitelist := ibt.BuildWhitelistFromStruct(DashboardData{})
	log.Printf("Whitelist: %v\n\n", whitelist)

	// Access fields directly on TelemetryTick (no map needed)
	log.Println("=== Direct Struct Field Access ===")
	tick := &ibt.TelemetryTick{
		Speed:    123.45,
		RPM:      5500.0,
		Gear:     4,
		Throttle: 0.85,
		Brake:    0.0,
		LapID:    5,
	}
	log.Printf("Speed: %.2f, RPM: %.0f, Gear: %d, Throttle: %.2f\n",
		tick.Speed, tick.RPM, tick.Gear, tick.Throttle)

	// Processor API with auto-whitelist
	log.Println("\n=== Processor API ===")
	log.Println("Implement Fields() to auto-extract whitelist:")
	log.Println("")
	log.Println("  func (p *MyProcessor) Fields() interface{} {")
	log.Println("      return struct {")
	log.Println("          Speed float64 `ibt:\"Speed\"`")
	log.Println("          RPM   float64 `ibt:\"RPM\"`")
	log.Println("      }{}")
	log.Println("  }")
}
