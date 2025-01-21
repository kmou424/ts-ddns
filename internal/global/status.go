package global

import (
	"github.com/gookit/goutil"
	"github.com/gookit/goutil/sysutil"
	"sync/atomic"
)

var RegistryDisabled atomic.Bool
var Debug atomic.Bool

func initStatus() {
	ddnsDebug, _ := goutil.ToBool(sysutil.Getenv("TS_DDNS_DEBUG", "false"))
	Debug.Store(ddnsDebug)
}
