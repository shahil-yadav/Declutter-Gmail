package middleware

import (
	"my-gmail-server/services/auth_service"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func SetAuthState() gin.HandlerFunc {
	return func(c *gin.Context) {
		rvalue := auth_service.RandToken()
		s := sessions.Default(c)
		s.Set(auth_service.SessionState, rvalue)
		s.Save()
	}
}
