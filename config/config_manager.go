package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// LoadConfig loads configuration from a given filepath and unmarshals it into the provided config structure.
func Load_Config(pkg string, label string, config interface{}) error {
	path := fmt.Sprintf("./%s/%s.json", pkg, label)
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, config)
}

// SaveConfig saves the configuration (provided as an interface{}) to a given filepath.
func Save_Config(pkg string, label string, config interface{}) error {
	path := fmt.Sprintf("./%s/%s.json", pkg, label)
	data, err := json.MarshalIndent(config, "", "    ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}
