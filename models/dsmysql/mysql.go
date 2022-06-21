package dsmysql

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/sflewis2970/datastore-service/config"
	"github.com/sflewis2970/datastore-service/models/data"
)

var mySQLModel *dbModel

type dbModel struct {
	cfgData *config.ConfigData
}

// Open database
func (dbm *dbModel) Open(sqlDriverName string) (*sql.DB, error) {
	log.Println("Opening MySQL database")

	// Open database connection
	dataSourceName := dbm.cfgData.MySQL.Connection + "maindb"
	db, openErr := sql.Open(sqlDriverName, dataSourceName)

	if openErr != nil {
		log.Println("Error opening MySQL database...")
		return nil, openErr
	}

	return db, nil
}

// Ping database server by verifying the database connection is active
func (dbm *dbModel) Ping() error {
	db, openErr := dbm.Open(dbm.cfgData.MySQL.DriverName)
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
	db, openErr := dbm.Open(dbm.cfgData.MySQL.DriverName)
	if openErr != nil {
		return data.RESULTS_DEFAULT, openErr
	}
	defer db.Close()

	log.Print("Adding a new record to MySQL database")
	queryStr := "INSERT INTO trivia VALUES (?, ?, ?, ?)"
	sqlDB, execErr := db.Exec(queryStr, qRequest.QuestionID, qRequest.Question, qRequest.Category, qRequest.Answer)
	if execErr != nil {
		log.Print("Error insert record into datbase")
	}

	rowsAffected, rowsAffectedErr := sqlDB.RowsAffected()
	if rowsAffectedErr != nil {
		log.Print("Error getting rows affected: ", rowsAffectedErr.Error())
	}

	return rowsAffected, nil
}

// Get a single record from table
func (dbm *dbModel) GetQuestion(questionID string) (data.QuestionTable, error) {
	db, openErr := dbm.Open(dbm.cfgData.MySQL.DriverName)
	if openErr != nil {
		return data.QuestionTable{}, openErr
	}
	defer db.Close()

	var qTable data.QuestionTable

	log.Print("Getting record from MySQL database")
	queryStr := "SELECT question, category, answer FROM trivia WHERE question_id = ?"
	queryErr := db.QueryRow(queryStr, questionID).Scan(&qTable.Question, &qTable.Category, &qTable.Answer)
	if queryErr != nil {
		if queryErr != sql.ErrNoRows {
			log.Print("Error getting results from MySQL database...: ", queryErr.Error())
		}

		return data.QuestionTable{}, queryErr
	}

	return qTable, nil
}

// Update a single record in table
func (dbm *dbModel) UpdateQuestion(qRequest data.QuestionRequest) (int64, error) {
	db, openErr := dbm.Open(dbm.cfgData.MySQL.DriverName)
	if openErr != nil {
		return data.RESULTS_DEFAULT, openErr
	}
	defer db.Close()

	log.Println("Updating record in MySQL database")
	queryStr := "UPDATE trivia SET question = ?, category = ?, answer = ? WHERE question_id = ?"
	sqlDB, execErr := db.Exec(queryStr, qRequest.Question, qRequest.Category, qRequest.Answer, qRequest.QuestionID)
	if execErr != nil {
		log.Print("Error updating database record")
		return data.RESULTS_DEFAULT, execErr
	}

	rowsAffected, rowsAffectedErr := sqlDB.RowsAffected()
	if rowsAffectedErr != nil {
		log.Print("Error getting rows affected: ", rowsAffectedErr.Error())
		return data.RESULTS_DEFAULT, rowsAffectedErr
	}

	return rowsAffected, nil
}

// Delete a single record from table
func (dbm *dbModel) DeleteQuestion(questionID string) (int64, error) {
	db, openErr := dbm.Open(dbm.cfgData.MySQL.DriverName)
	if openErr != nil {
		return data.RESULTS_DEFAULT, openErr
	}
	defer db.Close()

	log.Println("deleting record from MySQL database")
	queryStr := "DELETE FROM trivia WHERE question_id = ?"
	sqlDB, execErr := db.Exec(queryStr, questionID)
	if execErr != nil {
		log.Print("Error deleting record from database")
		return data.RESULTS_DEFAULT, execErr
	}

	rowsAffected, rowsAffectedErr := sqlDB.RowsAffected()
	if rowsAffectedErr != nil {
		log.Print("Error getting rows affected: ", rowsAffectedErr.Error())
		return data.RESULTS_DEFAULT, rowsAffectedErr
	}

	return rowsAffected, nil
}

func GetMySQLModel() *dbModel {
	if mySQLModel == nil {
		log.Print("Creating MySQL database model")

		mySQLModel = new(dbModel)
		cfg, getErr := config.GetConfig()
		if getErr != nil {
			log.Fatal("Error getting config info: ", getErr)
		}

		mySQLModel.cfgData = cfg.GetConfigData()
	}

	return mySQLModel
}
