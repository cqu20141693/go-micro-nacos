package main

import (
	ccMicro "github.com/cqu20141693/sip-server"
	"github.com/cqu20141693/sip-server/client"
	"github.com/cqu20141693/sip-server/server/handler"
	"github.com/gin-gonic/gin"
	"go-micro.dev/v4/logger"
	"go-micro.dev/v4/server"
)

func main() {

	service := ccMicro.CreateServiceWithHttpServer()
	service.Init()
	configRouter(service.Server())
	// Run service
	if err := service.Run(); err != nil {
		logger.Fatal(err)
	}
}

func configRouter(server server.Server) {
	gin.SetMode(gin.DebugMode)
	router := gin.New()
	router.Use(gin.Recovery())

	service := handler.NewSipService(client.NewSipClient())
	service.InitRouteMapper(router)
	hd := server.NewHandler(router)
	if err := server.Handle(hd); err != nil {
		logger.Fatal(err)
	}
}
