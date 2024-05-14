package config

import (
	"github.com/pmh-only/spoti2wall/utils"
	"gopkg.in/ini.v1"
)

var GlobalConfig *ini.File

func InitConfig() {
	GlobalConfig, _ = ini.Load(utils.GetConfigPath())
}

func GetClientId(defaultValue string) string {
	clientID := GlobalConfig.Section("").Key("client_id").String()
	if clientID == "" {
		return defaultValue
	} else {
		return clientID
	}
}

func GetClientSecret(defaultValue string) string {
	secret := GlobalConfig.Section("").Key("client_secret").String()
	if secret == "" {
		return defaultValue
	} else {
		return secret
	}
}
