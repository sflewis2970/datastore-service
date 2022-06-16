package common

import (
	"bytes"
	"net/http"
	"strings"
	"time"
)

// Build formatted time string
func GetFormattedTime(timeNow time.Time, timeFormat string) string {
	return timeNow.Format(timeFormat)
}

// Build UUID string
func BuildUUID(uuid string, delimiter string, nbrOfGroups int) string {
	newUUID := ""

	uuidList := strings.Split(uuid, delimiter)
	for key, value := range uuidList {
		if key < nbrOfGroups {
			newUUID = newUUID + value
		}
	}

	return newUUID
}

func GetRequest(destination string) (*http.Response, error) {
	response, getErr := http.Get(destination)
	if getErr != nil {
		return nil, getErr
	}

	return response, getErr
}

func PostRequest(destination string, jsonData []byte) (*http.Response, error) {
	response, postErr := http.Post(destination, "application/json", bytes.NewBuffer(jsonData))
	if postErr != nil {
		return nil, postErr
	}

	return response, postErr
}
