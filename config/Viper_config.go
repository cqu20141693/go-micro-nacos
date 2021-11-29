package config

import (
	"errors"
	"github.com/spf13/viper"
	"go-micro.dev/v4/logger"
	"log"
)

func init() {
	ReadLocalConfig()
	NacosInit()
}

func ReadLocalConfig() {
	// 读取本地配置
	viper.SetConfigName("bootstrap.yml")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./resource")
	viper.AddConfigPath("/etc/resource")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			logger.Info("Config file not found; ignore error if desired")

		} else {
			logger.Info("Config file was found but another error was produced")
		}
		log.Fatal(err)
	}
}

func GetAppConfigName(name, active string) (string, error) {
	if name == "" {
		return "", errors.New("cc.application.name not config")
	} else if active == "" {
		return "", errors.New("cc.profiles.active not config")
	}
	return name + "-" + active + "." + LocalNacosConfig.FileExtension, nil
}

func GetStringOrDefault(key, defaultVal string) string {
	str := viper.GetString(key)
	if str == "" {
		return defaultVal
	}
	return str
}
