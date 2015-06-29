package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

type PumpInfo struct {
	WaterLevels             []*WaterLevel
	StartTimeEpochSeconds   int64
	CurrentTimeEpochSeconds int64
	Uptime                  string
}

// startTime is the server's start time
var startTime time.Time = time.Now()

// buildSumpInfoHandler builds an HTTP handler that returns the service and sump pit info
func buildSumpInfoHandler(db Database, durationToShow time.Duration, rwMutex *sync.RWMutex) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {

		// read-only endpoint
		rwMutex.RLock()
		defer rwMutex.RUnlock()

		now := time.Now()
		waterLevels, err := db.FetchWaterLevelHistory(durationToShow)
		if err != nil {
			http.Error(w, `"error"`, http.StatusInternalServerError)
			return
		}

		response := PumpInfo{
			WaterLevels:             waterLevels,
			StartTimeEpochSeconds:   startTime.Unix(),
			CurrentTimeEpochSeconds: now.Unix(),
			Uptime:                  now.Round(1 * time.Second).Sub(startTime.Round(1 * time.Second)).String(),
		}
		writeToResponse(w, response)
	}
}

// buildSumpRegisterLevelsHandler builds an HTTP handler that receives new sump pit readings
func buildSumpRegisterLevelsHandler(db Database, panicWaterLevel float64, emailer *SMTPEmailer, serverSecret string, rwMutex *sync.RWMutex) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		requestTime := time.Now() // capture the time before we get the lock, so the lock doesn't interere with the data

		// read/write endpoint
		rwMutex.Lock()
		defer rwMutex.Unlock()

		if strings.ToLower(req.Method) != "post" {
			http.Error(w, `"HTTP POST requests only"`, http.StatusBadRequest)
			return
		}

		// parse the client request
		type ClientRequest struct {
			Secret string  `json:"secret"`
			Level  float64 `json:"level"`
		}
		decoder := json.NewDecoder(req.Body)
		var data ClientRequest
		err := decoder.Decode(&data)
		if err != nil {
			http.Error(w, `"Bad Request"`, http.StatusBadRequest)
			return
		}

		// make sure client is using the right secret
		if data.Secret != serverSecret {
			http.Error(w, `"Forbidden"`, http.StatusForbidden)
		}

		if emailer != nil && panicWaterLevel > 0.0 && data.Level > panicWaterLevel {
			emailer.SendEmail("*** SUMP PIT PANIC ***",
				fmt.Sprintf("The sump pit water is now at %f inches!", data.Level))
		}

		waterLevel := WaterLevel{
			Time:  requestTime,
			Level: float32(data.Level),
		}
		db.RecordWaterLevel(waterLevel)
		writeToResponse(w, "ok")
	}
}

// indexHtmlHandler is an HTTP handler to return index.html
func indexHtmlHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(index_html_str))
}

// writeToResponse writes the input response object to JSON
func writeToResponse(w http.ResponseWriter, response interface{}) {
	// convert to JSON and write to the client
	js, err := json.Marshal(response)
	if err != nil {
		log.Printf("Error writing response: %s", err)
		http.Error(w, `"error"`, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
