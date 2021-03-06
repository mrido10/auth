package config

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Server struct {
		Host        string `yaml:"host"`
		ServicePort string `yaml:"servicePort"`
		Protocol    string `yaml:"protocol"`
	} `yaml:"server"`
	Jwt struct {
		Key string `yaml:"key"`
	} `yaml:"jwt"`
	Hash struct {
		Secret string `yaml:"secret"`
	} `yaml:"hash"`
	SendEmail struct {
		SMTP_HOST     string `yaml:"SMTP_HOST"`
		SMTP_PORT     int    `yaml:"SMTP_PORT"`
		SENDER_NAME   string `yaml:"SENDER_NAME"`
		AUTH_EMAIL    string `yaml:"AUTH_EMAIL"`
		AUTH_PASSWORD string `yaml:"AUTH_PASSWORD"`
	} `yaml:"sendEmail"`
	EncryptDecrypt struct {
		Key string `yaml:"key"`
	} `yaml:"encryptDecrypt"`
	Postgres struct {
		Address string `yaml:"address"`
	} `yaml:"postgres"`
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
