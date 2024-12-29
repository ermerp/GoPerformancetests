package main

import (
	"math/rand"
	"time"
)

const (
	MAX_RETRIES    = 100
	RETRY_DELAY_MS = 1000
)

var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

// calculates the delay for the retry
func calculateRetryDelay(attempt int) time.Duration {
	return time.Duration(attempt*RETRY_DELAY_MS)*time.Millisecond + time.Duration(rng.Intn(501))*time.Millisecond
}
