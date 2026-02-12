package api

import (
	"my-gmail-server/pkg/app"
	"my-gmail-server/pkg/e"
	"my-gmail-server/services/auth_service"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func AuthCallback(c *gin.Context) {
	auth_service.RegisterUser(c)
}

func Logout(c *gin.Context) {
	appG := app.Gin{C: c}
	session := sessions.Default(c)
	session.Clear()

	// forgot to save the session while clearing it, dumb me
	session.Save()
	appG.Response(http.StatusOK, e.SUCCESS, gin.H{"message": "you are successfully logged out"})
}

// middleware
func SetAuthState() gin.HandlerFunc {
	return func(c *gin.Context) {
		rvalue := auth_service.RandToken()
		s := sessions.Default(c)
		s.Set(auth_service.SessionState, rvalue)
		s.Save()
	}
}

func Login(c *gin.Context) {
	appG := app.Gin{C: c}

	redirect := c.Query("redirect")
	stateValue := auth_service.RandToken()

	// create/get the session from request context
	session := sessions.Default(c)
	session.Set(auth_service.SessionState, stateValue)
	session.Save()

	switch redirect {
	case "true":
		c.Redirect(http.StatusPermanentRedirect, auth_service.GetLoginURL(stateValue))

	case "false":
		appG.Response(http.StatusOK, e.SUCCESS, auth_service.GetLoginURL(stateValue))

	default:
		// need the redirect query to function
		// becoz i wanted to test out this new appG.Response() method
		appG.Response(http.StatusBadRequest, e.STATUS_BAD_REQUEST, "Invalid params")
	}
}
