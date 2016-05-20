package main

import (
    "flag"
    "log"
    "os"
    // "os"
    // "os/signal"
    // "runtime/pprof"
    // "syscall"
    
    "github.com/BlueMedora/bluemedora-firehose-nozzle/logger"
)

func main() {
    flag.Parse()
    
    logger := logger.New()
    //Read in config
    //Start nozzle
}