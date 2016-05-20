/**Copyright Blue Medora Inc. 2016**/

package logger

import (
    "os"
    
    "github.com/cloudfoundry/gosteno"
)

//New logger
func New() *gosteno.Logger {
    createLogDirectory()
    
    
    loggingConfig := &gosteno.Config {
        Sinks:  []gosteno.Sink{
           gosteno.NewFileSink("./logs/bm_nozzle.log"),  
        },
        Level:      gosteno.LOG_DEBUG,
        Codec:      gosteno.NewJsonCodec(),
        EnableLOC:  true,
    }
    
    gosteno.Init(loggingConfig)
    return gosteno.NewLogger("bm_firehose_nozzle")
}

func createLogDirectory() {
    if _, err := os.Stat("./logs/"); err == nil {
        os.RemoveAll("./logs")
    }
    
    os.MkdirAll("./logs/", os.ModePerm)
}