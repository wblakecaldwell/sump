package main

import (
	"time"
)

type Database interface {
	// RecordWaterLevel records an individual reading of the sump water level
	RecordWaterLevel(level WaterLevel) error

	// FetchWaterLevelHistory returns a list of recent water level history
	FetchWaterLevelHistory(timeSpan time.Duration) ([]*WaterLevel, error)
}
