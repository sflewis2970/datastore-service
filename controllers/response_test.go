package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sflewis2970/datastore-service/models"
	"github.com/sflewis2970/datastore-service/models/data"
)

var ErrQuestionNotFound = errors.New("Question not found")

func AddQuestionTest(t *testing.T, jsonData []byte) []byte {
	// Create new request
	request, reqErr := http.NewRequest("GET", "/api/v1/ds/addquestion", bytes.NewBuffer(jsonData))
	if reqErr != nil {
		t.Errorf("Could not create request.\n")
	}

	// Setup recoder
	rRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(InsertQuestion)
	handler.ServeHTTP(rRecorder, request)

	// Check response code
	status := rRecorder.Code
	if status != http.StatusOK {
		t.Errorf("handler returned invalid status code: got %d, expected: %d\n", status, http.StatusOK)
	}

	// Unmarshal JSON
	return rRecorder.Body.Bytes()
}

func CheckAnswerTest(t *testing.T, jsonData []byte) []byte {
	// Create new request
	request, reqErr := http.NewRequest("GET", "/api/v1/ds/checkanswer", bytes.NewBuffer(jsonData))
	if reqErr != nil {
		t.Errorf("Could not create request.\n")
	}

	// Setup recoder
	rRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(CheckAnswer)
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

	// Initialize controllers object
	InitializeController()

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

func TestCheckAnswerWithCorrectAnswer(t *testing.T) {
	// Initialize logging
	log.SetFlags(log.Ldate | log.Lshortfile)

	// Initialize controllers object
	InitializeController()

	// Simulate a client sending a QuestionRequest to the datastore server
	// Question Request
	var qRequest data.QuestionRequest

	// Build Question Request
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
	bodyBytes := AddQuestionTest(t, jsonData)

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

	// Simulate a client sending an AnswerRequest to the datastore server
	// Answer Request
	var aRequest data.AnswerRequest

	// Build Question Request
	aRequest.QuestionID = qRequest.QuestionID
	aRequest.Response = qRequest.Answer
	jsonData, marshalErr = json.Marshal(aRequest)

	if marshalErr != nil {
		t.Errorf("New request error: %s", marshalErr.Error())
	}

	// Send AddQuestion request to datastore
	bodyBytes = CheckAnswerTest(t, jsonData)

	// Unmarshal data to AnswerResponse
	var aResponse data.AnswerResponse
	unmarshalErr = json.Unmarshal(bodyBytes, &aResponse)
	if unmarshalErr != nil {
		t.Errorf(unmarshalErr.Error())
	}

	// Check Question field
	if qRequest.Question != aResponse.Question {
		t.Errorf("The question fields do NOT match")
	}
}

func TestCheckAnswerWithIncorrectAnswer(t *testing.T) {
	// Initialize logging
	log.SetFlags(log.Ldate | log.Lshortfile)

	// Initialize controllers object
	InitializeController()

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
	bodyBytes := AddQuestionTest(t, jsonData)

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
	aRequest.Response = "Morpheous"
	jsonData, marshalErr = json.Marshal(aRequest)
	if marshalErr != nil {
		t.Errorf("New request error: %s", marshalErr.Error())
	}

	// Send AddQuestion request to datastore
	bodyBytes = CheckAnswerTest(t, jsonData)

	var aResponse data.AnswerResponse
	unmarshalErr = json.Unmarshal(bodyBytes, &aResponse)
	if unmarshalErr != nil {
		t.Errorf(unmarshalErr.Error())
	}

	// Check Question field
	if qRequest.Question != aResponse.Question {
		t.Errorf("The question fields unexpectedly does NOT match")
	}

	// Check Answer field
	if aRequest.Response == aResponse.Answer {
		t.Errorf("The response unexpectedly matches answer")
	}
}

func TestCheckAnswerWithoutAddQuestion(t *testing.T) {
	// Initialize logging
	log.SetFlags(log.Ldate | log.Lshortfile)

	// Initialize controllers object
	InitializeController()

	// Simulate a client sending a AnswerRequest to the datastore server
	// Build AnswerRequest
	var aRequest data.AnswerRequest

	// Build Question Request
	aRequest.QuestionID = "abcdefgh"
	aRequest.Response = "Morpheous"
	jsonData, marshalErr := json.Marshal(aRequest)
	if marshalErr != nil {
		t.Errorf("New request error: %s", marshalErr.Error())
	}

	// Send CheckAnswer request to datastore
	bodyBytes := CheckAnswerTest(t, jsonData)

	// Build AnswerResponse
	var aResponse data.AnswerResponse
	unmarshalErr := json.Unmarshal(bodyBytes, &aResponse)
	if unmarshalErr != nil {
		t.Errorf(unmarshalErr.Error())
	}

	// Check Error field
	log.Print("Answer Response error: ", aResponse.Error)
	if aResponse.Error == ErrQuestionNotFound.Error() {
		t.Errorf("An unexpectedly error occurred...")
	}
}
