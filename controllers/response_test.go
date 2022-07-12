package controllers

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/sflewis2970/datastore-service/config"
	"github.com/sflewis2970/datastore-service/models/messages"
)

// The envionment variables will be set on the server
// For testing, set the variables manually
func setConfigEnv(activeDriver string) error {
	// Set hostname environment variable
	setErr := os.Setenv(config.HOSTNAME, "")
	if setErr != nil {
		log.Print("Error setting config vars...")
		return setErr
	}

	// Set hostport environment variable
	setErr = os.Setenv(config.HOSTPORT, ":9090")
	if setErr != nil {
		log.Print("Error setting config vars...")
		return setErr
	}

	// Set activedriver environment variable
	setErr = os.Setenv(config.ACTIVEDRIVER, activeDriver)
	if setErr != nil {
		log.Print("Error setting config vars...")
		return setErr
	}

	// Set Go-cache environment variable
	switch os.Getenv(config.ACTIVEDRIVER) {
	case "go-cache":
		setErr = os.Setenv(config.DEFAULT_EXPIRATION, "1")
		if setErr != nil {
			log.Print("Error setting config vars...")
			return setErr
		}

		setErr = os.Setenv(config.CLEANUP_INTERVAL, "30")
		if setErr != nil {
			log.Print("Error setting config vars...")
			return setErr
		}
	case "mysql":
		setErr = os.Setenv(config.MYSQL_CONNECTION, "root:devStation@tcp(127.0.0.1:3306)/")
		if setErr != nil {
			log.Print("Error setting config vars...")
			return setErr
		}
	case "postgres":
		setErr = os.Setenv(config.POSTGRES_HOST, "127.0.0.1")
		if setErr != nil {
			log.Print("Error setting config vars...")
			return setErr
		}

		setErr = os.Setenv(config.POSTGRES_PORT, "5432")
		if setErr != nil {
			log.Print("Error setting config vars...")
			return setErr
		}

		setErr = os.Setenv(config.POSTGRES_USER, "postgres")
		if setErr != nil {
			log.Print("Error setting config vars...")
			return setErr
		}
	}

	return nil
}

func InsertTest(t *testing.T, jsonData []byte) []byte {
	// Create new request
	request, reqErr := http.NewRequest("GET", "/api/v1/ds/insert", bytes.NewBuffer(jsonData))
	if reqErr != nil {
		t.Errorf("Could not create request.\n")
	}

	// Setup recoder
	rRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(Insert)
	handler.ServeHTTP(rRecorder, request)

	// Check response code
	status := rRecorder.Code
	if status != http.StatusOK {
		t.Errorf("handler returned invalid status code: got %d, expected: %d\n", status, http.StatusOK)
	}

	// Unmarshal JSON
	return rRecorder.Body.Bytes()
}

func GetTest(t *testing.T, jsonData []byte) []byte {
	// Create new request
	request, reqErr := http.NewRequest("GET", "/api/v1/ds/get", bytes.NewBuffer(jsonData))
	if reqErr != nil {
		t.Errorf("Could not create request.\n")
	}

	// Setup recoder
	rRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(Get)
	handler.ServeHTTP(rRecorder, request)

	// Check response code
	status := rRecorder.Code
	if status != http.StatusOK {
		t.Errorf("handler returned invalid status code: got %d, expected: %d\n", status, http.StatusOK)
	}

	// Unmarshal JSON
	return rRecorder.Body.Bytes()
}

func TestStatus(t *testing.T) {
	// Initialize logging
	log.SetFlags(log.Ldate | log.Lshortfile)

	// Set config environment variables
	setConfigEnv(config.GOCACHE_DRIVER)

	// Initialize controllers object
	New(config.REFRESH_CONFIG_DATA)

	// Create new request
	request, reqErr := http.NewRequest("GET", "/api/v1/ds/status", nil)
	if reqErr != nil {
		t.Errorf("Could not create request.\n")
	}

	// Setup response recorder
	respRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(Status)
	handler.ServeHTTP(respRecorder, request)

	// Check response code
	status := respRecorder.Code
	if status != http.StatusOK {
		t.Errorf("handler returned invalid status code: got %d, expected: %d\n", status, http.StatusOK)
	}

	// Body bytes array
	bodyBytes := respRecorder.Body.Bytes()

	// Unmarshal StatusResponse JSON
	var sResponse messages.StatusResponse
	unmarshalErr := json.Unmarshal(bodyBytes, &sResponse)
	if unmarshalErr != nil {
		t.Errorf(unmarshalErr.Error())
	}

	// Check status field
	if sResponse.Status != messages.StatusCode(messages.DS_RUNNING) {
		t.Errorf("Server Status returned is not running, expected running status: got %d", sResponse.Status)
	}
}

func TestInsert(t *testing.T) {
	// Initialize logging
	log.SetFlags(log.Ldate | log.Lshortfile)

	// Set config environment variables
	setConfigEnv(config.GOCACHE_DRIVER)

	// Initialize controllers object
	New(config.REFRESH_CONFIG_DATA)

	// Simulate a client sending a QuestionRequest to the datastore server
	// Question Request
	var qRequest messages.QuestionRequest

	// Build Question Request
	qRequest.QuestionID = "aaaabbbb"
	qRequest.Question = "According to Greek mythology, who was the first woman on earth?"
	qRequest.Category = "general"
	qRequest.Answer = "Pandora"

	// Marshal QuestionRequest
	jsonData, marshalErr := json.Marshal(qRequest)
	if marshalErr != nil {
		t.Errorf("New request error: %s", marshalErr.Error())
	}

	// Send AddQuestion request to datastore
	bodyBytes := InsertTest(t, jsonData)

	// Unmarshal data to QuestionResponse
	var qResponse messages.QuestionResponse
	unmarshalErr := json.Unmarshal(bodyBytes, &qResponse)
	if unmarshalErr != nil {
		t.Errorf(unmarshalErr.Error())
	}

	// Check Error field
	if len(qResponse.Error) > 0 {
		t.Errorf("An error occurred inserting record...")
	}
}

func TestGetBeforeInsert(t *testing.T) {
	// Initialize logging
	log.SetFlags(log.Ldate | log.Lshortfile)

	// Set config environment variables
	setConfigEnv(config.GOCACHE_DRIVER)

	// Initialize controllers object
	New(config.REFRESH_CONFIG_DATA)

	// Simulate a client sending a AnswerRequest to the datastore server
	// Build AnswerRequest
	var aRequest messages.AnswerRequest

	// Build Question Request
	aRequest.QuestionID = "aaaacccc"
	jsonData, marshalErr := json.Marshal(aRequest)
	if marshalErr != nil {
		t.Errorf("New request error: %s", marshalErr.Error())
	}

	// Send AddQuestion request to datastore
	bodyBytes := GetTest(t, jsonData)

	var aResponse messages.AnswerResponse
	unmarshalErr := json.Unmarshal(bodyBytes, &aResponse)
	if unmarshalErr != nil {
		t.Errorf(unmarshalErr.Error())
	}

	// Check Question field
	if len(aResponse.Question) > 0 {
		t.Errorf("The question field unexpectedly has been found!")
	}

	// Check Question field
	if aResponse.Message != messages.NO_RESULTS_RETURNED_MSG {
		t.Errorf("The message unexpectedly returned the wrong message, message returned: %s", aResponse.Message)
	}
}

func TestGetAfterInsert(t *testing.T) {
	// Initialize logging
	log.SetFlags(log.Ldate | log.Lshortfile)

	// Set config environment variables
	setConfigEnv(config.GOCACHE_DRIVER)

	// Initialize controllers object
	New(config.REFRESH_CONFIG_DATA)

	// Simulate a client sending a QuestionRequest to the datastore server
	// Question Request
	var qRequest messages.QuestionRequest

	// Build Question Request
	qRequest.QuestionID = "aaaadddd"
	qRequest.Question = "According to Greek mythology, who was the first woman on earth?"
	qRequest.Category = "general"
	qRequest.Answer = "Pandora"

	// Marshal QuestionRequest
	jsonData, marshalErr := json.Marshal(qRequest)
	if marshalErr != nil {
		t.Errorf("New request error: %s", marshalErr.Error())
	}

	// Send AddQuestion request to datastore
	bodyBytes := InsertTest(t, jsonData)

	// Unmarshal data to QuestionResponse
	var qResponse messages.QuestionResponse
	unmarshalErr := json.Unmarshal(bodyBytes, &qResponse)
	if unmarshalErr != nil {
		t.Errorf(unmarshalErr.Error())
	}

	// Check Error field
	if len(qResponse.Error) > 0 {
		t.Errorf("An error occurred inserting record...")
	}

	// Simulate a client sending a AnswerRequest to the datastore server
	// Build AnswerRequest
	var aRequest messages.AnswerRequest

	// Build Question Request
	aRequest.QuestionID = qRequest.QuestionID
	jsonData, marshalErr = json.Marshal(aRequest)
	if marshalErr != nil {
		t.Errorf("New request error: %s", marshalErr.Error())
	}

	// Send CheckAnswer request to datastore
	bodyBytes = GetTest(t, jsonData)

	// Build AnswerResponse
	var aResponse messages.AnswerResponse
	unmarshalErr = json.Unmarshal(bodyBytes, &aResponse)
	if unmarshalErr != nil {
		t.Errorf(unmarshalErr.Error())
	}

	// Check Error field
	if len(aResponse.Question) == 0 {
		t.Errorf("Unexpectedly, question was NOT returned...")
	}
}

func TestGetAfterDelete(t *testing.T) {
	// Initialize logging
	log.SetFlags(log.Ldate | log.Lshortfile)

	// Set config environment variables
	setConfigEnv(config.GOCACHE_DRIVER)

	// Initialize controllers object
	New(config.REFRESH_CONFIG_DATA)

	// Simulate a client sending a QuestionRequest to the datastore server
	// Question Request
	var qRequest messages.QuestionRequest

	// Build Question Request
	qRequest.QuestionID = "aaaaeeee"
	qRequest.Question = "According to Greek mythology, who was the first woman on earth?"
	qRequest.Category = "general"
	qRequest.Answer = "Pandora"

	// Marshal QuestionRequest
	jsonData, marshalErr := json.Marshal(qRequest)
	if marshalErr != nil {
		t.Errorf("New request error: %s", marshalErr.Error())
	}

	// Send Insert request to datastore
	bodyBytes := InsertTest(t, jsonData)

	// Unmarshal data to QuestionResponse
	var qResponse messages.QuestionResponse
	unmarshalErr := json.Unmarshal(bodyBytes, &qResponse)
	if unmarshalErr != nil {
		t.Errorf(unmarshalErr.Error())
	}

	// Check Error field
	if len(qResponse.Error) > 0 {
		t.Errorf("An error occurred inserting record...")
	}

	// Simulate a client sending a AnswerRequest to the datastore server
	// Build AnswerRequest
	var aRequest messages.AnswerRequest

	// Build Question Request
	aRequest.QuestionID = qRequest.QuestionID
	jsonData, marshalErr = json.Marshal(aRequest)
	if marshalErr != nil {
		t.Errorf("New request error: %s", marshalErr.Error())
	}

	// Send Get request to datastore
	bodyBytes = GetTest(t, jsonData)

	// Build AnswerResponse
	var aResponse messages.AnswerResponse
	unmarshalErr = json.Unmarshal(bodyBytes, &aResponse)
	if unmarshalErr != nil {
		t.Errorf(unmarshalErr.Error())
	}

	// Check Error field
	if aResponse.Message == messages.NO_RESULTS_RETURNED_MSG {
		t.Errorf("An unexpectedly 'No results returned' message is returned...")
	}

	// Send Get(2nd) request to datastore
	bodyBytes = GetTest(t, jsonData)

	// Build AnswerResponse
	unmarshalErr = json.Unmarshal(bodyBytes, &aResponse)
	if unmarshalErr != nil {
		t.Errorf(unmarshalErr.Error())
	}

	log.Print("aResponse: ", aResponse)

	// Check Error field
	if aResponse.Message != messages.NO_RESULTS_RETURNED_MSG {
		t.Errorf("'No results returned' message did NOT returned...")
	}
}
