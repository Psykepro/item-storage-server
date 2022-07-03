package config

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const (
	Path = "./config/config.yml"
)

type Config struct {
	RabbitMQ *RabbitMQ `yaml:"rabbit_mq"`
	Logging  *Logging  `yaml:"logging"`
}

type RabbitMQ struct {
	Host     string   `yaml:"host"`
	Port     string   `yaml:"port"`
	User     string   `yaml:"user"`
	Password string   `yaml:"password"`
	Queue    Queue    `yaml:"queue"`
	Consumer Consumer `yaml:"consumer"`
}

type Queue struct {
	Name       string         `yaml:"name"`
	Durable    bool           `yaml:"durable"`
	AutoDelete bool           `yaml:"auto_delete"`
	Exclusive  bool           `yaml:"exclusive"`
	NoWait     bool           `yaml:"no_wait"`
	Args       map[string]any `yaml:"args"`
}

type Consumer struct {
	Tag     string         `yaml:"tag"`
	AutoAck bool           `yaml:"auto_ack"`
	NoLocal bool           `yaml:"no_local"`
	Args    map[string]any `yaml:"args"`
}

type Logging struct {
	File   Logger `yaml:"file"`
	StdOut Logger `yaml:"stdout"`
}

type Logger struct {
	Level    string `yaml:"level"`
	Mode     string `yaml:"mode"`
	Encoding string `yaml:"encoding"`
}

func GetConfig(path string) (*Config, error) {
	var config Config
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, err
	}

	filename, err := filepath.Abs(path)
	if err != nil {
		log.Fatal(err)
	}

	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		log.Fatal(err)
	}

	return &config, err
}
