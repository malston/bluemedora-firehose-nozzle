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
    
    //setup http handlers
    http.HandleFunc("/token", webserver.tokenHandler)
    
    return &webserver
}

//Start starts webserver listening. Should be run in goroutine
func (webserver *WebServer) Start(keyLocation string, certLocation string) error {
    err := http.ListenAndServeTLS(fmt.Sprintf(":%v", webserver.config.WebServerPort), certLocation, keyLocation, nil)
    return err
}

//TokenTimeout is a callback for when a token timesout to remove
func (webserver *WebServer) TokenTimeout(token *webtoken.Token) {
    webserver.mutext.Lock()
    delete(webserver.tokens, token.TokenValue)
    webserver.mutext.Unlock()
}

/**Handlers**/
func (webserver *WebServer) tokenHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodGet {
        username := r.Header.Get(headerUsernameKey)
        password := r.Header.Get(headerPasswordKey)
        
        //Check for username and password
        if username == "" || password == "" {
            w.WriteHeader(http.StatusBadRequest)
            io.WriteString(w, "username and/or password not found in header")
        } else {
            //Check validity of username and password
            if username != webserver.config.UAAUsername && password != webserver.config.UAAPassword {
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
            }
        }
    } else {
        w.WriteHeader(http.StatusMethodNotAllowed)
        io.WriteString(w, fmt.Sprintf("/token does not support %s http methods", r.Method))
    }
}

//handle login post
//handle each resource metric request
//handle token timeout