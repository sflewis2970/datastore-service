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
	"github.com/sflewis2970/datastore-service/models"
	"github.com/sflewis2970/datastore-service/models/data"
)

const NoResultsReturnedFoundMsg string = "No results returned"

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
	request, reqErr := http.NewRequest("GET", "/api/v1/ds/addquestion", bytes.NewBuffer(jsonData))
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
	request, reqErr := http.NewRequest("GET", "/api/v1/ds/checkanswer", bytes.NewBuffer(jsonData))
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
	setConfigEnv("go-cache")

	// Initialize controllers object
	Initialize(config.UPDATE_CONFIG_DATA)

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
	var sResponse data.StatusResponse
	unmarshalErr := json.Unmarshal(bodyBytes, &sResponse)
	if unmarshalErr != nil {
		t.Errorf(unmarshalErr.Error())
	}

	// Check status field
	if sResponse.Status != data.StatusCode(models.DS_RUNNING) {
		t.Errorf("Server Status returned is not running, expected running status: got %d", sResponse.Status)
	}
}

func TestInsert(t *testing.T) {
	// Initialize logging
	log.SetFlags(log.Ldate | log.Lshortfile)

	// Set config environment variables
	setConfigEnv("go-cache")

	// Initialize controllers object
	Initialize(config.UPDATE_CONFIG_DATA)

	// Simulate a client sending a QuestionRequest to the datastore server
	// Question Request
	var qRequest data.QuestionRequest

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
	var qResponse data.QuestionResponse
	unmarshalErr := json.Unmarshal(bodyBytes, &qResponse)
	if unmarshalErr != nil {
		t.Errorf(unmarshalErr.Error())
	}

	// Check Error field
	if len(qResponse.Error) != 0 {
		t.Errorf("An error occurred inserting record...")
	}
}

func TestGetBeforeInsert(t *testing.T) {
	// Initialize logging
	log.SetFlags(log.Ldate | log.Lshortfile)

	// Set config environment variables
	setConfigEnv("go-cache")

	// Initialize controllers object
	Initialize(config.UPDATE_CONFIG_DATA)

	// Simulate a client sending a QuestionRequest to the datastore server
	// Build Question Request
	var qRequest data.QuestionRequest

	// Build QuestionRequest
	qRequest.QuestionID = "abcdefgh"
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

	// Build QuestionResponse
	var qResponse data.QuestionResponse
	unmarshalErr := json.Unmarshal(bodyBytes, &qResponse)
	if unmarshalErr != nil {
		t.Errorf(unmarshalErr.Error())
	}

	// Check Error field
	if len(qResponse.Error) != 0 {
		t.Errorf("An error occurred inserting record...")
	}

	// Simulate a client sending a AnswerRequest to the datastore server
	// Build AnswerRequest
	var aRequest data.AnswerRequest

	// Build Question Request
	aRequest.QuestionID = qRequest.QuestionID
	jsonData, marshalErr = json.Marshal(aRequest)
	if marshalErr != nil {
		t.Errorf("New request error: %s", marshalErr.Error())
	}

	// Send AddQuestion request to datastore
	bodyBytes = GetTest(t, jsonData)

	var aResponse data.AnswerResponse
	unmarshalErr = json.Unmarshal(bodyBytes, &aResponse)
	if unmarshalErr != nil {
		t.Errorf(unmarshalErr.Error())
	}

	// Check Question field
	if qRequest.Question != aResponse.Question {
		t.Errorf("The question fields unexpectedly does NOT match")
	}
}

func TestGetAfterInsert(t *testing.T) {
	// Initialize logging
	log.SetFlags(log.Ldate | log.Lshortfile)

	// Set config environment variables
	setConfigEnv("go-cache")

	// Initialize controllers object
	Initialize(config.UPDATE_CONFIG_DATA)

	// Simulate a client sending a AnswerRequest to the datastore server
	// Build AnswerRequest
	var aRequest data.AnswerRequest

	// Build Question Request
	aRequest.QuestionID = "abcdefgh"
	jsonData, marshalErr := json.Marshal(aRequest)
	if marshalErr != nil {
		t.Errorf("New request error: %s", marshalErr.Error())
	}

	// Send CheckAnswer request to datastore
	bodyBytes := GetTest(t, jsonData)

	// Build AnswerResponse
	var aResponse data.AnswerResponse
	unmarshalErr := json.Unmarshal(bodyBytes, &aResponse)
	if unmarshalErr != nil {
		t.Errorf(unmarshalErr.Error())
	}

	// Check Error field
	if aResponse.Message != NoResultsReturnedFoundMsg {
		t.Errorf("An unexpectedly message returned...")
	}
}

func TestGetAfterDelete(t *testing.T) {
	// Initialize logging
	log.SetFlags(log.Ldate | log.Lshortfile)

	// Set config environment variables
	setConfigEnv("go-cache")

	// Initialize controllers object
	Initialize(config.UPDATE_CONFIG_DATA)

	// Simulate a client sending a AnswerRequest to the datastore server
	// Build AnswerRequest
	var aRequest data.AnswerRequest

	// Build Question Request
	aRequest.QuestionID = "abcdefgh"
	jsonData, marshalErr := json.Marshal(aRequest)
	if marshalErr != nil {
		t.Errorf("New request error: %s", marshalErr.Error())
	}

	// Send CheckAnswer request to datastore
	bodyBytes := GetTest(t, jsonData)

	// Build AnswerResponse
	var aResponse data.AnswerResponse
	unmarshalErr := json.Unmarshal(bodyBytes, &aResponse)
	if unmarshalErr != nil {
		t.Errorf(unmarshalErr.Error())
	}

	// Check Error field
	if aResponse.Message != NoResultsReturnedFoundMsg {
		t.Errorf("An unexpectedly message returned...")
	}
}
