// Copyright (c) 2016 Blue Medora, Inc. All rights reserved.
// This file is subject to the terms and conditions defined in the included file 'LICENSE.txt'.

package logger

import (
    "os"
    "fmt"
    "path/filepath"
    "log"

    "github.com/cloudfoundry/gosteno"
)

//New logger
func New(logDirectory string, logFile string, loggerName string) *gosteno.Logger {
    loggingConfig := &gosteno.Config {
        Sinks:  []gosteno.Sink{
           gosteno.NewFileSink(fmt.Sprintf("%s/%s", logDirectory, logFile)),  
        },
        Level:      gosteno.LOG_DEBUG,
        Codec:      gosteno.NewJsonCodec(),
        EnableLOC:  true,
    }
    
    gosteno.Init(loggingConfig)
    return gosteno.NewLogger(loggerName)
}

//CreateLogDirectory clears out old directory and creates a new one
func CreateLogDirectory(logDirectory string) {
    log.Print("Finding absolute path to log directory")
    absoluteDirectoryPath, err := filepath.Abs(logDirectory)

    if err != nil {
        log.Printf("Error getting absolute path to log directory using relative path due to %v", err)
        absoluteDirectoryPath = logDirectory
    }

    log.Printf("Using path %s", absoluteDirectoryPath)

    if _, err := os.Stat(fmt.Sprintf("%s/", absoluteDirectoryPath)); err == nil {
        os.RemoveAll(fmt.Sprintf("%s", absoluteDirectoryPath))
    }
    
    os.MkdirAll(fmt.Sprintf("%s/", absoluteDirectoryPath), os.ModePerm)
}