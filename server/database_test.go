package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

// runSimpleScenario runs through a basic scenario, recording water levels and pump runs
func runSimpleScenario(t *testing.T, db Database) {
	// initial state
	levels, err := db.FetchWaterLevelHistory(10 * time.Hour)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(levels))

	// record two water levels, including a duplicate
	now := time.Now()
	later := now.Add(5 * time.Second)
	err = db.RecordWaterLevel(WaterLevel{Time: now, Level: 1.5})
	assert.NoError(t, err)
	err = db.RecordWaterLevel(WaterLevel{Time: now, Level: 1.5})
	assert.NoError(t, err)
	err = db.RecordWaterLevel(WaterLevel{Time: later, Level: 1.6})
	assert.NoError(t, err)

	// check them
	levels, err = db.FetchWaterLevelHistory(10 * time.Hour)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(levels))
	assert.Equal(t, 1.5, levels[0].Level)
	assert.Equal(t, now, levels[0].Time)
	assert.Equal(t, 1.6, levels[1].Level)
	assert.Equal(t, later, levels[1].Time)
}
