package controllers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/sflewis2970/datastore-service/common"
	"github.com/sflewis2970/datastore-service/models"
	"github.com/sflewis2970/datastore-service/models/data"
)

func Status(rw http.ResponseWriter, r *http.Request) {
	// Display a log message
	log.Print("client requesting server status...")

	// Status Response
	var sResponse data.StatusResponse

	// Load config data
	log.Print("Getting active datastore driver from config data...")

	// DB Model
	dbModel := models.NewDBModel(ctrlr.cfgData.ActiveDriver)

	sResponse.Timestamp = common.GetFormattedTime(time.Now(), "Mon Jan 2 15:04:05 2006")
	if dbModel != nil {
		// Update Server response fields
		pingErr := dbModel.Ping()

		if pingErr != nil {
			sResponse.Error = pingErr.Error()
			sResponse.Status = data.StatusCode(models.DS_UNAVAILABLE)
		} else {
			log.Print("preparing status to be sent back to client...")
			sResponse.Status = data.StatusCode(models.DS_RUNNING)
		}
	} else {
		sResponse.Error = "dbModel not created"
		sResponse.Status = data.StatusCode(models.DS_UNAVAILABLE)
	}

	// Write JSON to stream
	json.NewEncoder(rw).Encode(sResponse)
}

func InsertQuestion(rw http.ResponseWriter, r *http.Request) {
	ctrlr.dbMutex.Lock()
	defer ctrlr.dbMutex.Unlock()

	// Display a log message
	log.Print("client requested insert action...")

	// Question Request
	var qRequest data.QuestionRequest

	// Decode request into JSON format
	json.NewDecoder(r.Body).Decode(&qRequest)

	dbModel := models.NewDBModel(ctrlr.cfgData.ActiveDriver)
	rowsAffected, insertErr := dbModel.InsertQuestion(qRequest)

	var qResponse data.QuestionResponse

	// Update timestamp
	qResponse.Timestamp = common.GetFormattedTime(time.Now(), "Mon Jan 2 15:04:05 2006")

	if insertErr != nil {
		// Display a log message
		errMsg := "Insertion error: " + insertErr.Error()
		log.Printf(errMsg)

		// Update response fields
		qResponse.Error = errMsg
	} else {
		if rowsAffected > 0 {
			log.Print("rows affected: ", rowsAffected)
		}

		// Display a log message
		log.Print("sending response to client...")

		// Build QuestionResponse message
		qResponse.Message = "Add question record to the database"
	}

	// Write JSON to stream
	json.NewEncoder(rw).Encode(qResponse)
}

func GetQuestion(rw http.ResponseWriter, r *http.Request) {
	ctrlr.dbMutex.Lock()
	defer ctrlr.dbMutex.Unlock()

	// Get question ID from query parameter
	questionID := r.URL.Query().Get("questionid")

	// Display a log message
	log.Print("client requested to getquestion...")

	// Connect to database
	dbModel := models.NewDBModel(ctrlr.cfgData.ActiveDriver)

	// Get item from database
	qt, getErr := dbModel.GetQuestion(questionID)

	var qResponse data.QuestionResponse
	if getErr != nil {
		qResponse.Error = getErr.Error()

		// Update status
		rw.WriteHeader(http.StatusInternalServerError)

		// Display a log message
		log.Print("Error getting question...")

	} else {
		// Display a log message
		log.Print("sending response to client...")

		// Build Response Message
		qResponse.QuestionID = questionID
		qResponse.Question = qt.Question
		qResponse.Category = qt.Category
		qResponse.Answer = qt.Answer
	}

	// Write JSON to stream
	json.NewEncoder(rw).Encode(qResponse)
}

func UpdateQuestion(rw http.ResponseWriter, r *http.Request) {
	ctrlr.dbMutex.Lock()
	defer ctrlr.dbMutex.Unlock()

	var question data.QuestionRequest

	// Display a log message
	log.Print("received update request from client...")

	// Decode request into JSON format
	json.NewDecoder(r.Body).Decode(&question)

	dbModel := models.NewDBModel(ctrlr.cfgData.ActiveDriver)

	var qResponse data.QuestionResponse
	rowsAffected, updateErr := dbModel.UpdateQuestion(question)
	if updateErr != nil {
		qResponse.Error = updateErr.Error()

		// Update status
		rw.WriteHeader(http.StatusInternalServerError)
	} else {
		if rowsAffected > 0 {
			log.Print("rows affected: ", rowsAffected)
		}

		// Display a log message
		log.Print("sending response to client...")

		// Build QuestionResponse message
		qResponse.Message = "Updated question record in database"
	}

	// Write JSON to stream
	json.NewEncoder(rw).Encode(qResponse)
}

func DeleteQuestion(rw http.ResponseWriter, r *http.Request) {
	ctrlr.dbMutex.Lock()
	defer ctrlr.dbMutex.Unlock()

	// Get question ID from query parameter
	questionID := r.URL.Query().Get("questionid")

	// Display a log message
	log.Print("data received from client...")

	dbModel := models.NewDBModel(ctrlr.cfgData.ActiveDriver)

	rowsAffected, delErr := dbModel.DeleteQuestion(questionID)

	var qResponse data.QuestionResponse
	if delErr != nil {
		qResponse.Error = delErr.Error()

		// Update status
		rw.WriteHeader(http.StatusInternalServerError)
	} else {
		if rowsAffected > 0 {
			log.Print("rows affected: ", rowsAffected)
		}

		// Display a log message
		log.Print("sending response to client...")

		qResponse.Message = "Question with QuestionID = " + questionID + " has been deleted"
	}

	// Write JSON to stream
	json.NewEncoder(rw).Encode(qResponse)
}

func CheckAnswer(rw http.ResponseWriter, r *http.Request) {
	ctrlr.dbMutex.Lock()
	defer ctrlr.dbMutex.Unlock()

	// Answer Request
	var aRequest data.AnswerRequest

	// Display a log message
	log.Print("data received from client...")

	// Decode request into JSON format
	json.NewDecoder(r.Body).Decode(&aRequest)

	// use dbModel to execute SQL command
	dbModel := models.NewDBModel(ctrlr.cfgData.ActiveDriver)
	qt, getErr := dbModel.GetQuestion(aRequest.QuestionID)

	// Build AnswerResponse message
	var aResponse data.AnswerResponse
	aResponse.Timestamp = common.GetFormattedTime(time.Now(), "Mon Jan 2 15:04:05 2006")
	aResponse.Question = qt.Question
	aResponse.Category = qt.Category
	aResponse.Answer = qt.Answer
	aResponse.Response = aRequest.Response

	// Since sql.QueryRow wraps no results inside error messages,
	// when an error is returned a check needs to be made
	// if ErrNoRows is returned.
	if getErr != nil && getErr != sql.ErrNoRows {
		aResponse.Error = getErr.Error()

		// Update status
		rw.WriteHeader(http.StatusInternalServerError)
	} else if len(qt.Question) > 0 {
		log.Print("Question retrieved processing message...")
		// Build Response Message
		if qt.Answer == aRequest.Response {
			aResponse.Correct = true
			aResponse.Message = ctrlr.cfgData.Messages.CongratsMsg
		} else {
			aResponse.Correct = false
			aResponse.Message = ctrlr.cfgData.Messages.TryAgainMsg
		}

		// delete record from DB once the client answers the question
		// Whether the answer is correct or not
		rowsAffected, delErr := dbModel.DeleteQuestion(aRequest.QuestionID)
		if delErr != nil {
			aResponse.Error = delErr.Error()

			// Update status
			rw.WriteHeader(http.StatusInternalServerError)
		} else {
			if rowsAffected > 0 {
				log.Print("rows affected: ", rowsAffected)
			}
		}
	} else {
		aResponse.Message = "No results returned"
	}

	// Write JSON to stream
	json.NewEncoder(rw).Encode(aResponse)
}