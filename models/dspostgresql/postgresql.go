package dspostgresql

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/sflewis2970/datastore-service/config"
	"github.com/sflewis2970/datastore-service/models/data"
)

var postgreSQLModel *dbModel

type dbModel struct {
	cfgData *config.ConfigData
}

// Open database
func (dbm *dbModel) Open(driverName string) (*sql.DB, error) {
	log.Println("Opening PostgreSQL database")

	// Open database connection
	dataSourceName := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", dbm.cfgData.PostGreSQL.Host, dbm.cfgData.PostGreSQL.Port, dbm.cfgData.PostGreSQL.User, "devStation", "main")
	db, openErr := sql.Open(driverName, dataSourceName)

	if openErr != nil {
		log.Print("Error opening database server: ", openErr.Error())
		return nil, openErr
	}

	return db, nil
}

// Ping database server by verifying the database connection is active
func (dbm *dbModel) Ping() error {
	db, openErr := dbm.Open(dbm.cfgData.PostGreSQL.DriverName)
	if openErr != nil {
		return openErr
	}
	defer db.Close()

	pingErr := db.Ping()
	if openErr != nil {
		log.Print("Error pinging database server")
		return pingErr
	}

	return nil
}

// Insert a single record into table
func (dbm *dbModel) InsertQuestion(qRequest data.QuestionRequest) (int64, error) {
	db, openErr := dbm.Open(dbm.cfgData.PostGreSQL.DriverName)
	if openErr != nil {
		return data.RESULTS_DEFAULT, openErr
	}
	defer db.Close()

	log.Print("Adding a new record to the database")
	queryStr := "insert into trivia VALUES ($1, $2, $3, $4);"
	sqlDB, execErr := db.Exec(queryStr, qRequest.QuestionID, qRequest.Question, qRequest.Category, qRequest.Answer)
	if execErr != nil {
		log.Print("Error inserting data into database table")
		return data.RESULTS_DEFAULT, execErr
	}

	rowsAffected, rowsAffectedErr := sqlDB.RowsAffected()
	if rowsAffectedErr != nil {
		log.Print("Error getting rows affected: ", rowsAffectedErr.Error())
		return data.RESULTS_DEFAULT, rowsAffectedErr
	}

	return rowsAffected, nil
}

// Get a single record from table
func (dbm *dbModel) GetQuestion(questionID string) (data.QuestionTable, error) {
	db, openErr := dbm.Open(dbm.cfgData.PostGreSQL.DriverName)
	if openErr != nil {
		return data.QuestionTable{}, openErr
	}
	defer db.Close()

	var qTable data.QuestionTable

	log.Print("Getting a single record from the database")
	queryStr := "SELECT question, category, answer FROM trivia WHERE question_id = $1;"
	scanErr := db.QueryRow(queryStr, questionID).Scan(&qTable.Question, &qTable.Category, &qTable.Answer)
	if scanErr != nil && scanErr != sql.ErrNoRows {
		log.Print("Error getting record from database: ", scanErr.Error())
		return data.QuestionTable{}, scanErr
	}

	return qTable, nil
}

// Update a single record in table
func (dbm *dbModel) UpdateQuestion(qRequest data.QuestionRequest) (int64, error) {
	db, openErr := dbm.Open(dbm.cfgData.PostGreSQL.DriverName)
	if openErr != nil {
		return data.RESULTS_DEFAULT, openErr
	}
	defer db.Close()

	log.Println("Updating a single record in the database")
	queryStr := "UPDATE trivia SET question = $2, category = $3, answer = $4 WHERE question_id = $1"
	sqlDB, execErr := db.Exec(queryStr, qRequest.QuestionID, qRequest.Question, qRequest.Category, qRequest.Answer)
	if execErr != nil {
		log.Print("Error getting data from database table")
		return data.RESULTS_DEFAULT, execErr
	}

	rowsAffected, rowsAffectedErr := sqlDB.RowsAffected()
	if rowsAffectedErr != nil {
		log.Print("Error getting rows affected: ", rowsAffectedErr.Error())
		return data.RESULTS_DEFAULT, nil
	}

	return rowsAffected, nil
}

// Delete a single record from table
func (dbm *dbModel) DeleteQuestion(questionID string) (int64, error) {
	db, openErr := dbm.Open(dbm.cfgData.PostGreSQL.DriverName)
	if openErr != nil {
		return data.RESULTS_DEFAULT, openErr
	}
	defer db.Close()

	log.Println("deleting a single record from the database")
	queryStr := "DELETE FROM trivia WHERE question_id = $1"
	sqlDB, execErr := db.Exec(queryStr, questionID)
	if execErr != nil {
		log.Print("Error deleting data from database table")
		return data.RESULTS_DEFAULT, execErr
	}

	rowsAffected, rowsAffectedErr := sqlDB.RowsAffected()
	if rowsAffectedErr != nil {
		log.Print("Error getting rows affected: ", rowsAffectedErr.Error())
		return data.RESULTS_DEFAULT, nil
	}

	return rowsAffected, nil
}

func GetPostGreSQLModel() *dbModel {
	if postgreSQLModel == nil {
		log.Print("Creating PostgreSQL database model")

		postgreSQLModel = new(dbModel)
		cfg, getErr := config.GetConfig()
		if getErr != nil {
			log.Fatal("Error getting config info: ", getErr)
		}

		postgreSQLModel.cfgData = cfg.GetConfigData()
	}

	return postgreSQLModel
}
