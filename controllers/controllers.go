package controllers

import (
	"log"
	"sync"

	"github.com/sflewis2970/datastore-service/config"
)

type controller struct {
	dbMutex sync.Mutex
	cfgData *config.ConfigData
}

var ctrlr *controller

func InitializeController() {
	if ctrlr == nil {
		ctrlr = new(controller)

		// Get config object
		cfg, getErr := config.GetConfig()
		if getErr != nil {
			log.Fatal("Error getting config")
		}

		// Load config data into property
		ctrlr.cfgData = cfg.GetConfigData()
	}
}
