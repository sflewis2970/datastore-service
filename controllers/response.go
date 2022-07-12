package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/sflewis2970/datastore-service/models/messages"
)

func Status(rw http.ResponseWriter, r *http.Request) {
	// Display a log message
	log.Print("client requesting server status...")

	// Get Datastore Server Status
	sResponse, statusErr := controller.dataModel.Status()
	if statusErr != nil {
		rw.WriteHeader(http.StatusInternalServerError)
	}

	// Write JSON to stream
	json.NewEncoder(rw).Encode(sResponse)
}

func Insert(rw http.ResponseWriter, r *http.Request) {
	controller.dbMutex.Lock()
	defer controller.dbMutex.Unlock()

	// Display a log message
	log.Print("Insert action requested...")

	// Question Request
	var qRequest messages.QuestionRequest

	// Decode request into JSON format
	json.NewDecoder(r.Body).Decode(&qRequest)

	// Send Insert request
	qResponse, insertErr := controller.dataModel.Insert(qRequest)
	if insertErr != nil {
		rw.WriteHeader(http.StatusInternalServerError)
	}

	// Write JSON to stream
	json.NewEncoder(rw).Encode(qResponse)
}

// Get receives a request and
func Get(rw http.ResponseWriter, r *http.Request) {
	controller.dbMutex.Lock()
	defer controller.dbMutex.Unlock()

	// Answer Request
	var aRequest messages.AnswerRequest

	// Display a log message
	log.Print("data received from client...")

	// Decode request into JSON format
	json.NewDecoder(r.Body).Decode(&aRequest)

	// Send Answer Request
	aResponse, getErr := controller.dataModel.Get(aRequest)
	if getErr != nil {
		rw.WriteHeader(http.StatusInternalServerError)
	}

	// Write JSON to stream
	json.NewEncoder(rw).Encode(aResponse)
}

func Update(rw http.ResponseWriter, r *http.Request) {
	controller.dbMutex.Lock()
	defer controller.dbMutex.Unlock()

	var question messages.QuestionRequest

	// Display a log message
	log.Print("received update request from client...")

	// Decode request into JSON format
	json.NewDecoder(r.Body).Decode(&question)

	// Update question
	qResponse, updateErr := controller.dataModel.Update(question)
	if updateErr != nil {
		rw.WriteHeader(http.StatusInternalServerError)
	}

	// Display a log message
	log.Print("sending response to client...")

	// Write JSON to stream
	json.NewEncoder(rw).Encode(qResponse)
}

func Delete(rw http.ResponseWriter, r *http.Request) {
	controller.dbMutex.Lock()
	defer controller.dbMutex.Unlock()

	// Display a log message
	log.Print("data received from client...")

	// Get question ID from query parameter
	questionID := r.URL.Query().Get("questionid")

	// Send delete request
	qResponse, delErr := controller.dataModel.Delete(questionID)
	if delErr != nil {
		rw.WriteHeader(http.StatusInternalServerError)
	}

	// Display a log message
	log.Print("sending response to client...")

	// Write JSON to stream
	json.NewEncoder(rw).Encode(qResponse)
}
