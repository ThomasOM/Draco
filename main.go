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
	wordOptions, err := content.ParseWordOption(os.Args[1])
	if err != nil {
		panic(err)
	}

	count, err:= strconv.Atoi(os.Args[2])
	if err != nil {
		panic(err)
	}

	// Create all resources
	content := content.FetchAndCreateContent(wordOptions, count)
	stats := stats.CreateStats()
	display := display.CreateDisplay(content, stats)

	// Locks the main thread until completion
	display.Start()
}
