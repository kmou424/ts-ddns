package providers

import (
	"context"
	"github.com/cloudflare/cloudflare-go"
	"github.com/kmou424/ero"
	"github.com/kmou424/ts-ddns/internal/core/dns"
	"github.com/kmou424/ts-ddns/internal/core/registry"
	"github.com/kmou424/ts-ddns/internal/global"
	"github.com/kmou424/ts-ddns/pkgs/typed"
	"github.com/samber/lo"
	"sync/atomic"
)

const (
	ParamCloudflareAPIToken = "api_token"
)

type CloudflareProvider struct {
	initialized atomic.Bool
	api         *cloudflare.API
	zoneID      string
}

func (c *CloudflareProvider) Init(params map[string]string) error {
	if c.initialized.Load() {
		return nil
	}
	apiToken, ok := params[ParamCloudflareAPIToken]
	if !ok {
		return ero.Newf("missing parameter %s", ParamCloudflareAPIToken)
	}
	var err error
	c.api, err = cloudflare.NewWithAPIToken(apiToken)
	if err != nil {
		return ero.Wrap(err, "failed to create cloudflare API client")
	}
	c.zoneID, err = c.api.ZoneIDByName(global.Config.SLD)
	if err != nil {
		return ero.Wrap(err, "failed to get zone ID for domain %s", global.Config.SLD)
	}
	return nil
}

func (c *CloudflareProvider) GetRecords() ([]dns.Record, error) {
	dnsRecords, _, err := c.api.ListDNSRecords(
		context.Background(),
		cloudflare.ZoneIdentifier(c.zoneID),
		cloudflare.ListDNSRecordsParams{},
	)
	if err != nil {
		return nil, ero.Wrap(err, "failed to get DNS dnsRecords")
	}
	records := lo.Map[cloudflare.DNSRecord, dns.Record](dnsRecords, func(record cloudflare.DNSRecord, _ int) dns.Record {
		extra := typed.NewMap[string]()
		return dns.Record{
			ID:      record.ID,
			Type:    record.Type,
			Domain:  record.Name,
			IP:      record.Content,
			TTL:     record.TTL,
			Comment: record.Comment,
			Extra:   extra,
		}
	})
	return records, nil
}

func (c *CloudflareProvider) DeleteRecord(record dns.Record) error {
	err := c.api.DeleteDNSRecord(
		context.Background(),
		cloudflare.ZoneIdentifier(c.zoneID),
		record.ID,
	)
	if err != nil {
		return ero.Wrap(err, "failed to delete DNS record %s", record.ID)
	}
	return nil
}

func (c *CloudflareProvider) CreateRecord(record dns.Record) error {
	_, err := c.api.CreateDNSRecord(
		context.Background(),
		cloudflare.ZoneIdentifier(c.zoneID),
		cloudflare.CreateDNSRecordParams{
			Type:    record.Type,
			Name:    record.Domain,
			Content: record.IP,
			TTL:     record.TTL,
			Proxied: typed.Pointer(false),
			Comment: record.Comment,
		},
	)
	if err != nil {
		return ero.Wrap(err, "failed to create DNS record for %s", record.Domain)
	}
	return nil
}

func (c *CloudflareProvider) UpdateRecord(record dns.Record) error {
	_, err := c.api.UpdateDNSRecord(
		context.Background(),
		cloudflare.ZoneIdentifier(c.zoneID),
		cloudflare.UpdateDNSRecordParams{
			ID:      record.ID,
			Type:    record.Type,
			Name:    record.Domain,
			Content: record.IP,
			TTL:     record.TTL,
			Proxied: typed.Pointer(false),
			Comment: typed.Pointer(record.Comment),
		},
	)
	if err != nil {
		return ero.Wrap(err, "failed to update DNS record %s", record.ID)
	}
	return nil
}

func init() {
	registry.RegisterDNSProvider("cloudflare", func() dns.IProvider {
		return &CloudflareProvider{}
	})
}
