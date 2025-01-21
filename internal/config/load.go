package config

import (
	"github.com/gookit/goutil/fsutil"
	"github.com/gookit/goutil/sysutil"
	"github.com/kmou424/ero"
	"github.com/pelletier/go-toml/v2"
	"os"
	"path/filepath"
)

var paths = []string{
	"./config.toml",
	filepath.Join(sysutil.HomeDir(), ".config/ts-ddns/config.toml"),
	"/etc/ts-ddns/config.toml",
}

func AutoLoadConfig() *Config {
	for _, path := range paths {
		if !fsutil.FileExists(path) {
			continue
		}
		if config, err := loadConfig(path); err != nil {
			panic(ero.Wrap(err, "failed to load config file: %s", path))
		} else {
			return config
		}
	}
	panic(ero.New("no config file found"))
	return nil
}

func loadConfig(path string) (*Config, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, ero.Wrap(err, "failed to read config file: %s", path)
	}

	config := &Config{}

	err = toml.Unmarshal(content, config)
	if err != nil {
		return nil, ero.Wrap(err, "failed to parse config file: %s", path)
	}

	config.setDefault()
	config.validate()

	return config, nil
}
