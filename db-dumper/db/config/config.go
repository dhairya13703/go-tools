package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type DBConfig struct {
	DefaultHost     string `json:"default_host"`
	DefaultPort     string `json:"default_port"`
	DefaultUser     string `json:"default_user"`
	BackupDirectory string `json:"backup_directory"`
	Compression     bool   `json:"compression"`
}

func LoadConfig() (*DBConfig, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	configPath := filepath.Join(homeDir, ".db-dumper", "config.json")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return &DBConfig{
			DefaultHost:     "localhost",
			DefaultPort:     "3306",
			BackupDirectory: "backups",
			Compression:     true,
		}, nil
	}

	file, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var config DBConfig
	if err := json.Unmarshal(file, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
