package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/sflewis2970/datastore-service/common"
)

const BASE_DIR_NAME string = "datastore-service"
const CONFIG_FILE_NAME string = "./config/config.json"

type GoCache struct {
	DriverName        string `json:"drivername"`
	DefaultExpiration int    `json:"expiration"`
	CleanupInterval   int    `json:"cleanup"`
}

type MySQL struct {
	DriverName string `json:"drivername"`
	Connection string `json:"connection"`
}

type PostGreSQL struct {
	DriverName string `json:"drivername"`
	Host       string `json:"host"`
	Port       int    `json:"port"`
	User       string `json:"user"`
}

type Message struct {
	CongratsMsg string `json:"congrats"`
	TryAgainMsg string `json:"tryagain"`
}

type ConfigData struct {
	HostName     string `json:"hostname"`
	HostPort     string `json:"hostport"`
	ActiveDriver string `json:"active"`
	GoCache      GoCache
	MySQL        MySQL
	PostGreSQL   PostGreSQL
	Messages     Message
}

type config struct {
	cfgData *ConfigData
}

var cfg *config

// Exported type functions
func (c *config) GetConfigData() *ConfigData {
	return c.cfgData
}

func (c *config) GetGoCacheData() *GoCache {
	return &c.cfgData.GoCache
}

func (c *config) GetMySQLData() *MySQL {
	return &c.cfgData.MySQL
}

func (c *config) GetPostGresData() *PostGreSQL {
	return &c.cfgData.PostGreSQL
}

// Unexported type functions
func (c *config) findBaseDir(currentDir string, targetDir string) int {
	level := 0
	dirs := strings.Split(currentDir, "\\")

	dirsSize := len(dirs)
	for idx := dirsSize - 1; idx >= 0; idx-- {
		if dirs[idx] == targetDir {
			break
		} else {
			level++
		}
	}

	return level
}

func (c *config) readConfigFile() error {
	// Get working directory
	wd, getErr := common.GetWorkingDir()
	if getErr != nil {
		log.Print("Error getting working directory")
		return getErr
	}

	// Find path
	levels := c.findBaseDir(wd, BASE_DIR_NAME)
	for levels > 0 {
		chErr := os.Chdir("..")
		if chErr != nil {
			log.Print("Error changind dir: ", chErr)
		}

		// Update levels
		levels--
	}

	// Read config file
	log.Print("reading config file...")
	data, readErr := ioutil.ReadFile(CONFIG_FILE_NAME)
	if readErr != nil {
		return readErr
	}

	c.cfgData = new(ConfigData)
	unmarshalErr := json.Unmarshal(data, c.cfgData)
	if unmarshalErr != nil {
		return unmarshalErr
	}

	return nil
}

// Exported package function
func GetConfig() (*config, error) {
	if cfg == nil {
		log.Print("creating config object")
		cfg = new(config)

		readErr := cfg.readConfigFile()
		if readErr != nil {
			log.Print("Error reading config file: ", readErr)
			return nil, readErr
		}
	}

	log.Print("returning config object")
	return cfg, nil
}
