package main

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type ConfigError struct {
	ErrorMsg string
}

func (e ConfigError) Error() string {
	return fmt.Sprint("ConfigFile: ", e.ErrorMsg)
}

type Config struct {
	ApiKey       string
	RemoteIp     string
	Logfile      string
	Logtype      string
	PollInterval int
}

func NewConfig() Config {
	return Config{
		ApiKey:       "",
		RemoteIp:     "localhost",
		Logfile:      "./log.csv",
		Logtype:      "csv",
		PollInterval: 60,
	}
}

func (c *Config) Read(path string) error {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return ConfigError{err.Error()}
	}

	var newConf Config
	err = yaml.Unmarshal(buf, &newConf)
	if err == nil {
		*c = newConf
	} else {
		return ConfigError{err.Error()}
	}

	return nil
}

func (c *Config) Write(path string) error {
	buf, err := yaml.Marshal(*c)
	if err != nil {
		return ConfigError{fmt.Sprint("Error saving config: ", err)}
	}

	err = ioutil.WriteFile(path, buf, 0600)
	if err != nil {
		return ConfigError{fmt.Sprint("Error writing config file: ", err)}
	}

	return nil
}
