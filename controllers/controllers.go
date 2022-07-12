package controllers

import (
	"log"
	"sync"

	"github.com/sflewis2970/datastore-service/config"
	"github.com/sflewis2970/datastore-service/models"
)

type Controller struct {
	dbMutex   sync.Mutex
	dataModel *models.Model
	cfgData   *config.ConfigData
}

var controller *Controller

func New(args ...string) {
	if len(args) == 0 {
		log.Print("Adding an empty parameter...")
		args = append(args, "")
	}

	if controller == nil {
		log.Print("Creating controller object...")
		controller = new(Controller)

		// Load config data
		var cfgDataErr error
		controller.cfgData, cfgDataErr = config.Get().GetData(args[0])
		if cfgDataErr != nil {
			log.Print("Error getting config data: ", cfgDataErr)
			return
		}

		// Create dataModel
		controller.dataModel = models.New()
	}
}
