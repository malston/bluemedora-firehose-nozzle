// Copyright (c) 2016 Blue Medora, Inc. All rights reserved.
// This file is subject to the terms and conditions defined in the included file 'LICENSE.txt'.

package nozzleconfiguration

import (
    "encoding/json"
    "io/ioutil"
    "fmt"
    
    "github.com/cloudfoundry/gosteno"
)

//NozzleConfiguration represents configuration file
type NozzleConfiguration struct {
    UAAURL                      string
    UAAUsername                 string
    UAAPassword                 string
    TrafficControllerURL        string
    DisableAccessControl        bool
    InsecureSSLSkipVerify       bool
    IdleTimeoutSeconds          uint32
    MetricCacheDurationSeconds  uint32
    WebServerPort               uint32
}

//New NozzleConfiguration
func New(configPath string, logger *gosteno.Logger) (*NozzleConfiguration, error) {
    configBuffer, err := ioutil.ReadFile(configPath)
    
    if err != nil {
        return nil, fmt.Errorf("Unable to load config file bluemedora-firehose-nozzle.json: %s", err)
    }
    
    var nozzleConfig NozzleConfiguration
    err = json.Unmarshal(configBuffer, &nozzleConfig)
    if err != nil {
        return nil, fmt.Errorf("Error parsing config file bluemedora-firehose-nozzle.json: %s", err)
    }
    
    logger.Debug(fmt.Sprintf("Loaded configuration to UAAURL <%s>, UAA Username <%s>, Traffic Controller URL <%s>, Disable Access Control <%v>, Insecure SSL Skip Verify <%v>", 
        nozzleConfig.UAAURL, nozzleConfig.UAAUsername, nozzleConfig.TrafficControllerURL, nozzleConfig.DisableAccessControl, nozzleConfig.InsecureSSLSkipVerify))
        
    return &nozzleConfig, nil
}