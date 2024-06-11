package stats

import (
	"time"
)

type Stats struct {
	TimeStarted            time.Time
	TimeFinished           time.Time
	CorrectCharactersTyped int
	TotalCharactersTyped   int
	TotalSpacesTyped		int
	CorrectSpacesTyped		int
	TotalInputsTyped       int
	WordsTyped             int
}

func CreateStats() *Stats {
	return &Stats{}
}

func (stats *Stats) SecondsPassed() float64 {
	toTime := time.Now()
	if !stats.TimeFinished.IsZero() {
		toTime = stats.TimeFinished
	}

	secondsPassed := 0.0
	if !stats.TimeStarted.IsZero() {
		secondsPassed = toTime.Local().Sub(stats.TimeStarted).Seconds()
	}

	return secondsPassed
}

func (stats *Stats) Accuracy() float64 {
	if stats.TotalCharactersTyped == 0 {
		return 1.0
	}

	return float64(stats.CorrectCharactersTyped) / float64(stats.TotalCharactersTyped)
}

func (stats *Stats) WordsPerMinuteRaw() float64 {
	return wpm(stats.TotalCharactersTyped + stats.TotalSpacesTyped, stats.SecondsPassed())
}

func (stats *Stats) WordsPerMinute() float64 {
	return wpm(stats.CorrectCharactersTyped + stats.CorrectSpacesTyped, stats.SecondsPassed())
}

func wpm(typed int, seconds float64) float64 {
	if seconds == 0 {
		return 0.0
	}

	cps := float64(typed) / seconds
	wps := cps / 5.0 // Words are counted as 5 symbols on average
	return wps * 60.0
}
