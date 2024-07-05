package utils

import (
	"encoding/json"
	"os"
	"strings"
)

type Config struct {
	SetType string `json:"set_type"`
}

const configPath = "config.json"

var cfg *Config

func GetConfig() *Config {
	if cfg == nil {
		cfg = loadConfig()
	}
	return cfg
}

func InitializeConfig(setType string) {
	setType = strings.ToLower(setType)
	cfg = GetConfig()
	cfg.SetType = setType
	saveConfig(cfg)
}

func GetSetName() string {
	switch setType := GetConfig().SetType; setType {
	case "pokemon":
		return "Pokémon"
	case "yugioh":
		return "Yu-Gi-Oh!"
	default:
		return setType
	}
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

func loadConfig() *Config {
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
		SetType: "pokemon",
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
