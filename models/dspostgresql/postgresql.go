package dspostgresql

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/sflewis2970/datastore-service/config"
	"github.com/sflewis2970/datastore-service/models/data"
)

const (
	POSTGRESQL_DB_NAME_MSG string = "POSTGRESQL: "
)

const (
	POSTGRESQL_GET_CONFIG_ERROR      string = "Getting config error...: "
	POSTGRESQL_GET_CONFIG_DATA_ERROR string = "Getting config data error...: "
	POSTGRESQL_OPEN_ERROR            string = "Error opening database..."
	POSTGRESQL_INSERT_ERROR          string = "Error inserting record..."
	POSTGRESQL_GET_ERROR             string = "Error getting record..."
	POSTGRESQL_UPDATE_ERROR          string = "Error updating record..."
	POSTGRESQL_DELETE_ERROR          string = "Error deleting record..."
	POSTGRESQL_RESULTS_ERROR         string = "Error getting results...: "
	POSTGRESQL_ROWS_AFFECTED_ERROR   string = "Error getting rows affected...: "
	POSTGRESQL_PING_ERROR            string = "Error pinging database server..."
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
		log.Print(POSTGRESQL_DB_NAME_MSG+POSTGRESQL_OPEN_ERROR, openErr.Error())
		return nil, openErr
	}

	return db, nil
}

// Ping database server by verifying the database connection is active
func (dbm *dbModel) Ping() error {
	db, openErr := dbm.Open(dbm.cfgData.ActiveDriver)
	if openErr != nil {
		log.Print(POSTGRESQL_DB_NAME_MSG+POSTGRESQL_OPEN_ERROR, openErr.Error())
		return openErr
	}
	defer db.Close()

	pingErr := db.Ping()
	if openErr != nil {
		log.Print(POSTGRESQL_DB_NAME_MSG + POSTGRESQL_PING_ERROR)
		return pingErr
	}

	return nil
}

// Insert a single record into table
func (dbm *dbModel) Insert(qRequest data.QuestionRequest) (int64, error) {
	db, openErr := dbm.Open(dbm.cfgData.ActiveDriver)
	if openErr != nil {
		log.Print(POSTGRESQL_DB_NAME_MSG+POSTGRESQL_OPEN_ERROR, openErr.Error())
		return data.RESULTS_DEFAULT, openErr
	}
	defer db.Close()

	log.Print("Adding a new record to the database")
	queryStr := "insert into trivia VALUES ($1, $2, $3, $4);"
	sqlDB, execErr := db.Exec(queryStr, qRequest.QuestionID, qRequest.Question, qRequest.Category, qRequest.Answer)
	if execErr != nil {
		log.Print(POSTGRESQL_DB_NAME_MSG + POSTGRESQL_INSERT_ERROR)
		return data.RESULTS_DEFAULT, execErr
	}

	rowsAffected, rowsAffectedErr := sqlDB.RowsAffected()
	if rowsAffectedErr != nil {
		log.Print(POSTGRESQL_DB_NAME_MSG+POSTGRESQL_ROWS_AFFECTED_ERROR, rowsAffectedErr.Error())
		return data.RESULTS_DEFAULT, rowsAffectedErr
	}

	return rowsAffected, nil
}

// Get a single record from table
func (dbm *dbModel) Get(questionID string) (data.QuestionTable, error) {
	db, openErr := dbm.Open(dbm.cfgData.ActiveDriver)
	if openErr != nil {
		log.Print(POSTGRESQL_DB_NAME_MSG+POSTGRESQL_OPEN_ERROR, openErr.Error())
		return data.QuestionTable{}, openErr
	}
	defer db.Close()

	var qTable data.QuestionTable

	log.Print("Getting a single record from the database")
	queryStr := "SELECT question, category, answer FROM trivia WHERE question_id = $1;"
	scanErr := db.QueryRow(queryStr, questionID).Scan(&qTable.Question, &qTable.Category, &qTable.Answer)
	if scanErr != nil && scanErr != sql.ErrNoRows {
		log.Print(POSTGRESQL_DB_NAME_MSG+POSTGRESQL_GET_ERROR, scanErr.Error())
		return data.QuestionTable{}, scanErr
	}

	return qTable, nil
}

// Update a single record in table
func (dbm *dbModel) Update(qRequest data.QuestionRequest) (int64, error) {
	db, openErr := dbm.Open(dbm.cfgData.ActiveDriver)
	if openErr != nil {
		log.Print(POSTGRESQL_DB_NAME_MSG+POSTGRESQL_OPEN_ERROR, openErr.Error())
		return data.RESULTS_DEFAULT, openErr
	}
	defer db.Close()

	log.Println("Updating a single record in the database")
	queryStr := "UPDATE trivia SET question = $2, category = $3, answer = $4 WHERE question_id = $1"
	sqlDB, execErr := db.Exec(queryStr, qRequest.QuestionID, qRequest.Question, qRequest.Category, qRequest.Answer)
	if execErr != nil {
		log.Print(POSTGRESQL_DB_NAME_MSG+POSTGRESQL_UPDATE_ERROR, execErr.Error())
		return data.RESULTS_DEFAULT, execErr
	}

	rowsAffected, rowsAffectedErr := sqlDB.RowsAffected()
	if rowsAffectedErr != nil {
		log.Print(POSTGRESQL_DB_NAME_MSG+POSTGRESQL_ROWS_AFFECTED_ERROR, rowsAffectedErr.Error())
		return data.RESULTS_DEFAULT, nil
	}

	return rowsAffected, nil
}

// Delete a single record from table
func (dbm *dbModel) Delete(questionID string) (int64, error) {
	db, openErr := dbm.Open(dbm.cfgData.ActiveDriver)
	if openErr != nil {
		log.Print(POSTGRESQL_DB_NAME_MSG+POSTGRESQL_OPEN_ERROR, openErr.Error())
		return data.RESULTS_DEFAULT, openErr
	}
	defer db.Close()

	log.Println("deleting a single record from the database")
	queryStr := "DELETE FROM trivia WHERE question_id = $1"
	sqlDB, execErr := db.Exec(queryStr, questionID)
	if execErr != nil {
		log.Print(POSTGRESQL_DB_NAME_MSG+POSTGRESQL_DELETE_ERROR, execErr.Error())
		return data.RESULTS_DEFAULT, execErr
	}

	rowsAffected, rowsAffectedErr := sqlDB.RowsAffected()
	if rowsAffectedErr != nil {
		log.Print(POSTGRESQL_DB_NAME_MSG+POSTGRESQL_ROWS_AFFECTED_ERROR, rowsAffectedErr.Error())
		return data.RESULTS_DEFAULT, nil
	}

	return rowsAffected, nil
}

func GetPostGreSQLModel() *dbModel {
	if postgreSQLModel == nil {
		log.Print("Creating PostgreSQL database model")

		postgreSQLModel = new(dbModel)
		cfg, getCfgErr := config.Get()
		if getCfgErr != nil {
			log.Print(POSTGRESQL_DB_NAME_MSG+POSTGRESQL_GET_CONFIG_ERROR, getCfgErr)
			return nil
		}

		// Load config data
		var getCfgDataErr error
		postgreSQLModel.cfgData, getCfgDataErr = cfg.GetData()
		if getCfgDataErr != nil {
			log.Print(POSTGRESQL_DB_NAME_MSG+POSTGRESQL_GET_CONFIG_DATA_ERROR, getCfgDataErr)
			return nil
		}

	}

	return postgreSQLModel
}
