package routers

import (
	"my-gmail-server/routers/api"
	v1 "my-gmail-server/routers/api/v1"
	"my-gmail-server/services/auth_service"
	"my-gmail-server/settings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// enable cors for my angular spa
func CorsHandler() gin.HandlerFunc {
	config := cors.DefaultConfig()
	config.AllowCredentials = true
	config.AllowOrigins = []string{"http://localhost:4200"}
	return cors.New(config)
}

func InitRouter() *gin.Engine {
	e := gin.Default()

	e.Static("/assets", "./assets")
	e.Use(CorsHandler())
	e.Use(auth_service.Session(settings.AuthSettings.SessionName))

	apiv1 := e.Group("/api/v1")

	auth := e.Group("/auth")
	{
		// initialise settings for google oauth2.0
		auth.GET("/login", api.Login)
		auth.GET("/logout", api.Logout)

		// google cloud redirects me to here, and register the user in the database
		auth.GET("/callback", api.AuthCallback)
	}

	// handles form submissions
	jobs := apiv1.Group("/job")
	{
		jobs.POST("/scan", v1.CreateScanJob)
		jobs.POST("/trash", v1.CreateTrashJob)
	}

	// rest endpoints
	{
		apiv1.GET("/users", api.SetAuthState(), v1.ListUsers)
		apiv1.GET("/users/:id", v1.ListUserById)
		apiv1.GET("/users/:id/info/scan", v1.ListActiveScanJobInfo)
		apiv1.GET("/users/:id/info/trash", v1.ListActiveTrashJobInfo)

		apiv1.GET("/health", v1.CheckHealth)
	}

	return e
}
