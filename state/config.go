package state

import (
	"github.com/pmh-only/spoti2wall/utils"
	"gopkg.in/ini.v1"
)

var GlobalConfig *ini.File

func InitConfig() {
	GlobalConfig, _ = ini.Load(utils.GetConfigPath())
}
