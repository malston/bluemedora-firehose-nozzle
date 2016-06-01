// Copyright (c) 2016 Blue Medora, Inc. All rights reserved.
// This file is subject to the terms and conditions defined in the included file 'LICENSE.txt'.

package webserver

import (
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/BlueMedora/bluemedora-firehose-nozzle/nozzleconfiguration"
	"github.com/BlueMedora/bluemedora-firehose-nozzle/webtoken"
	"github.com/cloudfoundry/gosteno"
	"github.com/cloudfoundry/sonde-go/events"
)

//Webserver Constants
const (
	DefaultCertLocation 	= "./certs/cert.pem"
	DefaultKeyLocation  	= "./certs/key.pem"
	headerUsernameKey   	= "username"
	headerPasswordKey   	= "password"
	headerTokenKey      	= "token"
	
	metronAgentOrigin			= "MetronAgent"
    syslogDrainBinderOrigin		= "syslog_drain_binder"
    tpsWatcherOrigin			= "tps_watcher"
    tpsListenerOrigin			= "tps_listener"
    stagerOrigin				= "stager"
    sshProxyOrigin				= "ssh-proxy"
    senderOrigin				= "sender"
    routeEmitterOrigin			= "route_emitters"
    repOrigin					= "rep"
    receptorOrigin				= "receptor"
    nsyncListenerOrigin			= "nsync_listener"
    nsyncBulkerOrigin			= "nsync_bulker"
    gardenLinuxOrigin			= "garden-linux"
    fileServerOrigin			= "file_server"
    fetcherOrigin				= "fetcher"
    convergerOrigin				= "converger"
    ccUploaderOrigin			= "cc_uploader"
    bbsOrigin					= "bbs"
	auctioneerOrigin			= "auctioneer"
	etcdOrigin					= "etcd"
	dopplerServerOrigin			= "DopplerServer"
	cloudControllerOrigin		= "cc"
	trafficControllerOrigin		= "LoggregatorTrafficController"
	goRouterOrigin				= "gorouter"
)

//WebServer REST endpoint for sending data
type WebServer struct {
	logger *gosteno.Logger
	mutext sync.Mutex
	config *nozzleconfiguration.NozzleConfiguration
	tokens map[string]*webtoken.Token //Maps token string to token object
	
	cache  map[string]map[string]Resource
}

//New creates a new WebServer
func New(config *nozzleconfiguration.NozzleConfiguration, logger *gosteno.Logger) *WebServer {
	webserver := WebServer{
		logger: logger,
		config: config,
		tokens: make(map[string]*webtoken.Token),
		cache: 	make(map[string]map[string]Resource),
	}

	webserver.logger.Info("Registering handlers")
	//setup http handlers
	http.HandleFunc("/token", webserver.tokenHandler)
	http.HandleFunc("/metron_agents", webserver.metronAgentsHandler)
	http.HandleFunc("/syslog_drains", webserver.syslogDrainBindersHandler)
	http.HandleFunc("/tps_watchers", webserver.tpsWatcherHandler)
	http.HandleFunc("/tps_listeners", webserver.tpsListenersHandler)
	http.HandleFunc("/stagers", webserver.stagerHandler)
	http.HandleFunc("/ssh_proxies", webserver.sshProxyHandler)
	http.HandleFunc("/senders", webserver.senderHandler)
	http.HandleFunc("/route_emitters", webserver.routeEmitterHandler)
	http.HandleFunc("/reps", webserver.repHandler)
	http.HandleFunc("/receptors", webserver.receptorHandler)
	http.HandleFunc("/nsync_listeners", webserver.nsyncListenerHandler)
	http.HandleFunc("/nsync_bulkers", webserver.nsyncBulkerHandler)
	http.HandleFunc("/garden_linuxs", webserver.gardenLinuxHandler)
	http.HandleFunc("/file_servers", webserver.fileServersHandler)
	http.HandleFunc("/fetchers", webserver.fetcherHandler)
	http.HandleFunc("/convergers", webserver.convergerHandler)
	http.HandleFunc("/cc_uploaders", webserver.ccUploaderHandler)
	http.HandleFunc("/bbs", webserver.bbsHandler)
	http.HandleFunc("/auctioneers", webserver.auctioneerHandler)
	http.HandleFunc("/etcds", webserver.etcdsHandler)
	http.HandleFunc("/doppler_servers", webserver.dopplerServersHandler)
	http.HandleFunc("/cloud_controllers", webserver.cloudControllersHandler)
	http.HandleFunc("/traffic_controllers", webserver.trafficControllersHandler)
	http.HandleFunc("/gorouters", webserver.gorouterHandler)

	return &webserver
}

//Start starts webserver listening
func (webserver *WebServer) Start(keyLocation string, certLocation string) <-chan error {
	webserver.logger.Infof("Start listening on port %v", webserver.config.WebServerPort)
	errors := make(chan error, 1)
	go func() {
		defer close(errors)
		errors <- http.ListenAndServeTLS(fmt.Sprintf(":%v", webserver.config.WebServerPort), certLocation, keyLocation, nil)
	}()
	return errors
}

//TokenTimeout is a callback for when a token timesout to remove
func (webserver *WebServer) TokenTimeout(token *webtoken.Token) {
	webserver.mutext.Lock()
	webserver.logger.Debugf("Removing token %s", token.TokenValue)
	delete(webserver.tokens, token.TokenValue)
	webserver.mutext.Unlock()
}

/**Handlers**/
func (webserver *WebServer) tokenHandler(w http.ResponseWriter, r *http.Request) {
	webserver.logger.Info("Received /token request")
	if r.Method == http.MethodGet {
		username := r.Header.Get(headerUsernameKey)
		password := r.Header.Get(headerPasswordKey)

		//Check for username and password
		if username == "" || password == "" {
			webserver.logger.Debug("No username or password in header")
			w.WriteHeader(http.StatusBadRequest)
			io.WriteString(w, "username and/or password not found in header")
		} else {
			//Check validity of username and password
			if username != webserver.config.UAAUsername || password != webserver.config.UAAPassword {
				webserver.logger.Debugf("Wrong username and password for user %s", username)
				w.WriteHeader(http.StatusUnauthorized)
				io.WriteString(w, "Invalid Username and/or Password")
			} else {
				//Successful login
				token := webtoken.New(webserver.TokenTimeout)

				webserver.mutext.Lock()
				webserver.tokens[token.TokenValue] = token
				webserver.mutext.Unlock()

				w.Header().Set(headerTokenKey, token.TokenValue)
				w.WriteHeader(http.StatusOK)

				webserver.logger.Debugf("Successful login generated token <%s>", token.TokenValue)
			}
		}
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		io.WriteString(w, fmt.Sprintf("/token does not support %s http methods", r.Method))
	}
}

func (webserver *WebServer) metronAgentsHandler(w http.ResponseWriter, r *http.Request) {
	webserver.logger.Info("Received /metron_agents request")
}

func (webserver *WebServer) syslogDrainBindersHandler(w http.ResponseWriter, r *http.Request) {
	webserver.logger.Info("Received /syslog_drains request")
}

func (webserver *WebServer) tpsWatcherHandler(w http.ResponseWriter, r *http.Request) {
	webserver.logger.Info("Received /tps_watchers request")
}

func (webserver *WebServer) tpsListenersHandler(w http.ResponseWriter, r *http.Request) {
	webserver.logger.Info("Received /tps_listeners request")
}

func (webserver *WebServer) stagerHandler(w http.ResponseWriter, r *http.Request) {
	webserver.logger.Info("Received /stagers request")
}

func (webserver *WebServer) sshProxyHandler(w http.ResponseWriter, r *http.Request) {
	webserver.logger.Info("Received /ssh_proxies request")
}

func (webserver *WebServer) senderHandler(w http.ResponseWriter, r *http.Request) {
	webserver.logger.Info("Received /senders request")
}

func (webserver *WebServer) routeEmitterHandler(w http.ResponseWriter, r *http.Request) {
	webserver.logger.Info("Received /route_emitters request")
}

func (webserver *WebServer) repHandler(w http.ResponseWriter, r *http.Request) {
	webserver.logger.Info("Received /reps request")
}

func (webserver *WebServer) receptorHandler(w http.ResponseWriter, r *http.Request) {
	webserver.logger.Info("Received /receptors request")
}

func (webserver *WebServer) nsyncListenerHandler(w http.ResponseWriter, r *http.Request) {
	webserver.logger.Info("Received /nsync_listeners request")
}

func (webserver *WebServer) nsyncBulkerHandler(w http.ResponseWriter, r *http.Request) {
	webserver.logger.Info("Received /nsync_bulkers request")
}

func (webserver *WebServer) gardenLinuxHandler(w http.ResponseWriter, r *http.Request) {
	webserver.logger.Info("Received /garden_linuxs request")
}

func (webserver *WebServer) fileServersHandler(w http.ResponseWriter, r *http.Request) {
	webserver.logger.Info("Received /file_servers request")
}

func (webserver *WebServer) fetcherHandler(w http.ResponseWriter, r *http.Request) {
	webserver.logger.Info("Received /fetchers request")
}

func (webserver *WebServer) convergerHandler(w http.ResponseWriter, r *http.Request) {
	webserver.logger.Info("Received /convergers request")
}

func (webserver *WebServer) ccUploaderHandler(w http.ResponseWriter, r *http.Request) {
	webserver.logger.Info("Received /cc_uploaders request")
}

func (webserver *WebServer) bbsHandler(w http.ResponseWriter, r *http.Request) {
	webserver.logger.Info("Received /bbs request")
}

func (webserver *WebServer) auctioneerHandler(w http.ResponseWriter, r *http.Request) {
	webserver.logger.Info("Received /auctioneers request")
}

func (webserver *WebServer) etcdsHandler(w http.ResponseWriter, r *http.Request) {
	webserver.logger.Info("Received /etcds request")
}

func (webserver *WebServer) dopplerServersHandler(w http.ResponseWriter, r *http.Request) {
	webserver.logger.Info("Received /doppler_servers request")
}

func (webserver *WebServer) cloudControllersHandler(w http.ResponseWriter, r *http.Request) {
	webserver.logger.Info("Received /cloud_controllers request")
}

func (webserver *WebServer) trafficControllersHandler(w http.ResponseWriter, r *http.Request) {
	webserver.logger.Info("Received /traffic_controllers request")
}

func (webserver *WebServer) gorouterHandler(w http.ResponseWriter, r *http.Request) {
	webserver.logger.Info("Received /gorouters request")
}

/**Cache Logic**/

//CacheEnvelope caches envelope by origin
func (webserver *WebServer) CacheEnvelope(envelope *events.Envelope) {
	webserver.mutext.Lock()
	defer webserver.mutext.Unlock()
	
	key := createEnvelopeKey(envelope)
	webserver.logger.Debugf("Caching envelope origin %s with key %s", envelope.GetOrigin(), key)
	
	//Find origin map
	var resourceCache map[string]Resource
	
	if value, ok := webserver.cache[envelope.GetOrigin()]; ok {
		resourceCache = value
	} else {
		resourceCache = make(map[string]Resource)
		webserver.cache[envelope.GetOrigin()] = resourceCache
	}
	
	//Check to see if resource exists in origin map
	var resource Resource
	if value, ok := resourceCache[key]; ok {
		resource = value
	} else {
		resource = Resource {
			Deployment:		envelope.GetDeployment(),
			Job:			envelope.GetJob(),
			Index:			envelope.GetIndex(),
			IP:				envelope.GetIp(),
			ValueMetrics:	make(map[string]float64),
			CounterMetrics:	make(map[string]float64),
		}
	}
	
	addMetric(envelope, resource.ValueMetrics, resource.CounterMetrics, webserver.logger)
	resourceCache[key] = resource
}

//ClearCache clears out cache for server
func (webserver *WebServer) ClearCache() {
	webserver.logger.Info("Flushing Cache")
	webserver.mutext.Lock()
	defer webserver.mutext.Unlock()
	
	webserver.cache = make(map[string]map[string]Resource)
}

