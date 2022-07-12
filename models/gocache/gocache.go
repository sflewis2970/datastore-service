package gocache

import (
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/sflewis2970/datastore-service/config"
	"github.com/sflewis2970/datastore-service/models/messages"
)

const (
	GOCACHE_DB_NAME_MSG      string = "GO_CACHE: "
	GOCACHE_CREATE_CACHE_MSG string = "Creating in-memory map to store data..."
)

const (
	GOCACHE_GET_CONFIG_ERROR      string = "Getting config error...: "
	GOCACHE_GET_CONFIG_DATA_ERROR string = "Getting config data error...: "
	GOCACHE_OPEN_ERROR            string = "Open method not implemented..."
	GOCACHE_INSERT_ERROR          string = "Insert error..."
	GOCACHE_GET_ERROR             string = "Get error..."
	GOCACHE_UPDATE_ERROR          string = "Update error..."
	GOCACHE_DELETE_ERROR          string = "Delete error..."
	GOCACHE_RESULTS_ERROR         string = "Results error...: "
	GOCACHE_ROWS_AFFECTED_ERROR   string = "Rows affected error...: "
	GOCACHE_PING_ERROR            string = "In-memory cache has not been created..."
	GOCACHE_CONVERSION_ERROR      string = "Conversion error...: "
)

var goCacheModel *dbModel

type dbModel struct {
	cfgData  *config.ConfigData
	memCache *cache.Cache
}

func (dbm *dbModel) Open(sqlDriverName string) (*sql.DB, error) {
	return nil, errors.New(GOCACHE_DB_NAME_MSG + GOCACHE_OPEN_ERROR)
}

// Ping database server, since this is local to the server make sure the object for storing data is created
func (dbm *dbModel) Ping() error {
	if dbm.memCache == nil {
		return errors.New(GOCACHE_DB_NAME_MSG + GOCACHE_PING_ERROR)
	}

	return nil
}

// Insert a single record into table
func (dbm *dbModel) Insert(qRequest messages.QuestionRequest) (int64, error) {
	var qt messages.QuestionTable
	qt.Question = qRequest.Question
	qt.Category = qRequest.Category
	qt.Answer = qRequest.Answer

	log.Print("Adding a new record to map, ID: ", qRequest.QuestionID)
	dbm.memCache.Set(qRequest.QuestionID, qt, cache.DefaultExpiration)

	return messages.RESULTS_DEFAULT, nil
}

// Get a single record from table
func (dbm *dbModel) Get(questionID string) (messages.QuestionTable, error) {
	log.Print("Getting record from the map, with ID: ", questionID)

	item, itemFound := dbm.memCache.Get(questionID)

	var qt messages.QuestionTable
	if itemFound {
		ok := false
		qt, ok = item.(messages.QuestionTable)
		if !ok {
			log.Print(GOCACHE_DB_NAME_MSG+GOCACHE_CONVERSION_ERROR, item)
		}
	} else {
		log.Print(messages.NO_RESULTS_RETURNED_MSG)
	}

	return qt, nil
}

// Update a single record in table
func (dbm *dbModel) Update(qRequest messages.QuestionRequest) (int64, error) {
	log.Println("Updating record in the map")

	var qt messages.QuestionTable
	qt.Question = qRequest.Question
	qt.Category = qRequest.Category
	qt.Answer = qRequest.Answer

	dbm.memCache.Set(qRequest.QuestionID, qt, cache.DefaultExpiration)

	return messages.RESULTS_DEFAULT, nil
}

// Delete a single record from table
func (dbm *dbModel) Delete(questionID string) (int64, error) {
	log.Print("Deleting record with ID: ", questionID)

	// Delete the record from map
	dbm.memCache.Delete(questionID)

	return messages.RESULTS_DEFAULT, nil
}

func GetGoCacheModel(cfgData *config.ConfigData) *dbModel {
	// Initialize go-cache in-memory cache model
	goCacheModel = new(dbModel)

	// Assign config data
	goCacheModel.cfgData = cfgData

	// Load config data
	var cfgDataErr error
	goCacheModel.cfgData, cfgDataErr = config.Get().GetData()
	if cfgDataErr != nil {
		log.Print(GOCACHE_DB_NAME_MSG+GOCACHE_GET_CONFIG_DATA_ERROR, cfgDataErr)
		return nil
	}

	log.Print(GOCACHE_DB_NAME_MSG + GOCACHE_CREATE_CACHE_MSG)
	goCacheModel.memCache = cache.New(time.Duration(goCacheModel.cfgData.GoCache.DefaultExpiration)*time.Minute, time.Duration(goCacheModel.cfgData.GoCache.CleanupInterval)*time.Minute)

	return goCacheModel
}
