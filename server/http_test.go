// TODO: Need to test only retrieving level info from within the input time

package main

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

// TestWaterLevelsEndpointSuccess tests posting new water levels to the server - success path
func TestWaterLevelsEndpointSuccess(t *testing.T) {
	clearConfig()
	defer clearConfig()

	rwMutex := sync.RWMutex{}
	db := &MemoryDatabase{}
	endpoint := buildSumpRegisterLevelsHandler(db, 2.0, nil, "abcdefg", &rwMutex)

	// post 3 water levels
	req1, err := http.NewRequest("POST", "http://foo.bar/levels", bytes.NewBuffer([]byte("{\"secret\":\"abcdefg\", \"level\":1.5}")))
	assert.NoError(t, err)

	req2, err := http.NewRequest("POST", "http://foo.bar/levels", bytes.NewBuffer([]byte("{\"secret\":\"abcdefg\", \"level\":1.6}")))
	assert.NoError(t, err)

	req3, err := http.NewRequest("POST", "http://foo.bar/levels", bytes.NewBuffer([]byte("{\"secret\":\"abcdefg\", \"level\":1.7}")))
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	endpoint(w, req1)
	time.Sleep(100 * time.Millisecond)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, `"ok"`, string(w.Body.Bytes()))

	w = httptest.NewRecorder()
	endpoint(w, req2)
	time.Sleep(100 * time.Millisecond)

	w = httptest.NewRecorder()
	endpoint(w, req3)
	time.Sleep(100 * time.Millisecond)

	// verify two records made it into the database
	levels, err := db.FetchWaterLevelHistory(10 * time.Hour)
	assert.NoError(t, err)
	assert.Equal(t, 3, len(levels))
	assert.Equal(t, float32(1.5), levels[0].Level)
	assert.Equal(t, float32(1.6), levels[1].Level)
	assert.Equal(t, float32(1.7), levels[2].Level)
}

// TestWaterLevelsEndpointFailure tests posting new water levels to the server - failures
func TestWaterLevelsEndpointFailure(t *testing.T) {
	clearConfig()
	defer clearConfig()

	db := &MemoryDatabase{}
	rwMutex := sync.RWMutex{}
	endpoint := buildSumpRegisterLevelsHandler(db, 2.0, nil, "abcdefg", &rwMutex)

	// POST bad JSON
	var jsonStr = []byte("bad JSON")
	req, err := http.NewRequest("POST", "http://foo.bar/levels", bytes.NewBuffer(jsonStr))
	assert.NoError(t, err)
	w := httptest.NewRecorder()
	endpoint(w, req)
	assert.Equal(t, 400, w.Code)
	assert.Equal(t, `"Bad Request"`+"\n", string(w.Body.Bytes()))

	// GET request
	req, err = http.NewRequest("GET", "http://foo.bar/levels", nil)
	assert.NoError(t, err)
	w = httptest.NewRecorder()
	endpoint(w, req)
	assert.Equal(t, 400, w.Code)
	assert.Equal(t, "\"HTTP POST requests only\"\n", string(w.Body.Bytes()))

	// verify nothing made it into the database
	levels, err := db.FetchWaterLevelHistory(10 * time.Hour)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(levels))
}

func TestSumpInfoHandlerSuccess(t *testing.T) {
	clearConfig()
	defer clearConfig()

	db := NewMemoryDatabase()
	rwMutex := sync.RWMutex{}
	endpoint := buildSumpInfoHandler(db, 10*time.Hour, &rwMutex)

	var err error

	// create 2 water level entries and 2 pump runs
	now := time.Now()
	later := now.Add(5 * time.Second)
	err = db.RecordWaterLevel(WaterLevel{Time: now, Level: 1.5})
	assert.NoError(t, err)
	err = db.RecordWaterLevel(WaterLevel{Time: later, Level: 1.6})
	assert.NoError(t, err)

	// set the start time (global)
	originalStartTime := startTime
	defer func() { startTime = originalStartTime }()
	startTime = time.Date(2015, time.February, 19, 10, 9, 5, 10, time.Local)

	// get the server info
	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "http://foo.bar/levels", nil)
	endpoint(w, req)
	out, _ := ioutil.ReadAll(w.Body)
	pumpInfo := PumpInfo{}
	err = json.Unmarshal(out, &pumpInfo)
	assert.NoError(t, err)

	// verify water levels
	assert.Equal(t, 2, len(pumpInfo.WaterLevels))
	assert.Equal(t, now.Unix(), pumpInfo.WaterLevels[0].Time.Unix())
	assert.Equal(t, float32(1.5), pumpInfo.WaterLevels[0].Level)
	assert.Equal(t, later.Unix(), pumpInfo.WaterLevels[1].Time.Unix())
	assert.Equal(t, float32(1.6), pumpInfo.WaterLevels[1].Level)

	// verify start time, current time, uptime
	assert.Equal(t, startTime.Unix(), pumpInfo.StartTimeEpochSeconds)
	assert.True(t, time.Now().Unix()-pumpInfo.CurrentTimeEpochSeconds < 1)
	assert.True(t, time.Now().Unix()-pumpInfo.CurrentTimeEpochSeconds >= 0)
	assert.True(t, len(pumpInfo.Uptime) > 0)
}
