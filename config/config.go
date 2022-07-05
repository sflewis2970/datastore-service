package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/sflewis2970/datastore-service/common"
)

const BASE_DIR_NAME string = "datastore-service"
const CONFIG_FILE_NAME string = "./config/config.json"
const UPDATE_CONFIG_DATA string = "update"

// database drivers
const (
	GOCACHE_DRIVER    string = "gocache"
	GOREDIS_DRIVER    string = "goredis"
	MYSQL_DRIVER      string = "mysql"
	POSTGRESQL_DRIVER string = "postgres"
)

// Config variable keys
const (
	HOSTNAME string = "HOSTNAME"
	HOSTPORT string = "HOSTPORT"

	// The choices for activedriver are: "go-cache", "mysql", "postgres"
	ACTIVEDRIVER       string = "ACTIVEDRIVER"
	DEFAULT_EXPIRATION string = "expiration"
	CLEANUP_INTERVAL   string = "cleanup"
	GOREDIS_HOST       string = "GOREDIS_HOST"
	GOREDIS_PORT       string = "GOREDIS_PORT"
	GOREDIS_PASSWORD   string = "GOREDIS_PASSWORD"
	MYSQL_CONNECTION   string = "mysql_connection"
	POSTGRES_HOST      string = "postgres_host"
	POSTGRES_PORT      string = "postgres_port"
	POSTGRES_USER      string = "postgres_user"
)

type GoCache struct {
	DefaultExpiration int `json:"expiration"`
	CleanupInterval   int `json:"cleanup"`
}

type GoRedis struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Password string `json:"password"`
}

type MySQL struct {
	Connection string `json:"connection"`
}

type PostGreSQL struct {
	Host string `json:"host"`
	Port int    `json:"port"`
	User string `json:"user"`
}

type ConfigData struct {
	HostName     string `json:"hostname"`
	HostPort     string `json:"hostport"`
	ActiveDriver string `json:"active"`
	GoCache      GoCache
	GoRedis      GoRedis
	MySQL        MySQL
	PostGreSQL   PostGreSQL
}

type config struct {
	cfgData *ConfigData
}

var cfg *config

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

	unmarshalErr := json.Unmarshal(data, c.cfgData)
	if unmarshalErr != nil {
		return unmarshalErr
	}

	return nil
}

func (c *config) getConfigEnv() error {
	// Loading config environment variables
	log.Print("loading config environment variables...")

	// Update config data
	// Base config settings
	c.cfgData.HostName = os.Getenv(HOSTNAME)
	c.cfgData.HostPort = os.Getenv(HOSTPORT)
	c.cfgData.ActiveDriver = os.Getenv(ACTIVEDRIVER)

	switch c.cfgData.ActiveDriver {
	case GOCACHE_DRIVER:
		// Go-cache settings
		log.Print("Setting go-cache environment variables...")
		strVal := os.Getenv(DEFAULT_EXPIRATION)
		if len(strVal) > 0 {
			value, convErr := strconv.Atoi(strVal)
			if convErr != nil {
				log.Print("Error converting string to int...")
				return convErr
			}
			c.cfgData.GoCache.CleanupInterval = value
		}

		strVal = os.Getenv(DEFAULT_EXPIRATION)
		if len(strVal) > 0 {
			value, convErr := strconv.Atoi(strVal)
			if convErr != nil {
				log.Print("Error converting string to int...")
				return convErr
			}
			c.cfgData.GoCache.DefaultExpiration = value
		}

	case GOREDIS_DRIVER:
		// Go-redis settings
		log.Print("Setting go-redis environment variables...")
		c.cfgData.GoRedis.Host = os.Getenv(GOREDIS_HOST)
		c.cfgData.GoRedis.Port = os.Getenv(GOREDIS_PORT)
		c.cfgData.GoRedis.Password = os.Getenv(GOREDIS_PASSWORD)

	case MYSQL_DRIVER:
		// MySQL settings
		c.cfgData.MySQL.Connection = os.Getenv(MYSQL_CONNECTION)

	case POSTGRESQL_DRIVER:
		// PostGres settings
		log.Print("Setting postgres environment variables...")
		c.cfgData.PostGreSQL.Host = os.Getenv(POSTGRES_HOST)
		strVal := os.Getenv(POSTGRES_PORT)
		if len(strVal) > 0 {
			value, convErr := strconv.Atoi(strVal)
			if convErr != nil {
				log.Print("Error converting string to int...")
				return convErr
			}
			c.cfgData.PostGreSQL.Port = value
		}
		c.cfgData.PostGreSQL.User = os.Getenv(POSTGRES_USER)
	default:
		log.Print("Could not find supported driver...")
		log.Print("no database environment variables set...")
	}

	return nil
}

// Exported type functions
func (c *config) GetData(args ...string) (*ConfigData, error) {
	if len(args) > 0 {
		if args[0] == UPDATE_CONFIG_DATA {
			useCfgFile := os.Getenv("USECONFIGFILE")
			if len(useCfgFile) > 0 {
				log.Print("Using config file to load config")

				readErr := cfg.readConfigFile()
				if readErr != nil {
					log.Print("Error reading config file: ", readErr)
					return nil, readErr
				}
			} else {
				log.Print("Using config environment to load config")

				getErr := cfg.getConfigEnv()
				if getErr != nil {
					log.Print("Error getting config environment data: ", getErr)
					return nil, getErr
				}
			}
		}
	}

	return c.cfgData, nil
}

// Exported package function
func Get() *config {
	if cfg == nil {
		log.Print("creating config object")

		// Initialize config
		cfg = new(config)

		// Initialize config data
		cfg.cfgData = new(ConfigData)
	}

	log.Print("returning config object")
	return cfg
}
