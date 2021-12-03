package config

import (
	"encoding/json"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/spf13/viper"
	"go-micro.dev/v4/logger"
	"log"
	"strconv"
	"strings"
)

var ConfigClient config_client.IConfigClient
var Err error
var DefaultGroup = "DEFAULT_GROUP"

var sc []constant.ServerConfig

func createServerConfig(ipAddr string, port uint64) constant.ServerConfig {
	return *constant.NewServerConfig(
		ipAddr,
		port,
		constant.WithScheme("http"),
		constant.WithContextPath("/nacos"),
	)
}

// 创建clientConfig
var cc constant.ClientConfig

type config struct {
	ServerAddr      string `json:"server-addr"`
	FileExtension   string `json:"file-extension"`
	Namespace       string `json:"namespace"`
	ClusterName     string `json:"cluster-name"`
	Group           string `json:"group"`
	RegisterEnabled bool   `json:"register-enabled"`
}

func newDefaultConfig() config {
	return config{ServerAddr: "localhost:8848", FileExtension: "yml", Namespace: "", ClusterName: "DEFAULT", Group: "DEFAULT_GROUP", RegisterEnabled: true}
}

var LocalNacosConfig = newDefaultConfig()

func NacosInit() {

	nacosConfig := viper.GetStringMap("cc.cloud.nacos.config")
	marshal, _ := json.Marshal(nacosConfig)
	json.Unmarshal(marshal, &LocalNacosConfig)
	splits := strings.Split(LocalNacosConfig.ServerAddr, ",")
	for i := range splits {
		host := strings.Split(splits[i], ":")
		if len(host) == 2 {
			port, err := strconv.Atoi(host[1])
			if err != nil {
				log.Fatal(err)
			}
			sc = append(sc, createServerConfig(host[0], uint64(port)))
		}
	}
	cc = constant.ClientConfig{
		NamespaceId: LocalNacosConfig.Namespace,
	}

	ConfigClient, Err = clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		},
	)
	if Err != nil {
		panic(Err)
	}
	if Err != nil {
		panic(Err)
	}
	ReadNacosConfig()
}

func ReadNacosConfig() {
	// 读取nacos配置
	name := viper.GetString("cc.application.name")
	active := viper.GetString("cc.profiles.active")
	appConfigName, err1 := GetAppConfigName(name, active)
	if err1 != nil {
		log.Fatal(err1)
	}
	content, err := GetConfigByDataId(appConfigName)
	if err != nil {
		logger.Info("getConfig error dataId=" + appConfigName)
	}
	err = viper.MergeConfig(strings.NewReader(content))
	if err != nil {
		logger.Error("Viper Failed to resolve configuration content=" + content)
	}
}

func GetConfigByDataId(dataId string) (content string, err error) {

	return GetConfig(dataId, DefaultGroup)
}

func GetConfig(dataId, group string) (content string, err error) {

	content, err = ConfigClient.GetConfig(vo.ConfigParam{
		DataId: dataId,
		Group:  group,
	})
	return
}

type ConfigListenHandler func(namespace, group, dataId, data string)
