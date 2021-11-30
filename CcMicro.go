package ccMicro

import (
	"context"
	"github.com/cqu20141693/sip-server/plugins/registry/nacos"
	"github.com/go-playground/validator/v10"
	"go-micro.dev/v4/client"
	"go-micro.dev/v4/logger"
	"go-micro.dev/v4/selector"

	httpClient "github.com/asim/go-micro/plugins/client/http/v4"
	//"github.com/asim/go-micro/plugins/registry/nacos/v3"
	httpServer "github.com/asim/go-micro/plugins/server/http/v4"
	"github.com/cqu20141693/sip-server/config"
	"github.com/gin-gonic/gin"
	"go-micro.dev/v4"
	"go-micro.dev/v4/registry"
	"go-micro.dev/v4/server"
	"strings"
)

type CcAPI interface {
	InitRouteMapper(router *gin.Engine)
}

type NacosNamespaceContextKey struct {
}

var PanicFunc = func() {
	if err := recover(); err != nil {
		logger.Info(err)
		logger.Info("occur panic")
	}
}
var Validate = validator.New()

func CreateRegister() registry.Registry {
	return nacos.NewRegistry(func(options *registry.Options) {
		addr := config.GetStringOrDefault("cc.cloud.nacos.config.server-Addr", "localhost:8848")
		addrs := strings.Split(addr, ",")
		ns := config.GetStringOrDefault("cc.cloud.nacos.config.namespace", "")
		options.Addrs = addrs
		options.Context = context.WithValue(context.Background(), &NacosNamespaceContextKey{}, ns)
	})
}

func CreateService() micro.Service {
	webAddr := config.GetStringOrDefault("server.port", "8081")
	appName := config.GetStringOrDefault("cc.application.name", "go.micro")
	return micro.NewService(
		micro.Address(":"+webAddr),
		micro.Name(appName),
		micro.Registry(CreateRegister()),
	)
}
func CreateServiceWithHttpServer() micro.Service {
	webAddr := config.GetStringOrDefault("server.port", "8080")
	appName := config.GetStringOrDefault("cc.application.name", "go.micro")
	srv := httpServer.NewServer(
		server.Name(appName),
		server.Address(":"+webAddr),
	)

	return micro.NewService(
		micro.Server(srv),
		micro.Name(appName),
		micro.Registry(CreateRegister()),
	)
}

func CreateClient() client.Client {

	s := selector.NewSelector(selector.Registry(CreateRegister()))

	return httpClient.NewClient(client.Selector(s),
		client.ContentType("application/json"))
}
