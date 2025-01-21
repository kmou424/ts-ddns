package service

import (
	"github.com/gookit/slog"
	"github.com/kmou424/ts-ddns/internal/core/registry"
	"github.com/kmou424/ts-ddns/internal/core/tailscale"
	"github.com/kmou424/ts-ddns/internal/global"
	"github.com/kmou424/ts-ddns/pkgs/typed"
	"time"
)

func Run() {
	slog.Info("DDNS service is starting")
	provider := registry.GetDNSProvider()
	tsClient := tailscale.NewClient()
	slog.Info("DDNS service is now alive")
	for {
		ctx := &Context{
			tsClient: tsClient,
			provider: provider,
			extra:    typed.NewMap[int](),
		}

		looping(ctx)

		duration := time.Duration(global.Config.Interval) * time.Second
		slog.Infof(
			"the next update will be at: %s",
			time.Now().Add(duration).Format("2006-01-02 15:04:05"),
		)
		time.Sleep(duration)
	}
}
