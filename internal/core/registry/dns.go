package registry

import (
	"github.com/gookit/slog"
	"github.com/kmou424/ero"
	"github.com/kmou424/ts-ddns/internal/core/dns"
	"github.com/kmou424/ts-ddns/internal/global"
	"sync/atomic"
)

var gDNSProviders = make(map[string]dns.IProvider)
var gDNSProvidersInitStatus = make(map[string]*atomic.Bool)

func RegisterDNSProvider(name string, getProvider func() dns.IProvider) {
	if global.RegistryDisabled.Load() {
		return
	}
	provider := getProvider()
	if provider == nil {
		panic(ero.New("provider is nil"))
	}
	gDNSProviders[name] = provider
	gDNSProvidersInitStatus[name] = new(atomic.Bool)
	slog.Infof("dns provider registered [%s]", name)
}

func GetDNSProvider() dns.IProvider {
	providerName := global.Config.DNS.Provider
	providerParams := global.Config.DNS.Params

	provider, ok := gDNSProviders[providerName]
	if !ok {
		panic(ero.Newf("provider %s not found", providerName))
	}
	if !gDNSProvidersInitStatus[providerName].Load() {
		slog.Infof("initializing dns provider [%s]...", providerName)
		err := provider.Init(providerParams)
		if err != nil {
			panic(ero.Wrap(err, "failed to initialize provider"))
		}
	}
	return provider
}
