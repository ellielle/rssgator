package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

type Config struct {
	DBURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

const configFile = ".gatorconfig.json"

// local database URI to initialize with
const DBURL = "postgres://lily:snowmogchocobo@localhost:5432/gator?sslmode=disable"

// Read opens the config file and parses the content
// into a Config struct that is returned
func Read() (Config, error) {
	filePath, err := getConfigPath()
	if err != nil {
		return Config{}, err
	}
	data, err := os.ReadFile(filePath)
	if err != nil {
		createConfigFile()
		fmt.Println("creating configuration, re-run the program")
		os.Exit(1)
	}

	var cfg Config
	err = json.Unmarshal(data, &cfg)
	if err != nil {
		return Config{}, errors.New("error unmarshaling config")
	}
	return cfg, nil
}

// createConfigFile creates a config file with the default information
func createConfigFile() error {
	config := Config{
		DBURL: DBURL,
	}
	err := write(&config)
	if err != nil {
		return err
	}
	return nil
}

// getConfigPath returns the user's rssgator config file path
func getConfigPath() (string, error) {
	dir, err := os.UserHomeDir()
	if err != nil {
		return "", errors.New("could not find home directory")
	}
	return dir + "/" + configFile, nil
}

// SetUser sets the current user of the database and saves it to the config
func (cfg *Config) SetUser(user string) error {
	cfg.CurrentUserName = user
	err := write(cfg)
	if err != nil {
		return err
	}
	return nil
}

// write gets the correct file path for the config and writes the config to file
func write(cfg *Config) error {
	filePath, err := getConfigPath()
	if err != nil {
		return err
	}

	config, err := json.Marshal(cfg)
	if err != nil {
		return errors.New("could not marshal configuration")
	}
	// write the config to file
	err = os.WriteFile(filePath, config, 0600)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return nil
}
