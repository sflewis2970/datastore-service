package models

import (
	"log"
	"testing"

	"github.com/sflewis2970/datastore-service/models/data"
)

const (
	GOCACHE_DRIVER    string = "go-cache"
	MYSQL_DRIVER      string = "mysql"
	POSTGRESQL_DRIVER string = "postgres"
)

func checkDBDriver(t *testing.T, driverName string, gotDBModel data.IDBModel) {
	if gotDBModel == nil {
		t.Errorf("NewDBModel(%v): returned an invalid object", gotDBModel)
		return
	}

	// Test insert question
	var qRequest data.QuestionRequest
	qRequest.QuestionID = "aaaaqqqq"
	qRequest.Question = "What is 4 / 2?"
	qRequest.Answer = "2"
	_, insertErr := gotDBModel.InsertQuestion(qRequest)
	if insertErr != nil {
		t.Error("Error inserting new record...")
		return
	}

	// Test get question
	qt, getErr := gotDBModel.GetQuestion(qRequest.QuestionID)
	if getErr != nil {
		t.Error("Error retrieving record...")
		return
	}

	if qRequest.Question != qt.Question {
		t.Error("Error request question does NOT match retrieved question...")
		return
	}

	if qRequest.Answer != qt.Answer {
		t.Error("Error request answer does NOT match retrieved answer...")
		return
	}

	// Test update question
	qRequest.Category = "general"
	updateRowsAffected, updateErr := gotDBModel.UpdateQuestion(qRequest)
	if updateErr != nil {
		t.Error("Error updating existing record...")
		return
	}

	if driverName != GOCACHE_DRIVER && updateRowsAffected == 0 {
		t.Error("No rows affected when attempting to delete existing record...")
		return
	}

	// Test delete question
	deletedRowsAffected, deleteErr := gotDBModel.DeleteQuestion(qRequest.QuestionID)
	if deleteErr != nil {
		t.Error("Error deleting record...")
	}

	if driverName != GOCACHE_DRIVER && deletedRowsAffected == 0 {
		t.Error("No rows affected when attempting to delete existing record...")
		return
	}

}

func checkInvalidDriver(t *testing.T, driverName string, gotDBModel data.IDBModel) {
	if gotDBModel != nil {
		t.Errorf("NewDBModel(%v): Invalid driver should not generate a valid object", gotDBModel)
	}
}

func TestNewSDBModel(t *testing.T) {
	// Initialize logging
	log.SetFlags(log.Ldate | log.Lshortfile)

	// Test cases
	testCases := []struct {
		testName   string
		driverName string
	}{
		{testName: "GoCache driver test", driverName: GOCACHE_DRIVER},
		{testName: "MySQL driver test", driverName: MYSQL_DRIVER},
		{testName: "PostgreSQL driver test", driverName: POSTGRESQL_DRIVER},
		{testName: "No driver test", driverName: ""},
		{testName: "Bad driver test", driverName: "baddrivername"},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			gotDBModel := NewDBModel(tc.driverName)

			switch tc.driverName {
			case "go-cache":
				fallthrough
			case "mysql":
				fallthrough
			case "postgres":
				checkDBDriver(t, tc.driverName, gotDBModel)
			default:
				checkInvalidDriver(t, tc.driverName, gotDBModel)
			}
		})
	}
}
