package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type Config struct {
	DBUrl           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

const CONFIG_FILE = ".gatorconfig.json"

func Read() Config {
	cfg := Config{}

	path, err := getConfigFilePath()
	if err != nil {
		log.Printf("Unable to get file path: %s", err.Error())
	}

	contentBytes, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Unable to read config file: %s", err.Error())
	}

	err = json.Unmarshal(contentBytes, &cfg)
	if err != nil {
		log.Fatalf("Unable to unmarshal config: %s", err.Error())
	}

	return cfg
}

func (cfg Config) SetUser(name string) {
	cfg.CurrentUserName = name

	contentBytes, err := json.Marshal(cfg)
	if err != nil {
		log.Fatalf("Unable to marshal config: %s", err.Error())
	}

	path, err := getConfigFilePath()
	if err != nil {
		log.Printf("Unable to get file path: %s", err.Error())
	}

  err = os.WriteFile(path, contentBytes, 066)
	if err != nil {
		log.Fatalf("Unable to write config: %s", err.Error())
	}
}

func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	path := fmt.Sprintf("%s/%s", homeDir, CONFIG_FILE)

	return path, nil
}
