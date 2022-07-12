package models

import (
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/sflewis2970/datastore-service/common"
	"github.com/sflewis2970/datastore-service/config"
	"github.com/sflewis2970/datastore-service/models/dspostgresql"
	"github.com/sflewis2970/datastore-service/models/gocache"
	"github.com/sflewis2970/datastore-service/models/goredis"
	"github.com/sflewis2970/datastore-service/models/messages"
)

type Model struct {
	cfgData *config.ConfigData
	dbModel messages.IDBModel
}

func (m *Model) Status() (messages.StatusResponse, error) {
	// Load config data
	log.Print("Getting active datastore driver from config data...")

	// Status Response
	var sResponse messages.StatusResponse

	// DB Model
	m.dbModel = m.NewDBModel(m.cfgData.ActiveDriver)

	sResponse.Timestamp = common.GetFormattedTime(time.Now(), "Mon Jan 2 15:04:05 2006")
	if m.dbModel != nil {
		// Update Server response fields
		pingErr := m.dbModel.Ping()

		if pingErr != nil {
			sResponse.Error = pingErr.Error()
			sResponse.Status = messages.StatusCode(messages.DS_UNAVAILABLE)
			return sResponse, pingErr
		} else {
			log.Print("preparing status to be sent back to client...")
			sResponse.Status = messages.StatusCode(messages.DS_RUNNING)
		}
	} else {
		sResponse.Error = "dbModel not created"
		sResponse.Status = messages.StatusCode(messages.DS_UNAVAILABLE)
	}

	return sResponse, nil
}

func (m *Model) Insert(qRequest messages.QuestionRequest) (messages.QuestionResponse, error) {
	m.dbModel = m.NewDBModel(m.cfgData.ActiveDriver)
	rowsAffected, insertErr := m.dbModel.Insert(qRequest)

	var qResponse messages.QuestionResponse

	// Update timestamp
	qResponse.Timestamp = common.GetFormattedTime(time.Now(), "Mon Jan 2 15:04:05 2006")

	if insertErr != nil {
		// Display a log message
		errMsg := "Insertion error: " + insertErr.Error()
		log.Printf(errMsg)

		// Update response fields
		qResponse.Error = errMsg

		return qResponse, errors.New(errMsg)
	} else {
		if rowsAffected > 0 {
			log.Print("rows affected: ", rowsAffected)
		}

		// Build QuestionResponse
		qResponse.QuestionID = qRequest.QuestionID
		qResponse.Question = qRequest.Question
		qResponse.Category = qRequest.Category

		// Display a log message
		log.Print("sending response to client...")

		// Build QuestionResponse message
		qResponse.Message = "Record added to the datastore"
	}

	return qResponse, nil
}

func (m *Model) Get(aRequest messages.AnswerRequest) (messages.AnswerResponse, error) {
	// use dbModel to execute SQL command
	m.dbModel = m.NewDBModel(m.cfgData.ActiveDriver)

	var aResponse messages.AnswerResponse
	qt, getErr := m.dbModel.Get(aRequest.QuestionID)
	if getErr != nil {
		// Display a log message
		errMsg := "Get error: " + getErr.Error()
		log.Printf(errMsg)

		// Update response fields
		aResponse.Error = errMsg

		return aResponse, errors.New(errMsg)
	}

	// Build AnswerResponse message
	aResponse.Timestamp = common.GetFormattedTime(time.Now(), "Mon Jan 2 15:04:05 2006")
	aResponse.Question = qt.Question
	aResponse.Category = qt.Category
	aResponse.Answer = qt.Answer

	// Since sql.QueryRow wraps no results inside error messages,
	// when an error is returned a check needs to be made
	// if ErrNoRows is returned.
	if getErr != nil && getErr != sql.ErrNoRows {
		errMsg := "Error getting record: " + getErr.Error()
		aResponse.Error = errMsg

		return aResponse, errors.New(errMsg)
	} else if len(qt.Question) > 0 {
		log.Print("Question retrieved processing message...")
		// Build Response Message

		// delete record from DB once the client answers the question
		// Whether the answer is correct or not
		_, delErr := m.dbModel.Delete(aRequest.QuestionID)
		if delErr != nil {
			errMsg := "Error deleting record: " + delErr.Error()
			aResponse.Error = delErr.Error()

			return aResponse, errors.New(errMsg)
		}
	} else {
		aResponse.Message = messages.NO_RESULTS_RETURNED_MSG
	}

	return aResponse, nil
}

func (m *Model) Update(qRequest messages.QuestionRequest) (messages.QuestionResponse, error) {
	m.dbModel = m.NewDBModel(m.cfgData.ActiveDriver)

	var qResponse messages.QuestionResponse
	_, updateErr := m.dbModel.Update(qRequest)
	if updateErr != nil {
		// Display a log message
		errMsg := "Error updating record: " + updateErr.Error()
		log.Printf(errMsg)

		// Update response fields
		qResponse.Error = errMsg

		return qResponse, errors.New(errMsg)
	} else {
		// Build QuestionResponse message
		qResponse.Message = "Updated question record in database"
	}

	return qResponse, nil
}

func (m *Model) Delete(questionID string) (messages.QuestionResponse, error) {
	m.dbModel = m.NewDBModel(m.cfgData.ActiveDriver)

	_, delErr := m.dbModel.Delete(questionID)

	var qResponse messages.QuestionResponse
	if delErr != nil {
		// Display a log message
		errMsg := "Error updating record: " + delErr.Error()
		log.Printf(errMsg)

		qResponse.Error = delErr.Error()

		// Update response fields
		qResponse.Error = errMsg

		return qResponse, errors.New(errMsg)

	} else {
		qResponse.Message = "Question with QuestionID = " + questionID + " has been deleted"
	}

	return qResponse, nil
}

func (m *Model) NewDBModel(activeDriver string) messages.IDBModel {
	if m.dbModel == nil {
		switch activeDriver {
		case config.GOCACHE_DRIVER:
			return gocache.GetGoCacheModel(m.cfgData)
		case config.GOREDIS_DRIVER:
			return goredis.GetGoRedisModel(m.cfgData)
		case config.POSTGRESQL_DRIVER:
			return dspostgresql.GetPostGreSQLModel(m.cfgData)
		default:
			log.Print("Unsupported database driver, active driver: ", activeDriver)
		}
	}

	return m.dbModel
}

func New() *Model {
	log.Print("Creating model object...")
	model := new(Model)

	// Load config data
	var cfgDataErr error
	model.cfgData, cfgDataErr = config.Get().GetData()
	if cfgDataErr != nil {
		log.Print("Error loading config data...: ", cfgDataErr)
		return nil
	}

	return model
}
