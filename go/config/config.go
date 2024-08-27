/*
 * Copyright (c) 2024 Intergreatme. All rights reserved.
 */

package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

const filename = "config.yaml"

type Configuration struct {
	CompanyID   string `yaml:"company_id"`
	PFXFilename string `yaml:"pfx"`
	Password    string `yaml:"password"`
	URL         string `yaml:"url"`
	CertDir     string `yaml:"cert_dir"`
}

func Read() (Configuration, error) {
	var config Configuration

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

func (cnf *Configuration) Write() error {
	out, err := yaml.Marshal(cnf)
	if err != nil {
		return err
	}
	err = os.WriteFile(filename, out, 0664)
	if err != nil {
		return err
	}
	return nil
}
