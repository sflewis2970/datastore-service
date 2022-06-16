package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

const CONFIG_FILE_NAME = "./config/config.json"

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

func (c *config) GetConfigData() *ConfigData {
	return c.cfgData
}

func (c *config) readConfigFile() error {
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

func GetConfig() *config {
	if cfg == nil {
		log.Print("creating config object")
		cfg = new(config)

		readErr := cfg.readConfigFile()
		if readErr != nil {
			log.Print("Error reading config file: ", readErr)
		}
	}

	log.Print("returning config object")
	return cfg
}
