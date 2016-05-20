package logger

import (
    "os"
    "strings"
    
    "github.com/cloudfoundry/gosteno"
)

func New() *gosteno.Logger {
    loggingConfig := &gosteno.Config {
        Sinks:  []gosteno.Sink{
           gosteno.NewFileSink("./logs/bm_nozzle.log"),  
        },
        Level:      gosteno.LOG_DEBUG,
        Codec:      gosteno.NewJsonCodec(),
        EnableLOC:  true,
    }
    
    gosteno.Init(loggingConfig)
    return gosteno.newLogger("bm_firehose_nozzle")
}