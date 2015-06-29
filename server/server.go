package main

import (
	"time"
)

// WaterLevel represents a reading from the sump pit client
type WaterLevel struct {
	Time  time.Time // time of reading, local time
	Level float32   // water level from bottom of pit in centimeters
}
