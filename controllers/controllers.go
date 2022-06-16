package controllers

import (
	"sync"

	"github.com/sflewis2970/datastore-service/config"
)

var dbMutex sync.Mutex
var cfgData *config.ConfigData

func LoadConfig() {
	cfgData = config.GetConfig().GetConfigData()
}
