package configurator

import (
	"encoding/json"
	"os"
)

type DiscordConfig struct {
	ClientID      string `json:"clientId"`
	Token         string `json:"token"`
	UpdateSeconds int64  `json:"updateSeconds"`
}

type ProjectConfig struct {
	Debug   bool           `json:"debug"`
	Discord *DiscordConfig `json:"discord"`
}

func LoadConfig(file string) *ProjectConfig {
	byteValue, err := os.ReadFile(file)
	if err != nil {
		return &ProjectConfig{Debug: false, Discord: nil}
	}
	var config ProjectConfig //nolint:wsl
	_ = json.Unmarshal(byteValue, &config)

	return &config
}
