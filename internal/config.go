package internal

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	Db_url            string `json:"db_url"`
	Current_user_name string `json:"current_user_name"`
}

func NewConfig() Config {
	return Config{}
}

func getConfigFilePath() (string, error) {
	path, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home dir: %v", err)
	}
	return filepath.Join(path, configFileName), nil
}

func (cfg *Config) Read() error {
	path, err := getConfigFilePath()
	if err != nil {
		return err
	}
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open gatorconfig: %v", err)
	}
	defer file.Close()

	if err = json.NewDecoder(file).Decode(cfg); err != nil {
		return fmt.Errorf("failed to decode file to struct: %v", err)
	}
	return nil
}

func (cfg *Config) SetUser(username string) error {
	cfg.Current_user_name = username

	path, err := getConfigFilePath()
	if err != nil {
		return err
	}

	file, err := os.OpenFile(path,
		os.O_WRONLY|os.O_CREATE|os.O_TRUNC,
		0644)
	if err != nil {
		return fmt.Errorf("failed to open config file for write: %v", err)
	}
	defer file.Close()

	if err = json.NewEncoder(file).Encode(cfg); err != nil {
		return fmt.Errorf("failed to encode config: %v", err)
	}
	return nil
}
