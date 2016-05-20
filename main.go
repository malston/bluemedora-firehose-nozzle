package main

import (
    "flag"
    
    "github.com/BlueMedora/bluemedora-firehose-nozzle/logger"
)

func main() {
    flag.Parse()
    
    logger := logger.New()
    logger.Debug("working log")
    //Read in config
    //Start nozzle
}