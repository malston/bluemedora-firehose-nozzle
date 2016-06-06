// Copyright (c) 2016 Blue Medora, Inc. All rights reserved.
// This file is subject to the terms and conditions defined in the included file 'LICENSE.txt'.

package webserver

import (
	"fmt"
	"testing"
	"net/http"
	"crypto/tls"

	"github.com/BlueMedora/bluemedora-firehose-nozzle/logger"
	"github.com/BlueMedora/bluemedora-firehose-nozzle/nozzleconfiguration"
	"github.com/BlueMedora/bluemedora-firehose-nozzle/testhelpers"
	"github.com/cloudfoundry/sonde-go/events"
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
	
	//Retrieve token for other endpoint test
	token := getToken(t, client, config)
	
	//Metron Agent tests
	metronAgentEndPointTest(t, client, token, config.WebServerPort, server)
	
	//Syslog Drain Binder tests
	syslogDrainBinderEndPointTest(t, client, token, config.WebServerPort, server)
	
	//TPS Watcher tests
	tpsWatcherEndPointTest(t, client, token, config.WebServerPort, server)
	
	//TPS Listener tests
	tpsListenerEndPointTest(t, client, token, config.WebServerPort, server)
	
	//Stager tests
	stagerEndPointTest(t, client, token, config.WebServerPort, server)
	
	//Route Emitter tests
	routeEmittersEndPointTest(t, client, token, config.WebServerPort, server)
	
	//Rep tests
	repEndPointTest(t, client, token, config.WebServerPort, server)
	
	//Receptor tests
	receptorEndPointTest(t, client, token, config.WebServerPort, server)
	
	//NSYNC Listener tests
	nsyncListenersEndPointTest(t, client, token, config.WebServerPort, server)
	
	//NSYNC Bulker tests
	nsyncBulkersEndPointTest(t, client, token, config.WebServerPort, server)
	
	//Garden Linux tests
	gardenLinuxEndPointTest(t, client, token, config.WebServerPort, server)
	
	//File Server tests
	fileServersEndPointTest(t, client, token, config.WebServerPort, server)
	
	//Fetcher tests
	fetchersEndPointTest(t, client, token, config.WebServerPort, server)
	
	//Converger tests
	convergerEndPointTest(t, client, token, config.WebServerPort, server)
	
	//CC Uploader tests
	ccUploaderEndPointTest(t, client, token, config.WebServerPort, server)
	
	//bbs tests
	bbsEndPointTest(t, client, token, config.WebServerPort, server)
	
	//Auctioneer tests
	auctioneerEndPointTest(t, client, token, config.WebServerPort, server)
	
	//etcd tests
	etcdEndPointTest(t, client, token, config.WebServerPort, server)
	
	//Doppler Servers tests
	dopplerServerEndPointTest(t, client, token, config.WebServerPort, server)
	
	//Cloud Controller tests
	cloudControllerEndPointTest(t, client, token, config.WebServerPort, server)
	
	//Traffic Controller tests
	trafficControllerEndPointTest(t, client, token, config.WebServerPort, server)
	
	//Go Router tests
	goRouterEndPointTest(t, client, token, config.WebServerPort, server)
}

/** Tests **/
func tokenEndPointTest(t *testing.T, client *http.Client, config *nozzleconfiguration.NozzleConfiguration) {
	t.Log("Running token request tests...")
	badCredentialTokenTest(t, client, config)
	noCredentialTokenTest(t, client, config)
	goodTokenRequestTest(t, client, config)
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

func metronAgentEndPointTest(t *testing.T, client *http.Client, token string, port uint32, server *WebServer) {
	cacheEnvelope(metronAgentOrigin, server)
	
	request := createResourceRequest(t, token, port, "metron_agents")
	
	t.Logf("Check if server response to valid /metron_agents request... (expecting status code: %v)", http.StatusOK)
	response, err := client.Do(request)
	
	if err != nil {
		t.Errorf("Error occured while hitting endpoint: %s", err.Error())
	} else if response.StatusCode != http.StatusOK {
		t.Errorf("Expecting status code %v, but received %v", http.StatusOK, response.StatusCode)
	}
}

func syslogDrainBinderEndPointTest(t *testing.T, client *http.Client, token string, port uint32, server *WebServer) {
	cacheEnvelope(syslogDrainBinderOrigin, server)
	
	request := createResourceRequest(t, token, port, "syslog_drains")
	
	t.Logf("Check if server response to valid /syslog_drains request... (expecting status code: %v)", http.StatusOK)
	response, err := client.Do(request)
	
	if err != nil {
		t.Errorf("Error occured while hitting endpoint: %s", err.Error())
	} else if response.StatusCode != http.StatusOK {
		t.Errorf("Expecting status code %v, but received %v", http.StatusOK, response.StatusCode)
	}
}

func tpsWatcherEndPointTest(t *testing.T, client *http.Client, token string, port uint32, server *WebServer) {
	cacheEnvelope(tpsWatcherOrigin, server)
	
	request := createResourceRequest(t, token, port, "tps_watchers")
	
	t.Logf("Check if server response to valid /tps_watchers request... (expecting status code: %v)", http.StatusOK)
	response, err := client.Do(request)
	
	if err != nil {
		t.Errorf("Error occured while hitting endpoint: %s", err.Error())
	} else if response.StatusCode != http.StatusOK {
		t.Errorf("Expecting status code %v, but received %v", http.StatusOK, response.StatusCode)
	}
}

func tpsListenerEndPointTest(t *testing.T, client *http.Client, token string, port uint32, server *WebServer) {
	cacheEnvelope(tpsListenerOrigin, server)
	
	request := createResourceRequest(t, token, port, "tps_listeners")
	
	t.Logf("Check if server response to valid /tps_listeners request... (expecting status code: %v)", http.StatusOK)
	response, err := client.Do(request)
	
	if err != nil {
		t.Errorf("Error occured while hitting endpoint: %s", err.Error())
	} else if response.StatusCode != http.StatusOK {
		t.Errorf("Expecting status code %v, but received %v", http.StatusOK, response.StatusCode)
	}
}

func stagerEndPointTest(t *testing.T, client *http.Client, token string, port uint32, server *WebServer) {
	cacheEnvelope(stagerOrigin, server)
	
	request := createResourceRequest(t, token, port, "stagers")
	
	t.Logf("Check if server response to valid /stagers request... (expecting status code: %v)", http.StatusOK)
	response, err := client.Do(request)
	
	if err != nil {
		t.Errorf("Error occured while hitting endpoint: %s", err.Error())
	} else if response.StatusCode != http.StatusOK {
		t.Errorf("Expecting status code %v, but received %v", http.StatusOK, response.StatusCode)
	}
}

func sshProxyEndPointTest(t *testing.T, client *http.Client, token string, port uint32, server *WebServer) {
	cacheEnvelope(sshProxyOrigin, server)
	
	request := createResourceRequest(t, token, port, "ssh_proxies")
	
	t.Logf("Check if server response to valid /ssh_proxies request... (expecting status code: %v)", http.StatusOK)
	response, err := client.Do(request)
	
	if err != nil {
		t.Errorf("Error occured while hitting endpoint: %s", err.Error())
	} else if response.StatusCode != http.StatusOK {
		t.Errorf("Expecting status code %v, but received %v", http.StatusOK, response.StatusCode)
	}
}

func sendersEndPointTest(t *testing.T, client *http.Client, token string, port uint32, server *WebServer) {
	cacheEnvelope(senderOrigin, server)
	
	request := createResourceRequest(t, token, port, "senders")
	
	t.Logf("Check if server response to valid /senders request... (expecting status code: %v)", http.StatusOK)
	response, err := client.Do(request)
	
	if err != nil {
		t.Errorf("Error occured while hitting endpoint: %s", err.Error())
	} else if response.StatusCode != http.StatusOK {
		t.Errorf("Expecting status code %v, but received %v", http.StatusOK, response.StatusCode)
	}
}

func routeEmittersEndPointTest(t *testing.T, client *http.Client, token string, port uint32, server *WebServer) {
	cacheEnvelope(routeEmitterOrigin, server)
	
	request := createResourceRequest(t, token, port, "route_emitters")
	
	t.Logf("Check if server response to valid /route_emitters request... (expecting status code: %v)", http.StatusOK)
	response, err := client.Do(request)
	
	if err != nil {
		t.Errorf("Error occured while hitting endpoint: %s", err.Error())
	} else if response.StatusCode != http.StatusOK {
		t.Errorf("Expecting status code %v, but received %v", http.StatusOK, response.StatusCode)
	}
}

func repEndPointTest(t *testing.T, client *http.Client, token string, port uint32, server *WebServer) {
	cacheEnvelope(repOrigin, server)
	
	request := createResourceRequest(t, token, port, "reps")
	
	t.Logf("Check if server response to valid /reps request... (expecting status code: %v)", http.StatusOK)
	response, err := client.Do(request)
	
	if err != nil {
		t.Errorf("Error occured while hitting endpoint: %s", err.Error())
	} else if response.StatusCode != http.StatusOK {
		t.Errorf("Expecting status code %v, but received %v", http.StatusOK, response.StatusCode)
	}
}

func receptorEndPointTest(t *testing.T, client *http.Client, token string, port uint32, server *WebServer) {
	cacheEnvelope(receptorOrigin, server)
	
	request := createResourceRequest(t, token, port, "receptors")
	
	t.Logf("Check if server response to valid /receptors request... (expecting status code: %v)", http.StatusOK)
	response, err := client.Do(request)
	
	if err != nil {
		t.Errorf("Error occured while hitting endpoint: %s", err.Error())
	} else if response.StatusCode != http.StatusOK {
		t.Errorf("Expecting status code %v, but received %v", http.StatusOK, response.StatusCode)
	}
}

func nsyncListenersEndPointTest(t *testing.T, client *http.Client, token string, port uint32, server *WebServer) {
	cacheEnvelope(nsyncListenerOrigin, server)
	
	request := createResourceRequest(t, token, port, "nsync_listeners")
	
	t.Logf("Check if server response to valid /nsync_listeners request... (expecting status code: %v)", http.StatusOK)
	response, err := client.Do(request)
	
	if err != nil {
		t.Errorf("Error occured while hitting endpoint: %s", err.Error())
	} else if response.StatusCode != http.StatusOK {
		t.Errorf("Expecting status code %v, but received %v", http.StatusOK, response.StatusCode)
	}
}

func nsyncBulkersEndPointTest(t *testing.T, client *http.Client, token string, port uint32, server *WebServer) {
	cacheEnvelope(nsyncBulkerOrigin, server)
	
	request := createResourceRequest(t, token, port, "nsync_bulkers")
	
	t.Logf("Check if server response to valid /nsync_bulkers request... (expecting status code: %v)", http.StatusOK)
	response, err := client.Do(request)
	
	if err != nil {
		t.Errorf("Error occured while hitting endpoint: %s", err.Error())
	} else if response.StatusCode != http.StatusOK {
		t.Errorf("Expecting status code %v, but received %v", http.StatusOK, response.StatusCode)
	}
}

func gardenLinuxEndPointTest(t *testing.T, client *http.Client, token string, port uint32, server *WebServer) {
	cacheEnvelope(gardenLinuxOrigin, server)
	
	request := createResourceRequest(t, token, port, "garden_linuxs")
	
	t.Logf("Check if server response to valid /garden_linuxs request... (expecting status code: %v)", http.StatusOK)
	response, err := client.Do(request)
	
	if err != nil {
		t.Errorf("Error occured while hitting endpoint: %s", err.Error())
	} else if response.StatusCode != http.StatusOK {
		t.Errorf("Expecting status code %v, but received %v", http.StatusOK, response.StatusCode)
	}
}

func fileServersEndPointTest(t *testing.T, client *http.Client, token string, port uint32, server *WebServer) {
	cacheEnvelope(fileServerOrigin, server)
	
	request := createResourceRequest(t, token, port, "file_servers")
	
	t.Logf("Check if server response to valid /file_servers request... (expecting status code: %v)", http.StatusOK)
	response, err := client.Do(request)
	
	if err != nil {
		t.Errorf("Error occured while hitting endpoint: %s", err.Error())
	} else if response.StatusCode != http.StatusOK {
		t.Errorf("Expecting status code %v, but received %v", http.StatusOK, response.StatusCode)
	}
}

func fetchersEndPointTest(t *testing.T, client *http.Client, token string, port uint32, server *WebServer) {
	cacheEnvelope(fetcherOrigin, server)
	
	request := createResourceRequest(t, token, port, "fetchers")
	
	t.Logf("Check if server response to valid /fetchers request... (expecting status code: %v)", http.StatusOK)
	response, err := client.Do(request)
	
	if err != nil {
		t.Errorf("Error occured while hitting endpoint: %s", err.Error())
	} else if response.StatusCode != http.StatusOK {
		t.Errorf("Expecting status code %v, but received %v", http.StatusOK, response.StatusCode)
	}
}

func convergerEndPointTest(t *testing.T, client *http.Client, token string, port uint32, server *WebServer) {
	cacheEnvelope(convergerOrigin, server)
	
	request := createResourceRequest(t, token, port, "convergers")
	
	t.Logf("Check if server response to valid /convergers request... (expecting status code: %v)", http.StatusOK)
	response, err := client.Do(request)
	
	if err != nil {
		t.Errorf("Error occured while hitting endpoint: %s", err.Error())
	} else if response.StatusCode != http.StatusOK {
		t.Errorf("Expecting status code %v, but received %v", http.StatusOK, response.StatusCode)
	}
}

func ccUploaderEndPointTest(t *testing.T, client *http.Client, token string, port uint32, server *WebServer) {
	cacheEnvelope(ccUploaderOrigin, server)
	
	request := createResourceRequest(t, token, port, "cc_uploaders")
	
	t.Logf("Check if server response to valid /cc_uploaders request... (expecting status code: %v)", http.StatusOK)
	response, err := client.Do(request)
	
	if err != nil {
		t.Errorf("Error occured while hitting endpoint: %s", err.Error())
	} else if response.StatusCode != http.StatusOK {
		t.Errorf("Expecting status code %v, but received %v", http.StatusOK, response.StatusCode)
	}
}

func bbsEndPointTest(t *testing.T, client *http.Client, token string, port uint32, server *WebServer) {
	cacheEnvelope(bbsOrigin, server)
	
	request := createResourceRequest(t, token, port, "bbs")
	
	t.Logf("Check if server response to valid /bbs request... (expecting status code: %v)", http.StatusOK)
	response, err := client.Do(request)
	
	if err != nil {
		t.Errorf("Error occured while hitting endpoint: %s", err.Error())
	} else if response.StatusCode != http.StatusOK {
		t.Errorf("Expecting status code %v, but received %v", http.StatusOK, response.StatusCode)
	}
}

func auctioneerEndPointTest(t *testing.T, client *http.Client, token string, port uint32, server *WebServer) {
	cacheEnvelope(auctioneerOrigin, server)
	
	request := createResourceRequest(t, token, port, "auctioneers")
	
	t.Logf("Check if server response to valid /auctioneers request... (expecting status code: %v)", http.StatusOK)
	response, err := client.Do(request)
	
	if err != nil {
		t.Errorf("Error occured while hitting endpoint: %s", err.Error())
	} else if response.StatusCode != http.StatusOK {
		t.Errorf("Expecting status code %v, but received %v", http.StatusOK, response.StatusCode)
	}
}

func etcdEndPointTest(t *testing.T, client *http.Client, token string, port uint32, server *WebServer) {
	cacheEnvelope(etcdOrigin, server)
	
	request := createResourceRequest(t, token, port, "etcds")
	
	t.Logf("Check if server response to valid /etcds request... (expecting status code: %v)", http.StatusOK)
	response, err := client.Do(request)
	
	if err != nil {
		t.Errorf("Error occured while hitting endpoint: %s", err.Error())
	} else if response.StatusCode != http.StatusOK {
		t.Errorf("Expecting status code %v, but received %v", http.StatusOK, response.StatusCode)
	}
}

func dopplerServerEndPointTest(t *testing.T, client *http.Client, token string, port uint32, server *WebServer) {
	cacheEnvelope(dopplerServerOrigin, server)
	
	request := createResourceRequest(t, token, port, "doppler_servers")
	
	t.Logf("Check if server response to valid /doppler_servers request... (expecting status code: %v)", http.StatusOK)
	response, err := client.Do(request)
	
	if err != nil {
		t.Errorf("Error occured while hitting endpoint: %s", err.Error())
	} else if response.StatusCode != http.StatusOK {
		t.Errorf("Expecting status code %v, but received %v", http.StatusOK, response.StatusCode)
	}
}

func cloudControllerEndPointTest(t *testing.T, client *http.Client, token string, port uint32, server *WebServer) {
	cacheEnvelope(cloudControllerOrigin, server)
	
	request := createResourceRequest(t, token, port, "cloud_controllers")
	
	t.Logf("Check if server response to valid /cloud_controllers request... (expecting status code: %v)", http.StatusOK)
	response, err := client.Do(request)
	
	if err != nil {
		t.Errorf("Error occured while hitting endpoint: %s", err.Error())
	} else if response.StatusCode != http.StatusOK {
		t.Errorf("Expecting status code %v, but received %v", http.StatusOK, response.StatusCode)
	}
}

func trafficControllerEndPointTest(t *testing.T, client *http.Client, token string, port uint32, server *WebServer) {
	cacheEnvelope(trafficControllerOrigin, server)
	
	request := createResourceRequest(t, token, port, "traffic_controllers")
	
	t.Logf("Check if server response to valid /traffic_controllers request... (expecting status code: %v)", http.StatusOK)
	response, err := client.Do(request)
	
	if err != nil {
		t.Errorf("Error occured while hitting endpoint: %s", err.Error())
	} else if response.StatusCode != http.StatusOK {
		t.Errorf("Expecting status code %v, but received %v", http.StatusOK, response.StatusCode)
	}
}

func goRouterEndPointTest(t *testing.T, client *http.Client, token string, port uint32, server *WebServer) {
	cacheEnvelope(goRouterOrigin, server)
	
	request := createResourceRequest(t, token, port, "gorouters")
	
	t.Logf("Check if server response to valid /gorouters request... (expecting status code: %v)", http.StatusOK)
	response, err := client.Do(request)
	
	if err != nil {
		t.Errorf("Error occured while hitting endpoint: %s", err.Error())
	} else if response.StatusCode != http.StatusOK {
		t.Errorf("Expecting status code %v, but received %v", http.StatusOK, response.StatusCode)
	}
}

/** Utility Functions **/
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

func getToken(t *testing.T, client *http.Client, config *nozzleconfiguration.NozzleConfiguration) string {
	t.Log("Requesting token...")
	tokenRequest := createTokenRequest(config.UAAUsername, config.UAAPassword, config.WebServerPort, t)
	response, err := client.Do(tokenRequest)
	
	if err != nil {
		t.Fatalf("Error occured while requsting token: %s", err.Error())
	}
	
	t.Log("Requested token")
	return response.Header.Get(headerTokenKey)
}

func createResourceRequest(t *testing.T, token string, port uint32, endpoint string) *http.Request {
	t.Logf("Creating resource request...")
	request, err := http.NewRequest("GET", fmt.Sprintf("https://localhost:%d/%s", port, endpoint), nil)
	
	if err != nil {
		t.Fatalf("Error occured while formatting request %s: %s", endpoint, err.Error())
	}
	
	request.Header.Add(headerTokenKey, token)
	
	t.Logf("Created resource request")
	return request
}

func cacheEnvelope(originType string, server *WebServer) {
	deployment := 	"deployment"
	eventType :=	events.Envelope_ValueMetric
	job := 			"job"
	index := 		"0"
	ip := 			"127.0.0.1"
	metricName := 	"metric"
	value := 		float64(100)
	unit := 		"unit"
	
	valueMetric := events.ValueMetric {
		Name:	&metricName,
		Value:	&value,
		Unit:	&unit,
	}
	
	envelope := events.Envelope {
		Origin:			&originType,
		EventType:		&eventType,	
		Deployment:		&deployment,
		Job:			&job,
		Index:			&index,
		Ip:				&ip,
		ValueMetric:	&valueMetric,	
	}
	
	server.CacheEnvelope(&envelope)
}
