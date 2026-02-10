package main

import (
	"context"
	"log"

	"github.com/OJPARKINSON/ibt"
	"github.com/OJPARKINSON/ibt/examples"
)

func main() {
	// Parse the files into stubs
	stubs, err := examples.ParseExampleStubs()
	if err != nil {
		log.Fatal(err)
	}

	// We group the stubs by iRacing session. This allows us to summarise results for
	// an entire session, instead of just a single ibt file.
	groups := stubs.Group()
	defer func() { _ = ibt.CloseAllStubs(groups) }()

	for groupIdx, group := range groups {
		// Create the instance(s) of your processor(s) for this group
		processor := newTrackTempProcessor()

		// Process the available telemetry using struct-based processing (50-60% faster)
		// This is currently only utilising the Track Temp processor,
		// but can include as many as you want.
		if err := ibt.Process(context.Background(), group, processor); err != nil {
			log.Printf("failed to process telemetry for group %d: %v", groupIdx, err)
			return
		}

		// Print the summarised track temperature
		processor.Print()
	}
}
