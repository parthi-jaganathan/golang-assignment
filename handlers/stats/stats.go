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

// statsSyncMap to be able to handle concurrency
var statsSyncMap = struct {
	sync.RWMutex
	m map[string]int
}{m: make(map[string]int)}

// Stats of the Password Hash API, total requests and total response time
type Stats struct {
	TotalRequests       int `json:"total"`
	AverageResponseTime int `json:"average"`
}

// UpdateStats updates the status map with the total time and increments the total requests - input value in nano seconds
func UpdateStats(duration time.Duration) {

	v := duration.Nanoseconds() / int64(time.Millisecond)
	log.Println("total time in (ms)", v)

	statsSyncMap.Lock()
	totalRequests, _ := statsSyncMap.m[tr]
	statsSyncMap.m[tr] = totalRequests + 1

	totalTime, _ := statsSyncMap.m[tt]
	statsSyncMap.m[tt] = totalTime + int(v)

	statsSyncMap.Unlock()
}

// GetStats ... TODO
func GetStats() *Stats {

	var averageTime = 0
	statsSyncMap.RLock()
	totalRequests, ok := statsSyncMap.m[tr]
	if !ok {
		totalRequests = 0
	}

	totalTime, ok := statsSyncMap.m[tt]
	log.Printf("total time %v, total request %v \n", totalTime, totalRequests)
	if ok && totalRequests > 0 {
		averageTime = totalTime / totalRequests
	}
	statsSyncMap.RUnlock()

	resp := &Stats{
		TotalRequests:       totalRequests,
		AverageResponseTime: averageTime,
	}

	return resp
}
