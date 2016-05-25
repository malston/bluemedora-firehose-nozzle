/**Copyright Blue Medora Inc. 2016**/

package webtoken

import (
    "time"
    "sync"
)

var (
    tokenTimeout = 60
)

//InvalidTokenError signals invalid token usage
type InvalidTokenError struct {
    s string
}

func (e *InvalidTokenError) Error() string {
    return e.s
}

//Token token used for webserver communication
type Token struct {
    TokenValue                  string
    validToken                  bool
    tokenUsedSinceLastTimout    bool
    timoutTicker                *time.Ticker
    mux                         sync.Mutex
}

//New creates a new token
func New() *Token {
    newToken := Token {
        TokenValue:                 GenerateTokenString(),
        validToken:                 true,
        tokenUsedSinceLastTimout:   false,
        timoutTicker:               time.NewTicker(time.Duration(tokenTimeout) * time.Second),
    }
    
    go newToken.startTimeout()
    
    return &newToken
}

func (token *Token) startTimeout() {
    select {
        case <-token.timoutTicker.C:
            token.mux.Lock()
            if !token.tokenUsedSinceLastTimout {
                token.validToken = false
                return
            }
             
            token.tokenUsedSinceLastTimout = false;
            token.mux.Unlock()
    }
}

//IsTokenValid validity of token
func (token *Token) IsTokenValid() bool {
    token.mux.Lock()
    defer token.mux.Unlock()
    return token.validToken
}

//UseToken marks token as used 
func (token *Token) UseToken() error {
    token.mux.Lock()
    defer token.mux.Unlock()
    
    if token.validToken {
        token.tokenUsedSinceLastTimout = true
    } else {
        return &InvalidTokenError{"Attempt to use invalid token"}
    }
    
    return nil
}