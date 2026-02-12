package main

import (
	"fmt"
	"log"
	"my-gmail-server/models"
	"my-gmail-server/routers"
	"my-gmail-server/services/auth_service"
	"my-gmail-server/services/jobs_service"
	"my-gmail-server/settings"

	"github.com/gin-gonic/gin"
)

func init() {
	settings.Setup()
	models.Setup()
	jobs_service.Setup()
	auth_service.Setup()
}

func main() {
	gin.SetMode(settings.ServerSetting.RunMode)

	r := routers.InitRouter()
	port := fmt.Sprintf(":%d", settings.ServerSetting.HttpPort)

	if err := r.Run(port); err != nil {
		log.Fatal(err)
	}
}
