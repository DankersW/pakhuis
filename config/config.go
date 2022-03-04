package config

import (
	"fmt"
	"io/ioutil"

	"github.com/DankersW/pakhuis/models"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

const (
	CONFIG_FILE_PATH = "config.yml"
)

func Get() models.Config {
	config, err := parseYamlFile(CONFIG_FILE_PATH)
	if err != nil {
		log.Fatal("Failed to parse config file '%s', %s", CONFIG_FILE_PATH, err.Error())
		return models.Config{}
	}
	return config
}

func parseYamlFile(file string) (models.Config, error) {
	buffer, err := ioutil.ReadFile(file)
	if err != nil {
		return models.Config{}, err
	}
	config := models.Config{}
	if yaml.Unmarshal(buffer, &config) != nil {
		return models.Config{}, fmt.Errorf("could not unmarshal, %s", err.Error())
	}
	return config, nil
}
