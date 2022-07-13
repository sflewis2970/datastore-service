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
	"github.com/sflewis2970/datastore-service/models/messages"
)

const (
	REDIS_DB_NAME_MSG      string = "GO_REDIS: "
	REDIS_CREATE_CACHE_MSG string = "Creating in-memory map to store data..."
)

const (
	REDIS_GET_CONFIG_ERROR      string = "Getting config error...: "
	REDIS_GET_CONFIG_DATA_ERROR string = "Getting config data error...: "
	REDIS_OPEN_ERROR            string = "Open method not implemented..."
	REDIS_MARSHAL_ERROR         string = "Marshaling error...: "
	REDIS_UNMARSHAL_ERROR       string = "Unmarshaling error...: "
	REDIS_INSERT_ERROR          string = "Insert error...: "
	REDIS_ITEM_NOT_FOUND_ERROR  string = "Item not found...: "
	REDIS_GET_ERROR             string = "Get error...: "
	REDIS_UPDATE_ERROR          string = "Update error..."
	REDIS_DELETE_ERROR          string = "Delete error...: "
	REDIS_RESULTS_ERROR         string = "Results error...: "
	REDIS_ROWS_AFFECTED_ERROR   string = "Rows affected error...: "
	REDIS_PING_ERROR            string = "Error pinging in-memory cache server...: "
	REDIS_CONVERSION_ERROR      string = "Conversion error...: "
)

var redisModel *dbModel

type dbModel struct {
	cfgData  *config.ConfigData
	memCache *redis.Client
}

func (dbm *dbModel) Open(sqlDriverName string) (*sql.DB, error) {
	return nil, errors.New(REDIS_DB_NAME_MSG + REDIS_OPEN_ERROR)
}

// Ping database server, since this is local to the server make sure the object for storing data is created
func (dbm *dbModel) Ping() error {
	ctx := context.Background()

	statusCmd := dbm.memCache.Ping(ctx)
	pingErr := statusCmd.Err()
	if pingErr != nil {
		log.Print(REDIS_DB_NAME_MSG+REDIS_PING_ERROR, pingErr)
		return pingErr
	}

	return nil
}

// Insert a single record into table
func (dbm *dbModel) Insert(qRequest messages.QuestionRequest) (int64, error) {
	ctx := context.Background()

	var qt messages.QuestionTable
	qt.Question = qRequest.Question
	qt.Category = qRequest.Category
	qt.Answer = qRequest.Answer

	byteStream, marshalErr := json.Marshal(qt)
	if marshalErr != nil {
		log.Print(REDIS_DB_NAME_MSG+REDIS_MARSHAL_ERROR, marshalErr)
		return messages.RESULTS_DEFAULT, marshalErr
	}

	log.Print("Adding a new record to map, ID: ", qRequest.QuestionID)
	setErr := dbm.memCache.Set(ctx, qRequest.QuestionID, byteStream, time.Duration(0)).Err()
	if setErr != nil {
		log.Print(REDIS_DB_NAME_MSG+REDIS_INSERT_ERROR, setErr)
		return messages.RESULTS_DEFAULT, setErr
	}

	return messages.RESULTS_DEFAULT, nil
}

// Get a single record from table
func (dbm *dbModel) Get(questionID string) (messages.QuestionTable, error) {
	log.Print("Getting record from the map, with ID: ", questionID)

	var qt messages.QuestionTable
	ctx := context.Background()
	getResult, getErr := dbm.memCache.Get(ctx, questionID).Result()
	if getErr == redis.Nil {
		log.Print(REDIS_DB_NAME_MSG + REDIS_ITEM_NOT_FOUND_ERROR)
		return messages.QuestionTable{}, nil
	} else if getErr != nil {
		log.Print(REDIS_DB_NAME_MSG+REDIS_GET_ERROR, getErr)
		return messages.QuestionTable{}, getErr
	} else {
		unmarshalErr := json.Unmarshal([]byte(getResult), &qt)
		if unmarshalErr != nil {
			log.Print(REDIS_DB_NAME_MSG+REDIS_UNMARSHAL_ERROR, unmarshalErr)
			return messages.QuestionTable{}, unmarshalErr
		}
	}

	return qt, nil
}

// Update a single record in table
func (dbm *dbModel) Update(qRequest messages.QuestionRequest) (int64, error) {
	log.Println("Updating record in the map")

	ctx := context.Background()

	var qt messages.QuestionTable
	qt.Question = qRequest.Question
	qt.Category = qRequest.Category
	qt.Answer = qRequest.Answer

	dbm.memCache.Set(ctx, qRequest.QuestionID, qt, 0)

	return messages.RESULTS_DEFAULT, nil
}

// Delete a single record from table
func (dbm *dbModel) Delete(questionID string) (int64, error) {
	log.Print("Deleting record with ID: ", questionID)

	// Delete the record from map
	ctx := context.Background()
	delErr := dbm.memCache.Del(ctx, questionID).Err()
	if delErr != nil {
		log.Print(REDIS_DB_NAME_MSG+REDIS_DELETE_ERROR, delErr)
	}

	return messages.RESULTS_DEFAULT, nil
}

func GetRedisModel(cfgData *config.ConfigData) *dbModel {
	// Initialize go-cache in-memory cache model
	log.Print("Creating goRedis dbModel object...")
	redisModel = new(dbModel)

	// Assign config data
	redisModel.cfgData = cfgData

	// Define go-redis cache settings
	log.Print(REDIS_DB_NAME_MSG + REDIS_CREATE_CACHE_MSG)

	// Define connection variables
	var redisOptions *redis.Options

	// The config package handles reading the environment variables and parsing the url.
	// Once the external packages access the values, the environment has lready been taken
	// care of.
	addr := redisModel.cfgData.Redis.URL + redisModel.cfgData.Redis.Port
	log.Print("The redis address is...: ", addr)

	redisOptions = &redis.Options{
		Addr:     addr, // redis Server Address,
		Password: "",   // no password set
		DB:       0,    // use default DB
	}

	// Create go-redis in-memory cache
	redisModel.memCache = redis.NewClient(redisOptions)

	return redisModel
}
