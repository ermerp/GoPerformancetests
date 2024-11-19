package main

import "time"

const (
	MAX_RETRIES = 100
)

// calculateRetryDelay berechnet die Verzögerung basierend auf der Anzahl der Versuche
func calculateRetryDelay(attempt int) time.Duration {
	return time.Duration(attempt*1000) * time.Millisecond
}
