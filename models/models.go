package models

import (
	"log"
	"math"

	"github.com/sflewis2970/datastore-service/models/data"
	"github.com/sflewis2970/datastore-service/models/dsmysql"
	"github.com/sflewis2970/datastore-service/models/dspostgresql"
	"github.com/sflewis2970/datastore-service/models/gocache"
)

// Datastore contants
const (
	// DS_NOT_STARTED -- Datastore server has not been started or initialized
	DS_NOT_STARTED int = iota
	// DS_RUNNING -- Datastore server has been started and is ready for messages
	DS_RUNNING
	// DS_INVALID_SERVER_NAME -- When requesting the Datastore server status the wrong server name was provided
	DS_INVALID_SERVER_NAME
	// DS_UNAVAILABLE -- When requesting the Datastore server status the server never responded or the connect was refused
	DS_UNAVAILABLE int = math.MaxInt
)

func NewDBModel(activeDriver string) data.IDBModel {
	var dbModel data.IDBModel

	switch activeDriver {
	case "go-cache":
		dbModel = gocache.GetGoCacheModel()
	case "mysql":
		dbModel = dsmysql.GetMySQLModel()
	case "postgres":
		dbModel = dspostgresql.GetPostGreSQLModel()
	default:
		log.Print("Unsupported database driver, active driver: ", activeDriver)
	}

	return dbModel
}
