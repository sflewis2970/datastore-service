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

func Initialize(args ...string) {
	if len(args) == 0 {
		log.Print("Adding an empty string...")
		args = append(args, "")
	}

	if ctrlr == nil {
		ctrlr = new(controller)

		// Get config object
		cfg, getCfgErr := config.Get()
		if getCfgErr != nil {
			log.Fatal("Error getting config: ", getCfgErr)
		}

		// Load config data
		var getCfgDataErr error
		ctrlr.cfgData, getCfgDataErr = cfg.GetData(args[0])
		if getCfgDataErr != nil {
			log.Fatal("Error getting config data: ", getCfgDataErr)
		}
	}
}
