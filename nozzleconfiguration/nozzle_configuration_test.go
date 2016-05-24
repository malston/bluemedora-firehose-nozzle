/**Copyright Blue Medora Inc. 2016**/

package nozzleconfiguration

import (
    "testing"
    "os"
    "fmt"
    "encoding/json"
    "io/ioutil"
    
    "github.com/BlueMedora/bluemedora-firehose-nozzle/logger"
) 

var (
    defaultLogDirector = "./logs"
    configFile = "../config/bluemedora-firehose-nozzle.json"
    tempConfigFile = "../config/bluemedora-firehose-nozzle.json.real"
    testUAAURL = "UAAURL"
    testUsername = "username"
    testPassword = "password"
    testTrafficControllerURL = "traffic_url"
    testDisableAccessControl = false
    testUseSSL = false
    testIdleTimeout = uint32(60)
)

func TestConfigParsing(t *testing.T) {
    //Setup Environment
    err := setupEnvironment(t)
    if err != nil {
        tearDownEnvironment(t)
        t.Fatalf("Setup failed due to: %s", err.Error())
    }
    
    t.Log("Creating configuration...")
    logger := logger.New(defaultLogDirector)
    
    //Create new configuration
    var config *NozzleConfiguration
    config, err = New("../config/bluemedora-firehose-nozzle.json", logger)
    
    if err != nil {
        tearDownEnvironment(t)
        t.Fatalf("Error occrued while creating configuration %s", err)
    }
    
    //Test values
    t.Log(fmt.Sprintf("Checking UAA URL... (expected value: %s)", testUAAURL))
    if config.UAAURL != testUAAURL {
        t.Errorf("Expected UAA URL of %s, but received %s", testUAAURL, config.UAAURL)
    }
    
    t.Log(fmt.Sprintf("Checking UAA Username... (expected value: %s)", testUsername))
    if config.UAAUsername != testUsername {
        t.Errorf("Expected UAA Username of %s, but received %s", testUsername, config.UAAUsername)
    }

    t.Log(fmt.Sprintf("Checking UAA Password... (expected value: %s)", testPassword))
    if config.UAAPassword != testPassword {
        t.Errorf("Expected UAA Password of %s, but received %s", testPassword, config.UAAPassword)
    }
    
    t.Log(fmt.Sprintf("Checking Traffic Controller URL... (expected value: %s)", testTrafficControllerURL))
    if config.TrafficControllerURL != testTrafficControllerURL {
        t.Errorf("Expected Traffic Controller URL of %s, but received %s", testTrafficControllerURL, config.TrafficControllerURL)
    }
    
    t.Log(fmt.Sprintf("Checking Disable Access Control... (expected value: %v)", testDisableAccessControl))
    if config.DisableAccessControl != testDisableAccessControl {
        t.Errorf("Expected Disable Access Control of %v, but received %v", testDisableAccessControl, config.DisableAccessControl)
    }

    t.Log(fmt.Sprintf("Checking Use SSL... (expected value: %v)", testUseSSL))
    if config.UseSSL != testUseSSL {
        t.Errorf("Expected Use SSL of %v, but received %v", testUseSSL, config.UseSSL)
    }
    
    t.Log(fmt.Sprintf("Checking Idle Timeout... (expected value: %v)", testIdleTimeout))
    if config.IdleTimeoutSeconds != testIdleTimeout {
        t.Errorf("Expected Idle Timeout of %v, but received %v", testIdleTimeout, config.IdleTimeoutSeconds)
    }
    
    err = tearDownEnvironment(t)
    if err != nil {
        t.Fatalf("Tear down failed due to: %s", err.Error())
    }
}

func setupEnvironment(t *testing.T) error {
    t.Log("Setting up environment...")
    
    err := os.Rename(configFile, tempConfigFile)
    if err != nil {
        return fmt.Errorf("Error renaming config file. Ensure bluemedora-firehose-nozzle.json exists in config directory: %s", err)
    }

    message := NozzleConfiguration{testUAAURL, testUsername, testPassword, testTrafficControllerURL, testDisableAccessControl, testUseSSL, testIdleTimeout}
    messageBytes, _ := json.Marshal(message)
    
    err = ioutil.WriteFile(configFile, messageBytes, os.ModePerm)
    if err != nil {
        return fmt.Errorf("Error creating new config file: %s", err)
    }
    
    t.Log("Setup test environment")
    return nil
}

func tearDownEnvironment(t *testing.T) error {
    t.Log("Tearing down test environment...")
    if _, err := os.Stat(tempConfigFile); os.IsNotExist(err) {
        t.Log("bluemedora-firehose-nozzle.json.real not found no clean up needed")
        return nil
    }
    
    if _, err := os.Stat(configFile); err == nil {
        err = os.Remove(configFile)
        if err != nil {
            return fmt.Errorf("Error removing test config file: %s", err)
        }
    }
    
    err := os.Rename(tempConfigFile, configFile)
    if err != nil {
        return fmt.Errorf("Error renaming config file. Ensure bluemedora-firehose-nozzle.json exists in config directory: %s", err)
    }
    
    t.Log("Tore down test environment")
    return nil
}