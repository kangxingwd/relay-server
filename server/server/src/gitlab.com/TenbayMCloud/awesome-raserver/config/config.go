package config

import (
	"fmt"
	log "github.com/cihub/seelog"
	ms "github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"os"
	"strings"
)

func exitf(format string, a ...interface{}) bool {
	fmt.Printf(format, a...)
	os.Exit(1)
	return true
}

func loadParseConfig()(*Settings, error) {
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var t Settings
	if err := ms.Decode(viper.AllSettings(), &t); err != nil {
		return nil, err
	}

	return &t, nil
}

func initSeelog() {
	logger, err := log.LoggerFromConfigAsString(viper.GetString("seelog"))
	_ = err != nil && exitf("init seelog fail %s", err.Error())

	log.ReplaceLogger(logger)
}

func doLoadConfig() {
	cfgFile := viper.GetString("config")
	_ = cfgFile == "" && exitf("missing config")

	fmt.Printf("WEB_CONFIG=%s\n", cfgFile)

	viper.SetConfigFile(cfgFile)

	t, err := loadParseConfig()
	_ = err != nil && exitf("loadParseConfig %s fail %s", cfgFile, err.Error())

	settings = t
}

func InitConfig() {
	viper.SetEnvPrefix("WEB")
	viper.AutomaticEnv() // read in environment variables that match
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetDefault("config", "./cfg.toml")

	doLoadConfig()
	initSeelog()
}

func HotReload() {
	t, err := loadParseConfig()
	if err != nil {
		log.Errorf("hot reload fail. %s", err.Error())
		return
	}

	settings = t

	log.Infof("hot reload config ok")
}

func GetMysql() *Mysql {
	return &settings.Mysql
}
func GetRedis() *Redis {
	return &settings.Redis
}
func GetCommon() *Common {
	return &settings.Common
}
func GetDNS() *DNS {
	return &settings.DNS
}
func GetWebSvr() *WebSvr {
	return &settings.WebSvr
}
