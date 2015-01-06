package main

import (
	flag "github.com/ogier/pflag"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"path/filepath"
)

type Config struct {
	Backend_timeout         int
	MaxKeepAliveConnections int
	Aerospike               struct {
		WriteTimeout      int
		ReadTimeout       int
		ConnectionTimeout int
		Hosts             []struct {
			Host string
			Port int
		}
	}
	Http struct {
		Port string
	}
	Memcache struct {
		Port string
	}
}

var config *Config
var config_filename string
var log_level string

func LoadConfigFromString(yaml_config []byte) *Config {
	var err error
	var config *Config

	if err = yaml.Unmarshal(yaml_config, &config); err != nil {
		panic(err)
	}
	return config
}

func LoadConfigFromFileName(config_filename string) *Config {
	var err error
	var yaml_config []byte
	filename, _ := filepath.Abs(config_filename)
	if yaml_config, err = ioutil.ReadFile(filename); err != nil {
		panic(err)
	}
	return LoadConfigFromString(yaml_config)
}

func ReadConfig() {
	config = LoadConfigFromFileName(config_filename)
}

func init() {
	flag.StringVar(&config_filename, "config", "config.yaml", "Config file")
	flag.StringVar(&log_level, "log-level", "INFO", "log level")
	flag.Parse()
}
