/**Copyright Blue Medora Inc. 2016**/

package webserver

import (
    // "net/http"
    "sync"
    
    "github.com/cloudfoundry/gosteno"
    "github.com/BlueMedora/bluemedora-firehose-nozzle/nozzleconfiguration"
    "github.com/BlueMedora/bluemedora-firehose-nozzle/webtoken"
)

//WebServer REST endpoint for sending data
type WebServer struct {
    logger  *gosteno.Logger
    mutext  sync.Mutex
    config  *nozzleconfiguration.NozzleConfiguration
    tokens  map[string]*webtoken.Token //Maps token string to token object
}

// func New(logger *gosteno.Logger)

//handle login post
//handle each resource metric request
//handle token timeout