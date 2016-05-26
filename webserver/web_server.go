/**Copyright Blue Medora Inc. 2016**/

package webserver

import (
    "net/http"
    "io"
    "sync"
    "fmt"
    
    "github.com/cloudfoundry/gosteno"
    "github.com/BlueMedora/bluemedora-firehose-nozzle/nozzleconfiguration"
    "github.com/BlueMedora/bluemedora-firehose-nozzle/webtoken"
)

//Webserver Constants
const (
    DefaultCertLocation = "./certs/cert.perm"
    DefaultKeyLocation = "./certs/key.perm"
    headerUsernameKey = "username"
    headerPasswordKey = "password"
    headerTokenKey = "token"
)

//WebServer REST endpoint for sending data
type WebServer struct {
    logger  *gosteno.Logger
    mutext  sync.Mutex
    config  *nozzleconfiguration.NozzleConfiguration
    tokens  map[string]*webtoken.Token //Maps token string to token object
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

//handle each resource metric request