package main

import (
	"testing"
)

func TestMemoryWaterLevels(t *testing.T) {
	db := &MemoryDatabase{}
	runSimpleScenario(t, db)
}
