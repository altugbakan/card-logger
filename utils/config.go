package utils

import (
	"encoding/json"
	"os"
)

type Config struct {
	Type string `json:"type"`
}

const configPath = "config.json"

var cfg *Config

func GetConfig() *Config {
	if cfg == nil {
		cfg = initializeConfig()
	}
	return cfg
}

func saveConfig(config *Config) {
	file, err := os.Create(configPath)
	if err != nil {
		LogError("Error creating JSON config file: %v", err)
	}

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(config); err != nil {
		LogError("Error encoding JSON config file: %v", err)
	}
}

func initializeConfig() *Config {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		config := getDefaultConfig()
		saveConfig(config)
		return config
	} else {
		return loadConfigFromJSON()
	}
}

func getDefaultConfig() *Config {
	return &Config{
		Type: "pokemon",
	}
}

func loadConfigFromJSON() *Config {
	file, err := os.ReadFile(configPath)
	if err != nil {
		LogError("Error reading JSON config file: %v", err)
	}

	var config Config
	if err := json.Unmarshal(file, &config); err != nil {
		LogError("Error unmarshalling JSON config file: %v", err)
	}

	return &config
}
