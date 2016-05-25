/**Copyright Blue Medora Inc. 2016**/

package webtoken

import (
    "math/rand"
)

var (
    tokenLength = 15
    tokenRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
)

//GenerateTokenString creates a token string to be used in webserver
func GenerateTokenString() string {
    token := make([]rune, tokenLength)
    for i:= range token {
        token[i] = tokenRunes[rand.Intn(len(tokenRunes))]
    }
    
    return string(token)
}
