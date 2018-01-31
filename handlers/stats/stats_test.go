package statsHandler

import (
	"log"
	"testing"
	"time"
)

func TestUpdateMap(t *testing.T) {
	for i := 0; i < 10; i++ {
		d, e := time.ParseDuration("10000ms")
		if e != nil {
			t.Error(e)
		}
		UpdateStats(d)
		stats := GetStats()
		_s := *stats
		if _s.AverageResponseTime != 10000 {
			t.Errorf("Invalid Avg response time")
		}
		if _s.TotalRequests != i+1 {
			t.Errorf("Invalid Avg response time")
		}
	}
}

func TestGetEmptyStats(t *testing.T) {
	stats := GetStats()
	_s := *stats
	if _s.AverageResponseTime > 0 {
		t.Errorf("Invalid Avg response time")
	}
	if _s.TotalRequests > 0 {
		t.Errorf("Invalid Avg response time")
	}
}

func TestConcurrentUpdate(t *testing.T) {
	for i := 0; i < 1000; i++ {
		d, e := time.ParseDuration("10000ms")
		if e != nil {
			t.Error(e)
		}
		go UpdateStats(d)
	}
}

func TestConcurrentRead(t *testing.T) {
	for i := 0; i < 1000; i++ {
		d, e := time.ParseDuration("10000ms")
		if e != nil {
			t.Error(e)
		}
		UpdateStats(d)
	}

	for i := 0; i < 1000; i++ {
		go func() {
			stats := GetStats()
			_s := *stats
			log.Println(_s)
			if _s.AverageResponseTime != 10000 {
				t.Errorf("Invalid Avg response time")
			}
			if _s.TotalRequests == 0 {
				t.Errorf("Invalid Avg response time")
			}
		}()
	}
}
