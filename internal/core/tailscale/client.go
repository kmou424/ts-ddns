package tailscale

import (
	"context"
	"fmt"
	"github.com/gookit/slog"
	"github.com/kmou424/ero"
	"github.com/kmou424/ts-ddns/internal/global"
	"github.com/kmou424/ts-ddns/pkgs/util"
	"golang.org/x/oauth2/clientcredentials"
	"net/http"
	"net/url"
)

type Client struct {
	client *http.Client

	baseURL string
}

func NewClient() *Client {
	slog.Info("creating tailscale client")
	tailscaleConfig := global.Config.Tailscale
	client := &http.Client{}

	baseURL := tailscaleConfig.BaseURL
	u, err := url.Parse(baseURL)
	if err != nil {
		panic(ero.New("invalid tailscale base url"))
	}

	switch tailscaleConfig.Method {
	case "token":
		// todo: implement token auth
	case "oauth2":
		u.Path = "/api/v2/oauth/token"
		var oauthConfig = &clientcredentials.Config{
			ClientID:     tailscaleConfig.ClientId,
			ClientSecret: tailscaleConfig.ClientSecret,
			TokenURL:     u.String(),
		}
		client = oauthConfig.Client(context.Background())
	}

	u.Path = ""
	return &Client{
		client:  client,
		baseURL: u.String(),
	}
}

func (c *Client) GetDevices() ([]Device, error) {
	u, _ := url.Parse(c.baseURL)
	u.Path = fmt.Sprintf("/api/v2/tailnet/%s/devices", global.Config.Tailscale.Tailnet)

	resp, err := c.client.Get(u.String())
	if err != nil {
		return nil, ero.Wrap(err, "failed to get devices")
	}

	marshalled, err := util.MarshalResp[struct {
		Devices []Device `json:"devices"`
	}](resp)
	if err != nil {
		return nil, ero.Wrap(err, "failed to unmarshal response")
	}

	return marshalled.Devices, nil
}
