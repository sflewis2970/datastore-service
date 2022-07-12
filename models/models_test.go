package models

import (
	"log"
	"os"
	"testing"

	"github.com/sflewis2970/datastore-service/config"
	"github.com/sflewis2970/datastore-service/models/messages"
)

func checkDBDriver(t *testing.T, driverName string, gotDBModel messages.IDBModel) {
	if gotDBModel == nil {
		t.Errorf("NewDBModel(%v): returned an invalid object", gotDBModel)
		return
	}

	// Test insert question
	var qRequest messages.QuestionRequest
	qRequest.QuestionID = "aaaaqqqq"
	qRequest.Question = "What is 4 / 2?"
	qRequest.Answer = "2"

	_, insertErr := gotDBModel.Insert(qRequest)
	if insertErr != nil {
		t.Error("Error inserting new record...")
		return
	}

	// Test get question
	qt, getErr := gotDBModel.Get(qRequest.QuestionID)
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
	updateRowsAffected, updateErr := gotDBModel.Update(qRequest)
	if updateErr != nil {
		t.Error("Error updating existing record...")
		return
	}

	if driverName != config.GOCACHE_DRIVER && updateRowsAffected == 0 {
		t.Error("No rows affected when attempting to update existing record...")
		return
	}

	// Test delete question
	deletedRowsAffected, deleteErr := gotDBModel.Delete(qRequest.QuestionID)
	if deleteErr != nil {
		t.Error("Error deleting record...")
	}

	if driverName != config.GOCACHE_DRIVER && deletedRowsAffected == 0 {
		t.Error("No rows affected when attempting to delete existing record...")
		return
	}

}

func checkInvalidDriver(t *testing.T, driverName string, gotDBModel messages.IDBModel) {
	if gotDBModel != nil {
		t.Errorf("NewDBModel(%v): Invalid driver should not generate a valid object", gotDBModel)
	}
}

// The envionment variables will be set on the server
// For testing, set the variables manually
func setConfigEnv(driverName string) error {
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
	setErr = os.Setenv(config.ACTIVEDRIVER, driverName)
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

func TestNewSDBModel(t *testing.T) {
	// Initialize logging
	log.SetFlags(log.Ldate | log.Lshortfile)

	// new model
	model := New()

	// Test cases
	testCases := []struct {
		testName   string
		testActive bool
		driverName string
	}{
		{testName: "GoCache driver test", testActive: true, driverName: config.GOCACHE_DRIVER},
		{testName: "GoRedis driver test", testActive: false, driverName: config.GOREDIS_DRIVER},
		{testName: "MySQL driver test", testActive: false, driverName: config.MYSQL_DRIVER},
		{testName: "PostgreSQL driver test", testActive: false, driverName: config.POSTGRESQL_DRIVER},
		{testName: "No driver test", testActive: true, driverName: ""},
		{testName: "Bad driver test", testActive: true, driverName: "baddrivername"},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			setConfigEnv(tc.driverName)

			_, getCfgDataErr := config.Get().GetData(config.REFRESH_CONFIG_DATA)
			if getCfgDataErr != nil {
				t.Errorf("Error getting config data...")
				return
			}

			gotDBModel := model.NewDBModel(tc.driverName)

			switch tc.driverName {
			case config.GOCACHE_DRIVER:
				fallthrough
			case config.MYSQL_DRIVER:
				fallthrough
			case config.POSTGRESQL_DRIVER:
				if tc.testActive {
					checkDBDriver(t, tc.driverName, gotDBModel)
				}
			default:
				if tc.testActive {
					checkInvalidDriver(t, tc.driverName, gotDBModel)
				}
			}
		})
	}
}
