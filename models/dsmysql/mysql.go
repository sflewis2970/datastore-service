package dsmysql

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/sflewis2970/datastore-service/config"
	"github.com/sflewis2970/datastore-service/models/data"
)

const (
	MYSQL_DB_NAME_MSG string = "MYSQL: "
)

const (
	MYSQL_GET_CONFIG_ERROR      string = "Getting config error...: "
	MYSQL_GET_CONFIG_DATA_ERROR string = "Getting config data error...: "
	MYSQL_OPEN_ERROR            string = "Error opening database...: "
	MYSQL_INSERT_ERROR          string = "Error inserting record...: "
	MYSQL_GET_ERROR             string = "Error getting record...: "
	MYSQL_UPDATE_ERROR          string = "Error updating record...: "
	MYSQL_DELETE_ERROR          string = "Error deleting record...: "
	MYSQL_RESULTS_ERROR         string = "Error getting results...: "
	MYSQL_ROWS_AFFECTED_ERROR   string = "Error getting rows affected...: "
	MYSQL_PING_ERROR            string = "Error pinging database server...: "
)

var mySQLModel *dbModel

type dbModel struct {
	cfgData *config.ConfigData
}

// Open database
func (dbm *dbModel) Open(sqlDriverName string) (*sql.DB, error) {
	log.Print(MYSQL_DB_NAME_MSG + "Opening MySQL database")

	// Open database connection
	dataSourceName := dbm.cfgData.MySQL.Connection + "maindb"
	db, openErr := sql.Open(sqlDriverName, dataSourceName)

	if openErr != nil {
		log.Print(MYSQL_DB_NAME_MSG+MYSQL_OPEN_ERROR, openErr)
		return nil, openErr
	}

	return db, nil
}

// Ping database server by verifying the database connection is active
func (dbm *dbModel) Ping() error {
	db, openErr := dbm.Open(dbm.cfgData.ActiveDriver)
	if openErr != nil {
		log.Print(MYSQL_DB_NAME_MSG+MYSQL_OPEN_ERROR, openErr)
		return openErr
	}
	defer db.Close()

	pingErr := db.Ping()
	if openErr != nil {
		log.Print(MYSQL_DB_NAME_MSG+MYSQL_PING_ERROR, pingErr)
		return pingErr
	}

	return nil
}

// Insert a single record into table
func (dbm *dbModel) Insert(qRequest data.QuestionRequest) (int64, error) {
	db, openErr := dbm.Open(dbm.cfgData.ActiveDriver)
	if openErr != nil {
		log.Print(MYSQL_DB_NAME_MSG+MYSQL_OPEN_ERROR, openErr)
		return data.OPEN_ERROR_CODE, openErr
	}
	defer db.Close()

	log.Print("Adding a new record to MySQL database")
	queryStr := "INSERT INTO trivia VALUES (?, ?, ?, ?)"
	sqlDB, execErr := db.Exec(queryStr, qRequest.QuestionID, qRequest.Question, qRequest.Category, qRequest.Answer)
	if execErr != nil {
		log.Print(MYSQL_DB_NAME_MSG+MYSQL_INSERT_ERROR, execErr)
		return data.INSERT_ERROR_CODE, execErr
	}

	rowsAffected, rowsAffectedErr := sqlDB.RowsAffected()
	if rowsAffectedErr != nil {
		log.Print(MYSQL_DB_NAME_MSG+MYSQL_ROWS_AFFECTED_ERROR, rowsAffectedErr.Error())
		return data.ROWS_AFFECTED_ERROR_CODE, rowsAffectedErr
	}

	return rowsAffected, nil
}

// Get a single record from table
func (dbm *dbModel) Get(questionID string) (data.QuestionTable, error) {
	db, openErr := dbm.Open(dbm.cfgData.ActiveDriver)
	if openErr != nil {
		log.Print(MYSQL_DB_NAME_MSG+MYSQL_OPEN_ERROR, openErr)
		return data.QuestionTable{}, openErr
	}
	defer db.Close()

	var qTable data.QuestionTable

	log.Print("Getting record from MySQL database")
	queryStr := "SELECT question, category, answer FROM trivia WHERE question_id = ?"
	queryErr := db.QueryRow(queryStr, questionID).Scan(&qTable.Question, &qTable.Category, &qTable.Answer)
	if queryErr != nil {
		if queryErr != sql.ErrNoRows {
			log.Print(MYSQL_DB_NAME_MSG+MYSQL_RESULTS_ERROR, queryErr.Error())
		}

		return data.QuestionTable{}, queryErr
	}

	return qTable, nil
}

// Update a single record in table
func (dbm *dbModel) Update(qRequest data.QuestionRequest) (int64, error) {
	db, openErr := dbm.Open(dbm.cfgData.ActiveDriver)
	if openErr != nil {
		log.Print(MYSQL_DB_NAME_MSG+MYSQL_OPEN_ERROR, openErr)
		return data.OPEN_ERROR_CODE, openErr
	}
	defer db.Close()

	log.Println("Updating record in MySQL database")
	queryStr := "UPDATE trivia SET question = ?, category = ?, answer = ? WHERE question_id = ?"
	sqlDB, execErr := db.Exec(queryStr, qRequest.Question, qRequest.Category, qRequest.Answer, qRequest.QuestionID)
	if execErr != nil {
		log.Print(MYSQL_DB_NAME_MSG+MYSQL_UPDATE_ERROR, execErr)
		return data.UPDATE_ERROR_CODE, execErr
	}

	rowsAffected, rowsAffectedErr := sqlDB.RowsAffected()
	if rowsAffectedErr != nil {
		log.Print(MYSQL_DB_NAME_MSG+MYSQL_ROWS_AFFECTED_ERROR, rowsAffectedErr.Error())
		return data.ROWS_AFFECTED_ERROR_CODE, rowsAffectedErr
	}

	return rowsAffected, nil
}

// Delete a single record from table
func (dbm *dbModel) Delete(questionID string) (int64, error) {
	db, openErr := dbm.Open(dbm.cfgData.ActiveDriver)
	if openErr != nil {
		log.Print(MYSQL_DB_NAME_MSG+MYSQL_OPEN_ERROR, openErr)
		return data.OPEN_ERROR_CODE, openErr
	}
	defer db.Close()

	log.Println("deleting record from MySQL database")
	queryStr := "DELETE FROM trivia WHERE question_id = ?"
	sqlDB, execErr := db.Exec(queryStr, questionID)
	if execErr != nil {
		log.Print(MYSQL_DB_NAME_MSG+MYSQL_DELETE_ERROR, execErr)
		return data.DELETE_ERROR_CODE, execErr
	}

	rowsAffected, rowsAffectedErr := sqlDB.RowsAffected()
	if rowsAffectedErr != nil {
		log.Print(MYSQL_DB_NAME_MSG+MYSQL_ROWS_AFFECTED_ERROR, rowsAffectedErr.Error())
		return data.ROWS_AFFECTED_ERROR_CODE, rowsAffectedErr
	}

	return rowsAffected, nil
}

func GetMySQLModel() *dbModel {
	if mySQLModel == nil {
		log.Print("Creating MySQL database model")

		mySQLModel = new(dbModel)
		cfg, getErr := config.Get()
		if getErr != nil {
			log.Print(MYSQL_DB_NAME_MSG+MYSQL_GET_CONFIG_ERROR, getErr)
			return nil
		}

		// Load config data
		var getCfgDataErr error
		mySQLModel.cfgData, getCfgDataErr = cfg.GetData()
		if getCfgDataErr != nil {
			log.Print(MYSQL_DB_NAME_MSG+MYSQL_GET_CONFIG_DATA_ERROR, getCfgDataErr)
			return nil
		}

	}

	return mySQLModel
}
