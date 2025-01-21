package service

import (
	"github.com/gookit/goutil/arrutil"
	"github.com/kmou424/ts-ddns/internal/core/dns"
	"github.com/kmou424/ts-ddns/internal/core/tailscale"
	"github.com/kmou424/ts-ddns/internal/global"
	"strings"
)

func toManagedRecords(records []dns.Record) ManagedRecords {
	var managedRecords ManagedRecords
	managedHost := global.Config.ManagedHost
	for _, record := range records {
		if strings.HasSuffix(record.Domain, managedHost) {
			managedRecords = append(managedRecords, record)
		}
	}
	return managedRecords
}

type ManagedRecords []dns.Record

func (r *ManagedRecords) deleteList(ipSetMap map[string]tailscale.DeviceIPSet) []dns.Record {
	toReserve := make([]dns.Record, 0)
	toDeletes := make([]dns.Record, 0)
	for _, record := range *r {
		switch record.Type {
		case dns.RecordTypeA:
			if !global.Config.IPv4 {
				toDeletes = append(toDeletes, record)
				continue
			}
		case dns.RecordTypeAAAA:
			if !global.Config.IPv6 {
				toDeletes = append(toDeletes, record)
				continue
			}
		}
		if ipSet, ok := ipSetMap[record.Domain]; !ok {
			toDeletes = append(toDeletes, record)
			continue
		} else {
			switch record.Type {
			case dns.RecordTypeA:
				if !arrutil.SliceHas(ipSet.IPv4, record.IP) {
					toDeletes = append(toDeletes, record)
					continue
				}
			case dns.RecordTypeAAAA:
				if !arrutil.SliceHas(ipSet.IPv6, record.IP) {
					toDeletes = append(toDeletes, record)
					continue
				}
			}
		}
		toReserve = append(toReserve, record)
	}
	*r = toReserve
	return toDeletes
}

func (r *ManagedRecords) createList(ipSetMap map[string]tailscale.DeviceIPSet) []dns.Record {
	toCreates := make([]dns.Record, 0)
	knownRecordsMap := make(map[string][]dns.Record)
	for _, record := range *r {
		recordTag := record.Domain + ":" + record.IP
		if _, ok := knownRecordsMap[recordTag]; !ok {
			knownRecordsMap[recordTag] = make([]dns.Record, 0)
		}
		knownRecordsMap[recordTag] = append(knownRecordsMap[recordTag], record)
	}
	for domain, ipSet := range ipSetMap {
		if global.Config.IPv4 {
			for _, ip := range ipSet.IPv4 {
				recordTag := domain + ":" + ip
				if _, ok := knownRecordsMap[recordTag]; ok {
					continue
				}
				record := dns.NewEmptyRecord()
				record.Type = dns.RecordTypeA
				record.Domain = domain
				record.IP = ip
				toCreates = append(toCreates, record)
			}
		}
		if global.Config.IPv6 {
			for _, ip := range ipSet.IPv6 {
				recordTag := domain + ":" + ip
				if _, ok := knownRecordsMap[recordTag]; ok {
					continue
				}
				record := dns.NewEmptyRecord()
				record.Type = dns.RecordTypeAAAA
				record.Domain = domain
				record.IP = ip
				toCreates = append(toCreates, record)
			}
		}
	}
	return toCreates
}
