package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	TelegramToken  string `json:"telegram_token"`
	TelegramChatID int64  `json:"telegram_chat_id"`
	FilmsPath      string `json:"films_path"`
	ShowsPath      string `json:"shows_path"`
	APIPort        string `json:"api_port"`
}

func New() *Config {
	cfg, err := loadFile()
	if err != nil {
		panic(err)
	}

	return cfg
}

func loadFile() (*Config, error) {
	filepath := os.Getenv("CONFIG_FILE")
	if filepath == "" {
		return nil, fmt.Errorf("CONFIG_FILE environment variable is not set")
	}

	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	buf, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	cfg := &Config{}
	err = json.Unmarshal(buf, cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
