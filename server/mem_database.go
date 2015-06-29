package main

import (
	"log"
	"time"
)

// Simple in-memory Database implementation, just keeping all data in a slice
type MemoryDatabase struct {
	waterLevels []*WaterLevel
}

// NewMemoryDatabase creates a new MemoryDatabase
func NewMemoryDatabase() *MemoryDatabase {
	return &MemoryDatabase{}
}

// RecordWaterLevel records an individual reading of the sump water level
func (db *MemoryDatabase) RecordWaterLevel(level WaterLevel) error {
	for _, storedLevel := range db.waterLevels {
		if storedLevel.Time == level.Time {
			return nil
		}
	}
	log.Printf("Recording water level: %f", level.Level)
	db.waterLevels = append(db.waterLevels, &level)
	return nil
}

// FetchWaterLevelHistory returns a list of recent water level history
func (db *MemoryDatabase) FetchWaterLevelHistory(timeSpan time.Duration) ([]*WaterLevel, error) {

	earliestTime := time.Now()
	earliestTime = earliestTime.Add(time.Duration(-1) * timeSpan)

	for i := range db.waterLevels {
		if !db.waterLevels[i].Time.Before(earliestTime) {
			return db.waterLevels[i:], nil
		}
	}

	return db.waterLevels, nil
}
