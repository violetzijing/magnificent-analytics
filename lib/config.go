package lib

import (
	"encoding/json"
	"fmt"
	"os"
)

// Config defines config content and type
type MagnificentConfig struct {
	URL          string  `json:"url"`
	Threshold    int     `json:"threshold"`
	Timeout      int     `json:"timeout"`
	IntervalTime int     `json:"interval_time"`
	HealthRatio  float64 `json:"health_ratio"`
}

// ParseConfig returns parsed config from config file
func ParseConfig() *MagnificentConfig {
	file, err := os.Open("config/development/config.json")
	if err != nil {
		panic(fmt.Sprintf("failed to get config file, err: %s", err.Error()))
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	cfg := &MagnificentConfig{}
	if err := decoder.Decode(cfg); err != nil {
		panic(fmt.Sprintf("failed to parse config, err: %s", err.Error()))
	}
	return cfg
}
