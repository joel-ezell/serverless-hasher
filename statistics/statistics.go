package statistics

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sync"
)

type statistics struct {
	// Stores the number of items averaged
	NumAveraged int64 `json:"total"`
	// Stores the average rounded to the nearest integer
	Average int64 `json:"average"`
}

var (
	s    *statistics
	once sync.Once
)

func getInstance() *statistics {
	if s == nil {
		once.Do(func() {
			s = &statistics{
				NumAveraged: 0,
				Average:     0}
		})
	}
	return s
}

// UpdateAverage This thread-safe function calculates and stores a new average, incorporating the provided value.
// The updated average is returned.
func UpdateAverage(newDuration int64) {
	// Another way to do this which would have slightly higher accuracy would be to simply maintain a running total and number of averaged values.
	updateStats(newDuration)
	// newTotal := (stats.floatAverage*float64(stats.NumAveraged) + float64(newDuration))
	// stats.NumAveraged++
	// stats.floatAverage = newTotal / float64(stats.NumAveraged)
	// stats.Average = int64(math.Round(stats.floatAverage))
	// return stats.Average
}

// GetStats returns a JSON encoded string representing the statistics
func GetStats() (string, error) {
	totalStats, err := getStats()
	stats := new(statistics)
	stats.NumAveraged = totalStats.totalCount
	stats.Average = totalStats.totalDuration / totalStats.totalCount;
	b, err := json.Marshal(stats)
	if err != nil {
		msg := fmt.Sprintf("Received error when marshalling JSON: %s", err)
		log.Printf(msg)
		return "", errors.New(msg)
	}
	return string(b), nil
}
