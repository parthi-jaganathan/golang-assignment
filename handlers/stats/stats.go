package statsHandler

import (
	"log"
	"sync"
	"time"
)

// Keys for the map
const (
	tr = "total_requests"
	tt = "total_time"
)

// statsSync to be able to handle concurrency
var statsSync = struct {
	sync.RWMutex
	dataMap map[string]int
}{dataMap: make(map[string]int)}

// Stats of the Password Hash API, total requests and total response time
type Stats struct {
	TotalRequests       int `json:"total"`
	AverageResponseTime int `json:"average"`
}

// UpdateStats updates the status map with the total time and increments the total requests - input value in nano seconds
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

// GetStats return the Stats data with total requests to /hash and /hash/{id} endpoint and average response time
func GetStats() *Stats {

	var averageTime = 0
	statsSync.RLock() // lock map before reading
	totalRequests, ok := statsSync.dataMap[tr]
	if !ok {
		totalRequests = 0
	}

	totalTime, ok := statsSync.dataMap[tt]
	log.Printf("total time %v, total request %v \n", totalTime, totalRequests)
	if ok && totalRequests > 0 {
		averageTime = totalTime / totalRequests
	}
	statsSync.RUnlock()

	resp := &Stats{
		TotalRequests:       totalRequests,
		AverageResponseTime: averageTime,
	}
	return resp
}

// TrackHashPasswordAPIMetrics should be called using defer and it needs to run at the end capturing the total time and updating the stats map
func TrackHashPasswordAPIMetrics(start time.Time) {
	fn := func() {
		elasped := time.Now().Sub(start) // calculate the total time elaspsed and call the stats handler
		UpdateStats(elasped)
	}
	fn()
}
