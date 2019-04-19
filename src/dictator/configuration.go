package main

import (
	"encoding/json"
	"io/ioutil"

	log "github.com/Sirupsen/logrus"
)

type NodeConfiguration struct {
	Name           string `json:"name"`
	Host           string `json:"host"`
	Port           int    `json:"port"`
	LoadingTimeout int    `json:"loading_timeout"`
}

type DictatorConfiguration struct {
	ServiceName   string            `json:"svc_name"`
	LogLevel      string            `json:"log_level"`
	Node          NodeConfiguration `json:"node"`
	HttpPort      int               `json:"http_port"`
	MasterService string            `json:"master_service"`
}

func NewDictatorConfiguration() DictatorConfiguration {
	log.Debug("Initialize configuration")

	return DictatorConfiguration{
		LogLevel: "INFO",
		HttpPort: 8000,
		Node: NodeConfiguration{
			Name:           "local",
			Host:           "localhost",
			Port:           6379,
			LoadingTimeout: 30,
		},
		MasterService: "",
	}
}

func (d *DictatorConfiguration) ReadConfigurationFile(configFilePath string) error {
	if configFilePath == "" {
		return nil
	}
	log.WithField("file", configFilePath).Debug("Reading configuration file")

	file, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		return err
	}

	err = json.Unmarshal(file, d)
	if err != nil {
		return err
	}

	return nil
}
