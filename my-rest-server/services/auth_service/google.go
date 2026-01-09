// https://github.com/zalando/gin-oauth2/blob/master/google/google.go

// Package google provides you access to Google's OAuth2
// infrastructure. The implementation is based on this blog post:
// http://skarlso.github.io/2016/06/12/google-signin-with-go/
package auth_service

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/gob"
	"encoding/json"
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/golang/glog"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func init() {
	gob.Register(oauth2.Token{})
}

// Credentials stores google client-ids.
type Credentials struct {
	ClientID     string `json:"clientid"`
	ClientSecret string `json:"secret"`
}

const (
	SessionState = "state"
	SessionToken = "ginoauth_google_token"

	GinCtxToken = "token"
)

var (
	conf  *oauth2.Config
	store sessions.Store
)

var loginURL string

func RandToken() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		glog.Fatalf("[Gin-OAuth] Failed to read rand: %v", err)
	}
	return base64.StdEncoding.EncodeToString(b)
}

// SetupGoogleConfig the authorization path
func SetupGoogleConfig(redirectURL, credFile string, scopes []string, secret []byte) {
	store = cookie.NewStore(secret)

	var c Credentials
	// todo: add reading from config instead of json
	// Localise credentials.json into conf/app.ini
	file, err := os.ReadFile(credFile)
	if err != nil {
		glog.Fatalf("[Gin-OAuth] File error: %v", err)
	}
	if err := json.Unmarshal(file, &c); err != nil {
		glog.Fatalf("[Gin-OAuth] Failed to unmarshal client credentials: %v", err)
	}

	conf = &oauth2.Config{
		ClientID:     c.ClientID,
		ClientSecret: c.ClientSecret,
		RedirectURL:  redirectURL,
		Scopes:       scopes,
		Endpoint:     google.Endpoint,
	}
}

func GetLoginURL(state string) string {
	// get refresh token
	return conf.AuthCodeURL(state, oauth2.AccessTypeOffline)
}
