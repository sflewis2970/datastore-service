package gocache

import (
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/sflewis2970/datastore-service/config"
	"github.com/sflewis2970/datastore-service/models/data"
)

var goCacheModel *dbModel

type dbModel struct {
	cfgData  *config.ConfigData
	memCache *cache.Cache
}

func (dbm *dbModel) Open(sqlDriverName string) (*sql.DB, error) {
	return nil, errors.New("Go-cache does not implement Open method")
}

// Ping database server, since this is local to the server make sure the object for storing data is created
func (dbm *dbModel) Ping() error {
	if dbm.memCache == nil {
		return errors.New("Go-cache has not been created")
	}

	return nil
}

// Insert a single record into table
func (dbm *dbModel) InsertQuestion(qRequest data.QuestionRequest) (int64, error) {
	var qt data.QuestionTable
	qt.Question = qRequest.Question
	qt.Category = qRequest.Category
	qt.Answer = qRequest.Answer

	log.Print("Adding a new record to map, ID: ", qRequest.QuestionID)
	dbm.memCache.Set(qRequest.QuestionID, qt, cache.DefaultExpiration)

	return data.RESULTS_DEFAULT, nil
}

// Get a single record from table
func (dbm *dbModel) GetQuestion(questionID string) (data.QuestionTable, error) {
	log.Print("Getting record from the map, with ID: ", questionID)

	item, itemFound := dbm.memCache.Get(questionID)

	var qt data.QuestionTable
	if itemFound {
		ok := false
		qt, ok = item.(data.QuestionTable)
		if !ok {
			log.Print("Error converting interface object: ", item)
		}
	} else {
		log.Print("No results returned")
	}

	return qt, nil
}

// Update a single record in table
func (dbm *dbModel) UpdateQuestion(qRequest data.QuestionRequest) (int64, error) {
	log.Println("Updating record in the map")

	var qt data.QuestionTable
	qt.Question = qRequest.Question
	qt.Category = qRequest.Category
	qt.Answer = qRequest.Answer

	dbm.memCache.Set(qRequest.QuestionID, qt, cache.DefaultExpiration)

	return data.RESULTS_DEFAULT, nil
}

// Delete a single record from table
func (dbm *dbModel) DeleteQuestion(questionID string) (int64, error) {
	log.Print("Deleting record with ID: ", questionID)

	// Delete the record from map
	dbm.memCache.Delete(questionID)

	return data.RESULTS_DEFAULT, nil
}

func GetGoCacheModel() *dbModel {
	if goCacheModel == nil {
		goCacheModel = new(dbModel)
		cfg, getErr := config.GetConfig()
		if getErr != nil {
			log.Fatal("Error getting config info: ", getErr)
		}

		goCacheModel.cfgData = cfg.GetConfigData()

		log.Print("Creating in-memory map to store data")
		goCacheModel.memCache = cache.New(time.Duration(goCacheModel.cfgData.GoCache.DefaultExpiration)*time.Minute, time.Duration(goCacheModel.cfgData.GoCache.CleanupInterval)*time.Minute)
	}

	return goCacheModel
}
