package main

import (
	"fmt"

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
	fmt.Println("=== BuildWhitelistFromStruct ===")
	whitelist := ibt.BuildWhitelistFromStruct(DashboardData{})
	fmt.Printf("Whitelist: %v\n\n", whitelist)

	// Access fields directly on TelemetryTick (no map needed)
	fmt.Println("=== Direct Struct Field Access ===")
	tick := &ibt.TelemetryTick{
		Speed:    123.45,
		RPM:      5500.0,
		Gear:     4,
		Throttle: 0.85,
		Brake:    0.0,
		LapID:    5,
	}
	fmt.Printf("Speed: %.2f, RPM: %.0f, Gear: %d, Throttle: %.2f\n",
		tick.Speed, tick.RPM, tick.Gear, tick.Throttle)

	// Processor API with auto-whitelist
	fmt.Println("\n=== Processor API ===")
	fmt.Println("Implement Fields() to auto-extract whitelist:")
	fmt.Println("")
	fmt.Println("  func (p *MyProcessor) Fields() interface{} {")
	fmt.Println("      return struct {")
	fmt.Println("          Speed float64 `ibt:\"Speed\"`")
	fmt.Println("          RPM   float64 `ibt:\"RPM\"`")
	fmt.Println("      }{}")
	fmt.Println("  }")
}
