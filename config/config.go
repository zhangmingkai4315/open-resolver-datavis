package config

import (
	"github.com/BurntSushi/toml"
)

var appConfig *AppConfig

// AppConfig will include all config items for use later
type AppConfig struct {
	GlobalConfig   globalConfig   `toml:"global"`
	DatabaseConfig databaseConfig `toml:"database"`
}

type globalConfig struct {
	Listen string `toml:"listen"`
	Worker int    `tomal:"worker"`
}

type databaseConfig struct {
	Type string `toml:"type"`
	URL  string `toml:"url"`
}

// NewAppConfig will read file and return a app config object
func NewAppConfig(path string) (*AppConfig, error) {
	var config AppConfig
	if _, err := toml.DecodeFile(path, &config); err != nil {
		return nil, err
	}
	appConfig = &config
	return &config, nil
}

// GetAppConfig return config object for application
func GetAppConfig() *AppConfig {
	if appConfig == nil {
		return nil
	}
	return appConfig
}
