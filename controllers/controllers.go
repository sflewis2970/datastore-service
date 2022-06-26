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

		// Load config data
		var getCfgDataErr error
		ctrlr.cfgData, getCfgDataErr = config.Get().GetData(args[0])
		if getCfgDataErr != nil {
			log.Print("Error getting config data: ", getCfgDataErr)
			return
		}
	}
}
