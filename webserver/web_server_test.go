// Copyright (c) 2016 Blue Medora, Inc. All rights reserved.
// This file is subject to the terms and conditions defined in the included file 'LICENSE.txt'.

package webserver

// TODO figure out how to generate certs in order for test to work

import (
	"fmt"
	"testing"
	"net/http"
	"crypto/tls"

	"github.com/BlueMedora/bluemedora-firehose-nozzle/logger"
	"github.com/BlueMedora/bluemedora-firehose-nozzle/nozzleconfiguration"
	"github.com/BlueMedora/bluemedora-firehose-nozzle/testhelpers"
)

const (
	defaultConfigLocation = "../config/bluemedora-firehose-nozzle.json"

	defaultLogDirectory = "../logs"
	webserverLogFile    = "bm_server.log"
	webserverLogName    = "bm_server"
	
	testCertLocation = "../certs/cert.pem"
	testKeyLocation  = "../certs/key.pem"
)

func TestEndpoints(t *testing.T) {
	server, config := createWebServer(t)
	
	t.Log("Setting up server envrionment...")
	testhelpers.GenerateCertFiles()
	errors := server.Start(testKeyLocation, testCertLocation)
	
	//Handle errors from server
	go func() {
		select {
			case err := <-errors:
				if err != nil {
					t.Fatalf("Error with server: %s", err.Error())
				}
		}
	}()
	
	client := createHTTPClient(t)

	//Token tests
	tokenEndPointTest(t, client, config)
	
}

func tokenEndPointTest(t *testing.T, client *http.Client, config *nozzleconfiguration.NozzleConfiguration) {
	t.Log("Running token request tests...")
	goodTokenRequestTest(t, client, config)
	badCredentialTokenTest(t, client, config)
	noCredentialTokenTest(t, client, config)
	t.Log("Finished token request tests")
}

func goodTokenRequestTest(t *testing.T, client *http.Client, config *nozzleconfiguration.NozzleConfiguration) {
	tokenRequest := createTokenRequest(config.UAAUsername, config.UAAPassword, config.WebServerPort, t)
	
	t.Logf("Check if server responses to good token request... (expecting status code: %v)", http.StatusOK)
	response, err := client.Do(tokenRequest)
	if err != nil {
		t.Fatalf("Error occured while requesting token: %s", err.Error())
	}
	
	if response.StatusCode != http.StatusOK {
		t.Errorf("Expecting status code %v, but received %v", http.StatusOK, response.StatusCode)
	}
}

func badCredentialTokenTest(t *testing.T, client *http.Client, config *nozzleconfiguration.NozzleConfiguration) {
	tokenRequest := createTokenRequest("baduser", "badPass", config.WebServerPort, t)
	
	t.Logf("Check if server responses to a bad credential token request... (expecting status code: %v)", http.StatusUnauthorized)
	response, err := client.Do(tokenRequest)
	if err != nil {
		t.Fatalf("Error occured while requesting token: %s", err.Error())
	}
	
	if response.StatusCode != http.StatusUnauthorized {
		t.Errorf("Expecting status code %v, but received %v", http.StatusUnauthorized, response.StatusCode)
	}
}

func noCredentialTokenTest(t *testing.T, client *http.Client, config *nozzleconfiguration.NozzleConfiguration) {
	tokenRequest := createTokenRequest("", "", config.WebServerPort, t)
	
	t.Logf("Check if server responses to a no credential token request... (expecting status code: %v)", http.StatusBadRequest)
	response, err := client.Do(tokenRequest)
	if err != nil {
		t.Fatalf("Error occured while requesting token: %s", err.Error())
	}
	
	if response.StatusCode != http.StatusBadRequest {
		t.Errorf("Expecting status code %v, but received %v", http.StatusBadRequest, response.StatusCode)
	}
}

func createWebServer(t *testing.T) (*WebServer, *nozzleconfiguration.NozzleConfiguration) {
	t.Log("Creating webserver...")
	logger.CreateLogDirectory(defaultLogDirectory)
	logger := logger.New(defaultLogDirectory, webserverLogFile, webserverLogName)

	config, err := nozzleconfiguration.New(defaultConfigLocation, logger)
	if err != nil {
		t.Fatalf("Error while loading configuration: %s", err.Error())
	}

	t.Log("Created webserver")
	return New(config, logger), config
}

func createHTTPClient(t *testing.T) *http.Client {
	t.Log("Creating a test client...")
	tr := &http.Transport {
		TLSClientConfig: &tls.Config {InsecureSkipVerify: true},
	}
	
	t.Log("Created a test client")
	return &http.Client {Transport: tr}
}

func createTokenRequest(username string, password string, port uint32, t *testing.T) *http.Request {
	t.Log("Creating token request...")
	request, err := http.NewRequest("GET", fmt.Sprintf("https://localhost:%d/token", port), nil)
	if err != nil {
		t.Fatalf("Error creating token request: %s", err.Error())
	}
	
	if username != "" || password != "" {
		request.Header.Add(headerUsernameKey, username)
		request.Header.Add(headerPasswordKey, password)
	}
	
	t.Log("Created token request")
	return request
}
