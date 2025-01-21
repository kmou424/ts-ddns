package service

import (
	"github.com/kmou424/ts-ddns/internal/core/dns"
	"github.com/kmou424/ts-ddns/internal/core/tailscale"
	"github.com/kmou424/ts-ddns/pkgs/typed"
)

// Context state-less context
type Context struct {
	tsClient *tailscale.Client
	provider dns.IProvider
	extra    typed.Map[int]
}
