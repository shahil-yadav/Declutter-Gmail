package routers

import (
	"my-gmail-server/pkg/gintemplrenderer"
	"my-gmail-server/routers/api"
	v1 "my-gmail-server/routers/api/v1"
	"my-gmail-server/routers/web"
	"my-gmail-server/services/auth_service"
	"my-gmail-server/settings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	e := gin.Default()

	// enable cors for angular
	e.Use(cors.Default())

	gintemplrenderer.Setup(e)

	// serving static files eg. image, video at 127.0.0.1/assets
	e.Static("/assets", "./assets")

	// initialise session for route group
	e.Use(auth_service.Session(settings.AuthSettings.SessionName))

	v1Prefixed := e.Group("/v1")
	auth := e.Group("/auth")
	{
		// initialise settings for google oauth2.0
		auth.GET("/login", api.Login)
		auth.GET("/logout", api.Logout)

		// google cloud redirects me to here, and register the user in the database
		auth.GET("/callback", api.AuthCallback)
	}

	// renders html routes
	html := e.Group("/")
	{
		html.GET("/", web.Index)
		html.GET("/:user", web.Dashboard)
	}

	// handles form submissions
	forms := v1Prefixed.Group("/job")
	{
		forms.POST("/scan", v1.CreateScanJob)
		forms.POST("/trash", v1.CreateTrashJob)
	}

	// rest endpoints
	{
		v1Prefixed.GET("/users", v1.ListUsers)
		v1Prefixed.GET("/users/:id", v1.ListUserById)
		v1Prefixed.GET("/users/:id/info/scan", v1.ListActiveScanJobInfo)
		v1Prefixed.GET("/users/:id/info/trash", v1.ListActiveTrashJobInfo)
	}

	return e
}
