package models

import (
	"log"
	"math"

	"github.com/sflewis2970/datastore-service/config"
	"github.com/sflewis2970/datastore-service/models/data"
	"github.com/sflewis2970/datastore-service/models/dsmysql"
	"github.com/sflewis2970/datastore-service/models/dspostgresql"
	"github.com/sflewis2970/datastore-service/models/gocache"
	"github.com/sflewis2970/datastore-service/models/goredis"
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
	case config.GOCACHE_DRIVER:
		dbModel = gocache.GetGoCacheModel()
	case config.GOREDIS_DRIVER:
		dbModel = goredis.GetGoRedisModel()
	case config.MYSQL_DRIVER:
		dbModel = dsmysql.GetMySQLModel()
	case config.POSTGRESQL_DRIVER:
		dbModel = dspostgresql.GetPostGreSQLModel()
	default:
		log.Print("Unsupported database driver, active driver: ", activeDriver)
	}

	return dbModel
}
