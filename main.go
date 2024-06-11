package main

import (
	"me/thomazz/draco/content"
	"me/thomazz/draco/display"
	"me/thomazz/draco/stats"
	"os"
	"strconv"
)

func main() {
	// Verify input arguments
	count, err:= strconv.Atoi(os.Args[1])
	if err != nil {
		panic(err)
	}

	// Create all resources
	content := content.FetchAndCreateContent(count)
	stats := stats.CreateStats()
	display := display.CreateDisplay(content, stats)

	// Locks the main thread until completion
	display.Start()
}
