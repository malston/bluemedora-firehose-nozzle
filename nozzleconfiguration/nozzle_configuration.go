/**Copyright Blue Medora Inc. 2016**/

package nozzleconfiguration

import (
    "encoding/json"
    "io/ioutil"
    "fmt"
    
    "github.com/cloudfoundry/gosteno"
)

//NozzleConfiguration represents configuration file
type NozzleConfiguration struct {
    UAAURL                  string
    UAAUsername             string
    UAAPassword             string
    TrafficControllerURL    string
    DisableAccessControl    bool
    UseSSL                  bool
    IdleTimeout             uint32
}

//New NozzleConfiguration
func New(logger *gosteno.Logger) (*NozzleConfiguration, error) {
    configBuffer, err := ioutil.ReadFile("./config/bluemedora-firehose-nozzle.json")
    
    if err != nil {
        return nil, fmt.Errorf("Unable to load config file bluemedora-firehose-nozzle.json: %s", err)
    }
    
    var nozzleConfig NozzleConfiguration
    err = json.Unmarshal(configBuffer, &nozzleConfig)
    if err != nil {
        return nil, fmt.Errorf("Error parsing config file bluemedora-firehose-nozzle.json: %s", err)
    }
    
    logger.Debug(fmt.Sprintf("Loaded configuration to UAAURL <%s>, UAA Username <%s>, Traffic Controller URL <%s>, Disable Access Control <%v>, Using SSL <%v>", 
        nozzleConfig.UAAURL, nozzleConfig.UAAUsername, nozzleConfig.TrafficControllerURL, nozzleConfig.DisableAccessControl, nozzleConfig.UseSSL))
        
    return &nozzleConfig, nil
}