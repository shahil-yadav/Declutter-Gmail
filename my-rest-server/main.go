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
	"golang.org/x/sync/errgroup"
)

func init() {
	settings.Setup()
	models.Setup()
	jobs_service.Setup()
	auth_service.Setup()
}

func main() {
	gin.SetMode(settings.ServerSetting.RunMode)

	var g errgroup.Group
	g.Go(func() error {
		r := routers.InitRouter()
		endPoint := fmt.Sprintf(":%d", settings.ServerSetting.HttpPort)

		return r.Run(endPoint)
	})

	g.Go(func() error {
		return jobs_service.RunWebUi(9000)
	})

	//fix: if mux server has errors, errgroup.Group doesn't report
	if err := g.Wait(); err != nil {
		log.Fatal("failed to run either gin or mux", err)
	}
}
