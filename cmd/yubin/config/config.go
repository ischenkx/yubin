package config

import (
	"encoding/json"
	"errors"
	"gopkg.in/yaml.v2"
	"os"
	"strings"
)

type Config struct {
	DB struct {
		Postgres struct {
			Host     string `json:"host,omitempty" yaml:"host"`
			Port     int    `json:"port,omitempty" yaml:"port"`
			DB       string `json:"db,omitempty" yaml:"db"`
			User     string `json:"user,omitempty" yaml:"user"`
			Password string `json:"password,omitempty" yaml:"password"`
		} `json:"postgres" yaml:"postgres"`
	} `json:"db" yaml:"db"`

	ViewStat struct {
		Redis struct {
			Channel  string `json:"channel,omitempty" yaml:"channel"`
			Addr     string `json:"addr,omitempty" yaml:"addr"`
			Username string `json:"username,omitempty" yaml:"username"`
			Password string `json:"password,omitempty" yaml:"password"`
			DB       int    `json:"db,omitempty" yaml:"db"`
		} `json:"redis" yaml:"redis"`

		URLs []string `json:"urls,omitempty" yaml:"urls"`
	} `json:"view_stat" yaml:"view_stat"`

	API struct {
		HTTP struct {
			Host string `json:"host,omitempty" yaml:"host"`
			Port int    `json:"port,omitempty" yaml:"port"`
		} `json:"http" yaml:"http"`
	} `json:"api" yaml:"api"`
}

func Read(path string) (Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return Config{}, err
	}
	defer file.Close()

	var config Config

	switch {
	case strings.HasSuffix(path, ".yml"), strings.HasSuffix(path, ".yaml"):
		if err := yaml.NewDecoder(file).Decode(&config); err != nil {
			return Config{}, err
		}
	case strings.HasSuffix(path, ".json"):
		if err := json.NewDecoder(file).Decode(&config); err != nil {
			return Config{}, err
		}
	default:
		return Config{}, errors.New("unsupported format")
	}
	return config, nil
}
