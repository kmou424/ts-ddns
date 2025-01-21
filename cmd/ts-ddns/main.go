package main

import (
	"fmt"
	"github.com/gookit/slog"
	"github.com/kmou424/ero"
	_ "github.com/kmou424/ts-ddns/internal/core/dns/providers"
	"github.com/kmou424/ts-ddns/internal/global"
	"github.com/kmou424/ts-ddns/internal/service"
)

func handleEroPanic() {
	if r := recover(); r != nil {
		err := r.(error)
		if global.Debug.Load() && ero.IsEro(err) {
			trace := ero.AllTrace(err, true)
			slog.Error(trace)
		} else {
			slog.Error(fmt.Sprintf("panic: %v", err))
		}
	}
}

func main() {
	defer handleEroPanic()
	global.RegistryDisabled.Store(true)
	service.Run()
}
