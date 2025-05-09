package config

import (
	"fmt"
	"regexp"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Server   Server   `yaml:"server"`
	Log      Log      `yaml:"log"`
	Loki     Loki     `yaml:"loki"`
	Messages Messages `yaml:"messages"`
	Storage  Storage  `yaml:"storage"`
}

type Server struct {
	UDP UDP `yaml:"udp"`
	TCP TCP `yaml:"tcp"`
}

type UDP struct {
	Enabled bool   `yaml:"enabled" env:"UDP_ENABLED" envDefault:"false"`
	IP      string `yaml:"ip" env:"UDP_IP" envDefault:"0.0.0.0"`
	Port    int    `yaml:"port" env:"UDP_PORT" envDefault:"1514"`
}

type TCP struct {
	Enabled bool   `yaml:"enabled" env:"TCP_ENABLED" envDefault:"false"`
	IP      string `yaml:"ip" env:"TCP_IP" envDefault:"0.0.0.0"`
	Port    int    `yaml:"port" env:"TCP_PORT" envDefault:"1514"`
}

type Log struct {
	Level      string `yaml:"level" env:"LOG_LEVEL" envDefault:"info"`
	ShowSource bool   `yaml:"show_source" env:"LOG_SHOW_SOURCE" envDefault:"false"`
	Format     string `yaml:"format" env:"LOG_FORMAT" envDefault:"json"`
}

type Loki struct {
	Server  string `yaml:"server" env:"LOKI_SERVER" envDefault:"http://127.0.0.1:3100"`
	Timeout string `yaml:"timeout" env:"LOKI_TIMEOUT" envDefault:"10s"`
}

type Storage struct {
	LogLifetime   string `yaml:"log_lifetime" env:"LOG_LIFETIME" envDefault:"14d"`
	PathToStorage string `yaml:"path_to_storage" env:"PATH_TO_STORAGE" envDefault:"./aglog.db"`
}

type Messages struct {
	Formats []string `yaml:"formats"`
}

func ReadConfig(path string) Config {
	var conf Config

	err := cleanenv.ReadConfig(path, &conf)
	if err != nil {
		panic(err)
	}

	checkFormats(conf.Messages.Formats)

	return conf
}

func checkFormats(formats []string) {
	for _, format := range formats {
		_, err := regexp.Compile(format)
		if err != nil {
			panic(fmt.Errorf("string %s is not valid format: %w", format, err))
		}
	}
}
