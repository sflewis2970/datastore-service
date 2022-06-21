package common

import (
	"log"
	"os"
	"time"
)

// Build formatted time string
func GetFormattedTime(timeNow time.Time, timeFormat string) string {
	return timeNow.Format(timeFormat)
}

// Get working directory
func GetWorkingDir() (string, error) {
	workingDir, getErr := os.Getwd()
	if getErr != nil {
		log.Fatal("Error getting working directory...")
		return "", getErr
	}

	return workingDir, nil
}
