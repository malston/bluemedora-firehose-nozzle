// Copyright (c) 2016 Blue Medora, Inc. All rights reserved.
// This file is subject to the terms and conditions defined in the included file 'LICENSE.txt'.

package logger

import (
    "os"
    "fmt"
    
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
    if _, err := os.Stat(fmt.Sprintf("%s/", logDirectory)); err == nil {
        os.RemoveAll(fmt.Sprintf("%s", logDirectory))
    }
    
    os.MkdirAll(fmt.Sprintf("%s/", logDirectory), os.ModePerm)
}