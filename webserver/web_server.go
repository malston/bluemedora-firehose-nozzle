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
)

//Webserver Constants
const (
	DefaultCertLocation = "./certs/cert.pem"
	DefaultKeyLocation  = "./certs/key.pem"
	headerUsernameKey   = "username"
	headerPasswordKey   = "password"
	headerTokenKey      = "token"
)

//WebServer REST endpoint for sending data
type WebServer struct {
	logger *gosteno.Logger
	mutext sync.Mutex
	config *nozzleconfiguration.NozzleConfiguration
	tokens map[string]*webtoken.Token //Maps token string to token object
}

//New creates a new WebServer
func New(config *nozzleconfiguration.NozzleConfiguration, logger *gosteno.Logger) *WebServer {
	webserver := WebServer{
		logger: logger,
		config: config,
		tokens: make(map[string]*webtoken.Token),
	}

	webserver.logger.Info("Registering handlers")
	//setup http handlers
	http.HandleFunc("/token", webserver.tokenHandler)
	http.HandleFunc("/metron_agents", webserver.metronAgentsHandler)
	http.HandleFunc("/syslog_drains", webserver.syslogDrainBindersHandler)
	http.HandleFunc("/etcds", webserver.etcdsHandler)
	http.HandleFunc("/doppler_servers", webserver.dopplerServersHandler)
	http.HandleFunc("/diegos", webserver.diegosHandler)
	http.HandleFunc("/cloud_controllers", webserver.cloudControllersHandler)
	http.HandleFunc("/traffic_controllers", webserver.trafficControllersHandler)

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
			if username != webserver.config.UAAUsername && password != webserver.config.UAAPassword {
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

func (webserver *WebServer) etcdsHandler(w http.ResponseWriter, r *http.Request) {
	webserver.logger.Info("Received /etcds request")
}

func (webserver *WebServer) dopplerServersHandler(w http.ResponseWriter, r *http.Request) {
	webserver.logger.Info("Received /doppler_servers request")
}

func (webserver *WebServer) diegosHandler(w http.ResponseWriter, r *http.Request) {
	webserver.logger.Info("Received /diegos request")
}

func (webserver *WebServer) cloudControllersHandler(w http.ResponseWriter, r *http.Request) {
	webserver.logger.Info("Received /cloud_controllers request")
}

func (webserver *WebServer) trafficControllersHandler(w http.ResponseWriter, r *http.Request) {
	webserver.logger.Info("Received /traffic_controllers request")
}
