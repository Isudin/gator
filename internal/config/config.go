package config

import (
	"encoding/json"
	"fmt"
	"os"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	DbUrl       string `json:"db_url"`
	CurrentUser string `json:"current_user_name"`
}

func (c *Config) SetUser(username string) error {
	if username == "" {
		return fmt.Errorf("username cannot be empty")
	}

	c.CurrentUser = username

	dir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	bytes, err := json.Marshal(c)
	if err != nil {
		return err
	}

	err = os.WriteFile(dir+"/"+configFileName, bytes, 0770)
	if err != nil {
		return err
	}

	return nil
}

func Read() (Config, error) {
	dir, err := os.UserHomeDir()
	if err != nil {
		return Config{}, err
	}

	file, err := os.ReadFile(dir + "/" + configFileName)
	if err != nil {
		return Config{}, err
	}

	var config Config
	err = json.Unmarshal(file, &config)
	return config, err
}
