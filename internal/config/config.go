package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	DsnDb         string `json:"db_dsn"`
	DriverDb      string `json:"db_driver"`
	Port          string `json:"port"`
	MigrationPath string `json:"migration_path"`
	InitDataPath  string `json:"init_data_path"`
}

func NewConfig() (Config, error) {
	var cfg Config

	fileJson, err := os.Open("config.json")
	if err != nil {
		return cfg, err
	}

	defer fileJson.Close()

	decoder := json.NewDecoder(fileJson)
	err = decoder.Decode(&cfg)

	if err != nil {
		return cfg, err
	}

	return cfg, nil
}
