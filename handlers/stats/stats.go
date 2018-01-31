package statsHandler

import (
	"log"
	"sync"
	"time"
)

// Keys for the map that stores request time and total requests
const (
	tr = "total_requests"
	tt = "total_time"
)

// statsSync to be able to handle concurrency using sync.RWMutex
var statsSync = struct {
	sync.RWMutex
	dataMap map[string]int
}{dataMap: make(map[string]int)}

// Stats of the Password Hash API, total requests and total response time
type Stats struct {
	TotalRequests       int `json:"total"`
	AverageResponseTime int `json:"average"`
}

// UpdateStats updates the stat map with the total time and increments the total # of request count. Duration is in nanoseconds.
func UpdateStats(duration time.Duration) {
	v := duration.Nanoseconds() / int64(time.Millisecond)
	log.Println("total time in (ms)", v)

	statsSync.Lock() // lock the stats map before writing as we need to handle concurrent requests
	totalRequests, _ := statsSync.dataMap[tr]
	statsSync.dataMap[tr] = totalRequests + 1

	totalTime, _ := statsSync.dataMap[tt]
	statsSync.dataMap[tt] = totalTime + int(v)

	statsSync.Unlock()
}

// GetStats return the stat data with total # of requests to /hash and /hash/{id} endpoint and average response time by those endpoints
// The stats capture only successfuly API requests
func GetStats() *Stats {

	var averageTime = 0
	statsSync.RLock() // lock map before reading as we need to handle concurrent requests
	totalRequests, ok := statsSync.dataMap[tr]
	if !ok {
		totalRequests = 0 // default request #
	}

	totalTime, ok := statsSync.dataMap[tt]
	log.Printf("total time %v, total request %v \n", totalTime, totalRequests)
	if ok && totalRequests > 0 { // calculate the average time based on total time and total # of requests
		averageTime = totalTime / totalRequests
	}
	statsSync.RUnlock() // unlock the map once done

	resp := &Stats{
		TotalRequests:       totalRequests,
		AverageResponseTime: averageTime,
	}
	return resp
}

// TrackHashPasswordAPIMetrics should be called using defer and it needs to run at the end capturing the total time and updating the stats map
func TrackHashPasswordAPIMetrics(start time.Time) {
	trackFn := func() {
		elasped := time.Now().Sub(start) // calculate the total time elaspsed and call the stats handler
		UpdateStats(elasped)
	}
	trackFn()
}
