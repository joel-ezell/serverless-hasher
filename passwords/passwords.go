package passwords

import (
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/joel-ezell/gohasher/statistics"
)

type passwords struct {
	// Stores the number of requested indices
	count   int
	passMap map[int]string
	mu      sync.RWMutex
}

const delaySecs = 5

var instance *passwords
var hashWg sync.WaitGroup

var once sync.Once

func getInstance() *passwords {
	// This ensures that the singleton is instantiated only once, even if multiple initial requests arrive at the same time
	once.Do(func() {
		instance = &passwords{
			passMap: make(map[int]string),
			count:   0}
	})
	return instance
}

// HashAndStore Computes a SHA-512 hash of the specified password, encodes it in Base64, then stores the password in a map
func HashAndStore(pwd string) (int, error) {
	start := time.Now()
	index := nextIndex()
	fmt.Print("Starting worker\n")
	go hashWorker(index, pwd, start)
	return index, nil
}

func hashWorker(index int, pwd string, start time.Time) {
	fmt.Printf("Before sleep\n")
	time.Sleep(delaySecs * time.Second)
	fmt.Printf("After sleep\n")
	sha := sha512.New()
	sha.Write([]byte(pwd))
	encodedPwd := base64.StdEncoding.EncodeToString(sha.Sum(nil))
	err := putHash(index, encodedPwd)
	if err != nil {
		log.Printf("Error received when storing password in DynamoDB: %s", err)
	}
	duration := time.Since(start)
	statistics.UpdateAverage(duration.Nanoseconds() / 1000)
}

// GetHashedPassword Returns the hashed password at the specified index
func GetHashedPassword(index int) (string, error) {
	hashedPwd, err := getHash(index)

	return hashedPwd, err
}
