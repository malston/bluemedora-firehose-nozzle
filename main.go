/**Copyright Blue Medora Inc. 2016**/

package main

import (
    "github.com/BlueMedora/bluemedora-firehose-nozzle/logger"
    "github.com/BlueMedora/bluemedora-firehose-nozzle/nozzleconfiguration"
    "github.com/BlueMedora/bluemedora-firehose-nozzle/bluemedorafirehosenozzle"
)

var (
    defaultConfigLocation = "./config/bluemedora-firehose-nozzle.json"
    defaultLogDirector = "./logs"
)

func main() {
    logger := logger.New(defaultLogDirector)
    logger.Debug("working log")
    
    //Read in config
    config, err := nozzleconfiguration.New(defaultConfigLocation, logger)
    if err != nil {
        logger.Fatalf("Error parsing config file: %s", err.Error())
    }
    
    nozzle := bluemedorafirehosenozzle.New(config, logger)
    nozzle.Start()
    
    if err != nil {
        logger.Fatalf("Error while running nozzle: %s", err.Error())
    }
    
    //Start nozzle
}