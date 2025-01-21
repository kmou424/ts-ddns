package tailscale

import (
	"github.com/kmou424/ts-ddns/pkgs/util"
	"strings"
	"time"
)

type Device struct {
	Addresses                 []string  `json:"addresses"`
	Id                        string    `json:"id"`
	NodeId                    string    `json:"nodeId"`
	User                      string    `json:"user"`
	TailNetAddr               string    `json:"name"`
	Hostname                  string    `json:"hostname"`
	ClientVersion             string    `json:"clientVersion"`
	UpdateAvailable           bool      `json:"updateAvailable"`
	Os                        string    `json:"os"`
	Created                   time.Time `json:"created"`
	LastSeen                  time.Time `json:"lastSeen"`
	KeyExpiryDisabled         bool      `json:"keyExpiryDisabled"`
	Expires                   time.Time `json:"expires"`
	Authorized                bool      `json:"authorized"`
	IsExternal                bool      `json:"isExternal"`
	MachineKey                string    `json:"machineKey"`
	NodeKey                   string    `json:"nodeKey"`
	BlocksIncomingConnections bool      `json:"blocksIncomingConnections"`
	EnabledRoutes             []string  `json:"enabledRoutes"`
	AdvertisedRoutes          []string  `json:"advertisedRoutes"`
	ClientConnectivity        struct {
		Endpoints []string `json:"endpoints"`
		Latency   map[string]map[string]struct {
			Preferred bool    `json:"preferred"`
			LatencyMs float64 `json:"latencyMs"`
		} `json:"latency"`
		MappingVariesByDestIP bool `json:"mappingVariesByDestIP"`
		ClientSupports        struct {
			HairPinning bool `json:"hairPinning"`
			Ipv6        bool `json:"ipv6"`
			Pcp         bool `json:"pcp"`
			Pmp         bool `json:"pmp"`
			Udp         bool `json:"udp"`
			Upnp        bool `json:"upnp"`
		} `json:"clientSupports"`
	} `json:"clientConnectivity"`
	Tags             []string `json:"tags"`
	TailnetLockError string   `json:"tailnetLockError"`
	TailnetLockKey   string   `json:"tailnetLockKey"`
	PostureIdentity  struct {
		SerialNumbers []string `json:"serialNumbers"`
	} `json:"postureIdentity"`
}

type DeviceIPSet struct {
	Name string   `json:"name"`
	IPv4 []string `json:"ipv4"`
	IPv6 []string `json:"ipv6"`
}

func (d *Device) GetName() string {
	return strings.Split(d.TailNetAddr, ".")[0]
}

func (d *Device) ToIPSet() DeviceIPSet {
	ipSet := DeviceIPSet{
		Name: d.GetName(),
		IPv4: make([]string, 0),
		IPv6: make([]string, 0),
	}

	for _, addr := range d.Addresses {
		if util.IsIPv4(addr) {
			ipSet.IPv4 = append(ipSet.IPv4, addr)
		} else if util.IsIPv6(addr) {
			ipSet.IPv6 = append(ipSet.IPv6, addr)
		}
	}

	return ipSet
}
