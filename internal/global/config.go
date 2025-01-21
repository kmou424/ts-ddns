package global

import (
	"github.com/kmou424/ts-ddns/internal/config"
)

var Config *config.Config

func initConfig() {
	Config = config.AutoLoadConfig()
}
