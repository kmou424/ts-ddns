package service

import (
	"fmt"
	"github.com/gookit/slog"
	"github.com/kmou424/ero"
	"github.com/kmou424/ts-ddns/internal/core/tailscale"
	"github.com/kmou424/ts-ddns/internal/global"
	"github.com/kmou424/ts-ddns/pkgs/util"
)

const (
	attrDeviceIPSetMap = iota
	attrManagedRecords
)

func looping(ctx *Context) {
	defer func() {
		if e := recover(); e != nil {
			err := e.(error)
			slog.Error("failed to update DNS records: %s", func() string {
				if global.Debug.Load() && ero.IsEro(err) {
					trace := ero.AllTrace(err, true)
					return trace
				} else {
					return fmt.Sprintf("%v", err)
				}
			}())
		}
	}()
	slog.Info("start DDNS loop")
	beforeLoop(ctx)
	onLoop(ctx)
	afterLoop(ctx)
	slog.Info("done DDNS loop")
}

func beforeLoop(ctx *Context) {
	slog.Info("getting devices from tailscale")
	devices, err := ctx.tsClient.GetDevices()
	if err != nil {
		panic(err)
	}

	var deviceIPSetMap = make(map[string]tailscale.DeviceIPSet)
	for _, device := range devices {
		ipSet := device.ToIPSet()
		domain := util.JoinDomains(ipSet.Name, global.Config.ManagedHost)
		deviceIPSetMap[domain] = ipSet
	}
	ctx.extra.Set(attrDeviceIPSetMap, deviceIPSetMap)
}

func onLoop(ctx *Context) {
	slog.Info("getting DNS records from provider")
	records, err := ctx.provider.GetRecords()
	if err != nil {
		panic(err)
	}
	managedRecords := toManagedRecords(records)
	ctx.extra.Set(attrManagedRecords, managedRecords)
}

func afterLoop(ctx *Context) {
	managedRecords := ctx.extra.MustGet(attrManagedRecords).(ManagedRecords)
	cleanRecords(ctx, managedRecords)
	createRecords(ctx, managedRecords)
}

func cleanRecords(ctx *Context, managedRecords ManagedRecords) {
	ipSetMap := ctx.extra.MustGet(attrDeviceIPSetMap).(map[string]tailscale.DeviceIPSet)
	deleteList := managedRecords.deleteList(ipSetMap)
	if len(deleteList) > 0 {
		for _, record := range deleteList {
			slog.Infof("delete DNS record: <%s %s:%s>", record.Type, record.Domain, record.IP)
			err := ctx.provider.DeleteRecord(record)
			if err != nil {
				slog.Errorf("failed to delete DNS record: %s", err)
			}
		}
	} else {
		slog.Info("no DNS records need to delete")
	}
}

func createRecords(ctx *Context, managedRecords ManagedRecords) {
	ipSetMap := ctx.extra.MustGet(attrDeviceIPSetMap).(map[string]tailscale.DeviceIPSet)
	createList := managedRecords.createList(ipSetMap)
	if len(createList) > 0 {
		for _, record := range createList {
			slog.Infof("create DNS record: <%s %s:%s>", record.Type, record.Domain, record.IP)
			err := ctx.provider.CreateRecord(record)
			if err != nil {
				slog.Errorf("failed to insert DNS record: %s", err)
			}
		}
	} else {
		slog.Info("no DNS records need to create")
	}
}
