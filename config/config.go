package config

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	MySql struct {
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		DbName   string `yaml:"dbName"`
	} `yaml:"mySql"`
	Server struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
	} `yaml:"server"`
	Jwt struct {
		Key string `yaml:"key"`
	} `yaml:"jwt"`
	Hash struct {
		Secret string `yaml:"secret"`
	} `yaml:"hash"`
}

func GetConfig() (*Config, error) {
	buf, err := ioutil.ReadFile("config/config.yaml")
	if err != nil {
		return nil, err
	}

	c := &Config{}
	err = yaml.Unmarshal(buf, c)
	if err != nil {
		return nil, fmt.Errorf("in file %q: %v", "config.yaml", err)
	}

	return c, nil
}
