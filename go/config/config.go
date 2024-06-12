package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type ConfigFile struct {
	Password string `yaml:"password"`
	ConfigID string `yaml:"config_id"`
}

func ReadConfigFile(filename string) (ConfigFile, error) {
	var config ConfigFile

	data, err := os.ReadFile(filename)
	if err != nil {
		return config, err
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return config, err
	}

	return config, nil
}
