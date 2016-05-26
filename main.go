/**Copyright Blue Medora Inc. 2016**/

package main

import (
	"github.com/BlueMedora/bluemedora-firehose-nozzle/bluemedorafirehosenozzle"
	"github.com/BlueMedora/bluemedora-firehose-nozzle/logger"
	"github.com/BlueMedora/bluemedora-firehose-nozzle/nozzleconfiguration"
	"github.com/BlueMedora/bluemedora-firehose-nozzle/webserver"
)

const (
	defaultConfigLocation = "./config/bluemedora-firehose-nozzle.json"

	defaultLogDirectory = "./logs"
	nozzleLogFile       = "bm_nozzle.log"
	nozzleLogName       = "bm_firehose_nozzle"

	webserverLogFile = "bm_server.log"
	webserverLogName = "bm_server"
)

func main() {
	logger := logger.New(defaultLogDirectory, nozzleLogFile, nozzleLogName)
	logger.Debug("working log")

	//Read in config
	config, err := nozzleconfiguration.New(defaultConfigLocation, logger)
	if err != nil {
		logger.Fatalf("Error parsing config file: %s", err.Error())
	}

	server := createWebServer(config)

	nozzle := bluemedorafirehosenozzle.New(config, server, logger)
	nozzle.Start()

	if err != nil {
		logger.Fatalf("Error while running nozzle: %s", err.Error())
	}

	//Start nozzle
}

func createWebServer(config *nozzleconfiguration.NozzleConfiguration) *webserver.WebServer {
	logger := logger.New(defaultLogDirectory, webserverLogFile, webserverLogName)
	return webserver.New(config, logger)
}
