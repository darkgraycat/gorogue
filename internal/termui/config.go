package termui

import (
	"encoding/json"
	"os"
)

type Characters struct {
	Cursor  string `json:"cursor"`
	Borders string `json:"borders"`
}

type Colors struct {
	Primary   byte `json:"primary"`
	Secondary byte `json:"secondary"`
}

type Config struct {
	Characters Characters `json:"characters"`
	Colors     Colors     `json:"colors"`
}

func LoadConfig(path string) (Config, error) {
	var cfg Config
	data, err := os.ReadFile(path)
	if err != nil {
		return cfg, err
	}
	err = json.Unmarshal(data, &cfg)
	return cfg, err
}
