package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/iwannay/jiaweb/utils/file"
)

type (
	Config struct {
		App     *AppConfig
		Session *SessionConfig
		Route   *RouteConfig
		Group   *GroupConfig
	}

	AppConfig struct {
		LogPath     string
		EnableLog   bool
		RunMode     string
		PprofPort   int
		EnablePprof bool
	}

	SessionConfig struct {
	}

	RouteConfig struct {
	}

	GroupConfig struct {
	}
)

const (
	errorPrefix = "Jiaweb:config"
)

func defaultAppConfig() *AppConfig {
	return &AppConfig{}
}

func defaultSessionConfig() *SessionConfig {
	return &SessionConfig{}
}

func defaultRouteConfig() *RouteConfig {
	return &RouteConfig{}
}

func defaultGroupConfig() *GroupConfig {
	return &GroupConfig{}
}

// MustInitConfig 初始化配置文件否则panic
func MustInitConfig(configfile string, configType string) *Config {
	conf, err := InitConfig(configfile, configType)
	if err != nil {
		panic(err)
	}
	return conf
}

// InitConfig 初始化配置文件
func InitConfig(configFile string, configType string) (*Config, error) {
	if !file.Exist(configFile) {
		configFile = file.GetCurrentDirectory() + "/" + configFile
		if !file.Exist(configFile) {
			configFile = file.GetCurrentDirectory() + "/config/" + configFile
			if !file.Exist(configFile) {
				return nil, fmt.Errorf("%s file %s not exists", errorPrefix, configFile)
			}
		}
	}

	bytes, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("%s config file [%s] cannot parse %s", errorPrefix, configFile, err)
	}

	var config Config

	switch configType {
	case "json":
		err = readFromJson(bytes, &config)
	default:
		err = readFromJson(bytes, &config)

	}

	if err != nil {
		return nil, fmt.Errorf("%s config file [%s] cannot parse %s", errorPrefix, configFile, err)
	}

	if config.App == nil {
		config.App = defaultAppConfig()
	}

	if config.Group == nil {
		config.Group = defaultGroupConfig()
	}

	if config.Route == nil {
		config.Route = defaultRouteConfig()
	}

	if config.Session == nil {
		config.Session = defaultSessionConfig()
	}

	return &config, nil

}

func readFromJson(bytes []byte, v interface{}) error {
	return json.Unmarshal(bytes, v)
}
