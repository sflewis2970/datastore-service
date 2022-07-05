package goredis

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/sflewis2970/datastore-service/config"
	"github.com/sflewis2970/datastore-service/models/data"
)

const (
	GOREDIS_DB_NAME_MSG      string = "GO_REDIS: "
	GOREDIS_CREATE_CACHE_MSG string = "Creating in-memory map to store data..."
)

const (
	GOREDIS_GET_CONFIG_ERROR      string = "Getting config error...: "
	GOREDIS_GET_CONFIG_DATA_ERROR string = "Getting config data error...: "
	GOREDIS_OPEN_ERROR            string = "Open method not implemented..."
	GOREDIS_MARSHAL_ERROR         string = "Marshaling error...: "
	GOREDIS_UNMARSHAL_ERROR       string = "Unmarshaling error...: "
	GOREDIS_INSERT_ERROR          string = "Insert error...: "
	GOREDIS_ITEM_NOT_FOUND_ERROR  string = "Item not found...: "
	GOREDIS_GET_ERROR             string = "Get error...: "
	GOREDIS_UPDATE_ERROR          string = "Update error..."
	GOREDIS_DELETE_ERROR          string = "Delete error..."
	GOREDIS_RESULTS_ERROR         string = "Results error...: "
	GOREDIS_ROWS_AFFECTED_ERROR   string = "Rows affected error...: "
	GOREDIS_PING_ERROR            string = "Error pinging in-memory cache server...: "
	GOREDIS_CONVERSION_ERROR      string = "Conversion error...: "
)

var goRedisModel *dbModel

type dbModel struct {
	cfgData  *config.ConfigData
	memCache *redis.Client
}

func (dbm *dbModel) Open(sqlDriverName string) (*sql.DB, error) {
	return nil, errors.New(GOREDIS_DB_NAME_MSG + GOREDIS_OPEN_ERROR)
}

// Ping database server, since this is local to the server make sure the object for storing data is created
func (dbm *dbModel) Ping() error {
	ctx := context.Background()

	statusCmd := dbm.memCache.Ping(ctx)
	pingErr := statusCmd.Err()
	if pingErr != nil {
		log.Print(GOREDIS_DB_NAME_MSG+GOREDIS_PING_ERROR, pingErr)
		return pingErr
	}

	return nil
}

// Insert a single record into table
func (dbm *dbModel) Insert(qRequest data.QuestionRequest) (int64, error) {
	ctx := context.Background()

	var qt data.QuestionTable
	qt.Question = qRequest.Question
	qt.Category = qRequest.Category
	qt.Answer = qRequest.Answer

	byteStream, marshalErr := json.Marshal(qt)
	if marshalErr != nil {
		log.Print(GOREDIS_DB_NAME_MSG+GOREDIS_MARSHAL_ERROR, marshalErr)
		return data.RESULTS_DEFAULT, marshalErr
	}

	log.Print("Adding a new record to map, ID: ", qRequest.QuestionID)
	setErr := dbm.memCache.Set(ctx, qRequest.QuestionID, byteStream, time.Duration(0)).Err()
	if setErr != nil {
		log.Print(GOREDIS_DB_NAME_MSG+GOREDIS_INSERT_ERROR, setErr)
		return data.RESULTS_DEFAULT, setErr
	}

	return data.RESULTS_DEFAULT, nil
}

// Get a single record from table
func (dbm *dbModel) Get(questionID string) (data.QuestionTable, error) {
	log.Print("Getting record from the map, with ID: ", questionID)

	var qt data.QuestionTable
	ctx := context.Background()
	getResult, getErr := dbm.memCache.Get(ctx, questionID).Result()
	if getErr == redis.Nil {
		log.Print(GOREDIS_DB_NAME_MSG + GOREDIS_ITEM_NOT_FOUND_ERROR)
		return data.QuestionTable{}, nil
	} else if getErr != nil {
		log.Print(GOREDIS_DB_NAME_MSG+GOREDIS_GET_ERROR, getErr)
		return data.QuestionTable{}, getErr
	} else {
		unmarshalErr := json.Unmarshal([]byte(getResult), &qt)
		if unmarshalErr != nil {
			log.Print(GOREDIS_DB_NAME_MSG+GOREDIS_UNMARSHAL_ERROR, unmarshalErr)
			return data.QuestionTable{}, unmarshalErr
		}
	}

	return qt, nil
}

// Update a single record in table
func (dbm *dbModel) Update(qRequest data.QuestionRequest) (int64, error) {
	log.Println("Updating record in the map")

	ctx := context.Background()

	var qt data.QuestionTable
	qt.Question = qRequest.Question
	qt.Category = qRequest.Category
	qt.Answer = qRequest.Answer

	dbm.memCache.Set(ctx, qRequest.QuestionID, qt, 0)

	return data.RESULTS_DEFAULT, nil
}

// Delete a single record from table
func (dbm *dbModel) Delete(questionID string) (int64, error) {
	log.Print("Deleting record with ID: ", questionID)

	// Delete the record from map
	ctx := context.Background()
	delErr := dbm.memCache.Del(ctx, questionID).Err()
	if delErr != nil {

	}

	return data.RESULTS_DEFAULT, nil
}

func GetGoRedisModel() *dbModel {
	if goRedisModel == nil {
		// Initialize go-cache in-memory cache model
		goRedisModel = new(dbModel)

		// Load config data
		var getCfgDataErr error
		goRedisModel.cfgData, getCfgDataErr = config.Get().GetData()
		if getCfgDataErr != nil {
			log.Print(GOREDIS_DB_NAME_MSG+GOREDIS_GET_CONFIG_DATA_ERROR, getCfgDataErr)
			return nil
		}

		// Define go-redis cache settings
		log.Print(GOREDIS_DB_NAME_MSG + GOREDIS_CREATE_CACHE_MSG)
		addr := goRedisModel.cfgData.GoRedis.Host + goRedisModel.cfgData.GoRedis.Port
		log.Print("The goredis address is...: ", addr)

		// Create go-redis in-memory cache
		goRedisModel.memCache = redis.NewClient(&redis.Options{
			Addr:     addr, //GoRedis Server Address,
			Password: "",   // no password set
			DB:       0,    // use default DB
		})
	}

	return goRedisModel
}
