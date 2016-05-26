/**Copyright Blue Medora Inc. 2016**/

package webserver

//TODO figure out how to generate certs in order for test to work

// import (
//     "testing"
    
//     "github.com/BlueMedora/bluemedora-firehose-nozzle/logger"
//     "github.com/BlueMedora/bluemedora-firehose-nozzle/nozzleconfiguration"
// )

// const (
//     defaultConfigLocation = "../config/bluemedora-firehose-nozzle.json"

// 	defaultLogDirectory = "../logs"
//     webserverLogFile = "bm_server.log"
// 	webserverLogName = "bm_server"
// )

// var (
//     config *nozzleconfiguration.NozzleConfiguration
// )

// func TestTokenGeneration(t *testing.T) {
//     server := createWebServer(t)
//     server.Start(DefaultKeyLocation, DefaultCertLocation)
// }

// func createWebServer(t *testing.T) *WebServer {
//     t.Log("Creating webserver...")
//     logger := logger.New(defaultLogDirectory, webserverLogFile, webserverLogName)
    
//     config, err := nozzleconfiguration.New(defaultConfigLocation, logger)
//     if err != nil {
//         t.Fatalf("Error while loading configuration: %s", err.Error())
//     }
    
//     t.Log("Created webserver")
//     return New(config, logger)
// }