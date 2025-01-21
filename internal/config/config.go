package config

import (
	"github.com/kmou424/ero"
	"github.com/kmou424/ts-ddns/pkgs/util"
	"github.com/mcuadros/go-defaults"
)

type Config struct {
	ManagedHost string `toml:"managed_host"`
	SLD         string `toml:"-"`
	IPv4        bool   `toml:"ipv4" default:"true"`
	IPv6        bool   `toml:"ipv6" default:"false"`
	Interval    int    `toml:"interval" default:"300"`

	Tailscale struct {
		Method       string `toml:"method" default:"token"`
		Tailnet      string `toml:"tailnet" default:""`
		BaseURL      string `toml:"base_url" default:"https://api.tailscale.com/"`
		APIKey       string `toml:"api_key"`
		ClientId     string `toml:"client_id"`
		ClientSecret string `toml:"client_secret"`
	} `toml:"tailscale"`

	DNS struct {
		Provider string            `toml:"provider" default:"cloudflare"`
		Params   map[string]string `toml:"params"`
	} `toml:"dns"`
}

func (c *Config) setDefault() {
	defaults.SetDefaults(c)

	var err error
	c.SLD, err = util.GetSLD(c.ManagedHost)
	if err != nil {
		panic(ero.Wrap(err, "failed to get SLD from managed host"))
	}
}

func (c *Config) validate() {

}
