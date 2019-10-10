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
// As I separate out into separate lambda functions I may collapse the DynamoDB code into the same files as the business logic in some cases.
func UpdateAverage(newDuration int64) {
	updateStats(newDuration)
}

// GetStats returns a JSON encoded string representing the statistics
func GetStats() (string, error) {
	totalStats, err := getStats()
	stats := new(statistics)
	stats.NumAveraged = totalStats.totalCount
	stats.Average = totalStats.totalDuration / totalStats.totalCount
	b, err := json.Marshal(stats)
	if err != nil {
		msg := fmt.Sprintf("Received error when marshalling JSON: %s", err)
		log.Printf(msg)
		return "", errors.New(msg)
	}
	return string(b), nil
}
