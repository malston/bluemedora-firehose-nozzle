// Copyright (c) 2016 Blue Medora, Inc. All rights reserved.
// This file is subject to the terms and conditions defined in the included file 'LICENSE.txt'.

package nozzleconfiguration

import (
    "testing"
    "os"
    "fmt"
    "encoding/json"
    "io/ioutil"
    "strings"
    
    "github.com/BlueMedora/bluemedora-firehose-nozzle/logger"
) 

const (
    defaultLogDirectory = "../logs"
	nozzleLogFile       = "bm_nozzle.log"
	nozzleLogName       = "bm_firehose_nozzle"
    
    configFile = "../config/bluemedora-firehose-nozzle.json"
    tempConfigFile = "../config/bluemedora-firehose-nozzle.json.real"
    
    testUAAURL = "UAAURL"
    testUsername = "username"
    testPassword = "password"
    testTrafficControllerURL = "traffic_url"
    testDisableAccessControl = false
    testInsecureSSLSkipVerify = false
    testIdleTimeout = uint32(60)
    testMetricCacheDuration = uint32(60)
    testWebServerPort = uint32(8081)
)

func TestConfigParsing(t *testing.T) {
    //Setup Environment
    err := setupGoodEnvironment(t)
    if err != nil {
        tearDownEnvironment(t)
        t.Fatalf("Setup failed due to: %s", err.Error())
    }
    
    t.Log("Creating configuration...")
    logger.CreateLogDirectory(defaultLogDirectory)
    logger := logger.New(defaultLogDirectory, nozzleLogFile, nozzleLogName)
    
    //Create new configuration
    var config *NozzleConfiguration
    config, err = New(configFile, logger)
    
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

    t.Log(fmt.Sprintf("Checking Insecure SSL Skip Verify... (expected value: %v)", testInsecureSSLSkipVerify))
    if config.InsecureSSLSkipVerify != testInsecureSSLSkipVerify {
        t.Errorf("Expected Insecure SSL Skip Verify of %v, but received %v", testInsecureSSLSkipVerify, config.InsecureSSLSkipVerify)
    }
    
    t.Log(fmt.Sprintf("Checking Idle Timeout... (expected value: %v)", testIdleTimeout))
    if config.IdleTimeoutSeconds != testIdleTimeout {
        t.Errorf("Expected Idle Timeout of %v, but received %v", testIdleTimeout, config.IdleTimeoutSeconds)
    }
    
    t.Log(fmt.Sprintf("Checking Metric Cache Duration... (expected value: %v)", testMetricCacheDuration))
    if config.MetricCacheDurationSeconds != testMetricCacheDuration {
        t.Errorf("Expected Metric Cache Duration of %v, but received %v", testMetricCacheDuration, config.MetricCacheDurationSeconds)
    }
    
    t.Log(fmt.Sprintf("Checking Web Server Port... (expected value: %v)", testWebServerPort))
    if config.WebServerPort != testWebServerPort {
        t.Errorf("Expected Web Server Port of %v, but received %v", testWebServerPort, config.WebServerPort)
    }
    
    err = tearDownEnvironment(t)
    if err != nil {
        t.Fatalf("Tear down failed due to: %s", err.Error())
    }
}

func TestBadConfigFile(t *testing.T) {
    err := setupBadEnvironment(t)
    if err != nil {
        tearDownEnvironment(t)
        t.Fatalf("Setup failed due to: %s", err.Error())
    }
    
    logger.CreateLogDirectory(defaultLogDirectory)
    logger := logger.New(defaultLogDirectory, nozzleLogFile, nozzleLogName)
    
    //Create new configuration
    t.Log("Checking loading of bad config file... (expecting error)")
    _, err = New(configFile, logger)
    
    if err != nil {
        if !strings.Contains(err.Error(), "Error parsing config file bluemedora-firehose-nozzle.json:") {
            t.Errorf("Expected error containing %s, but received %s", "Error parsing config file bluemedora-firehose-nozzle.json:", err.Error())
        }
    } else {
        t.Errorf("Expected error from loading a bad config file, but loaded correctly")
    }
    
    err = tearDownEnvironment(t)
    if err != nil {
        t.Fatalf("Tear down failed due to: %s", err.Error())
    }
}

func TestNoConfigFile(t *testing.T) {
    t.Log("Creating configuration...")
    logger.CreateLogDirectory(defaultLogDirectory)
    logger := logger.New(defaultLogDirectory, nozzleLogFile, nozzleLogName)
    
    //Create new configuration
    t.Log("Checking loading of non-existent file... (expecting error)")
    _, err := New("fake_file.json", logger)
    
    if err != nil {
        if !strings.Contains(err.Error(), "Unable to load config file bluemedora-firehose-nozzle.json:") {
            t.Errorf("Expected error containing %s, but received %s", "Unable to load config file bluemedora-firehose-nozzle.json:", err.Error())
        }
    } else {
        t.Errorf("Expected error from loading non-existsent file, but loaded correctly")
    }
}

func setupGoodEnvironment(t *testing.T) error {
    t.Log("Setting up good environment...")
    
    err := renameConfigFile(t)
    if err != nil {
        return err
    }

    err = createGoodConfigFile(t)
    if err != nil {
        return err
    }
    
    t.Log("Setup good test environment")
    return nil
}

func setupBadEnvironment(t *testing.T) error {
    t.Log("Setting up bad envrionment...")
    
    err := renameConfigFile(t)
    if err != nil {
        return err
    }
    
    err = createBadConfigFile(t)
    if err != nil {
        return err
    }
    
    
    t.Log("Setup bad test envrionment")
    return nil
}

func renameConfigFile(t *testing.T) error {
    t.Log("Renaming real config file...")
    
    err := os.Rename(configFile, tempConfigFile)
    if err != nil {
        return fmt.Errorf("Error renaming config file. Ensure bluemedora-firehose-nozzle.json exists in config directory: %s", err)
    }
    
    t.Log("Renamed real config file")
    return nil
}

func createGoodConfigFile(t *testing.T) error {
    t.Log("Creating good config file...")
    
    message := NozzleConfiguration{
        testUAAURL, testUsername, 
        testPassword, testTrafficControllerURL, 
        testDisableAccessControl, testInsecureSSLSkipVerify, 
        testIdleTimeout, testMetricCacheDuration,
        testWebServerPort}
        
    messageBytes, _ := json.Marshal(message)
    
    err := ioutil.WriteFile(configFile, messageBytes, os.ModePerm)
    if err != nil {
        return fmt.Errorf("Error creating good config file: %s", err)
    }
    
    t.Log("Created good config file")
    return nil
}

func createBadConfigFile(t *testing.T) error {
    t.Log("Creating bad config file...")
    
    _, err := os.Create(configFile)
    if err != nil {
        return fmt.Errorf("Error creating bad config file: %s", err)
    }
    
    t.Log("Created bad config file")
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