/**Copyright Blue Medora Inc. 2016**/

package main

import (
    "github.com/BlueMedora/bluemedora-firehose-nozzle/logger"
    "github.com/BlueMedora/bluemedora-firehose-nozzle/nozzleconfiguration"
    "github.com/BlueMedora/bluemedora-firehose-nozzle/bluemedorafirehosenozzle"
)

func main() {
    logger := logger.New()
    logger.Debug("working log")
    
    //Read in config
    config, err := nozzleconfiguration.New(logger)
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