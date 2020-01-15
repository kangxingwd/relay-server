package DnsReloadSrv

import (
	"context"
	logger "github.com/cihub/seelog"
	"gitlab.com/TenbayMCloud/awesome-raserver/common/dns"
	cfg "gitlab.com/TenbayMCloud/awesome-raserver/config"
	"time"
)

type Context = context.Context

type DnsReloadSrv struct {
	ctx Context
	addr   string
}

func NewmDnsReloadSrv(ctx Context) *DnsReloadSrv {
	return &DnsReloadSrv{ctx: ctx}
}

func (mm *DnsReloadSrv) Initialize() error {
	return nil
}

func (mm *DnsReloadSrv) Name() string {
	return "DnsReloadSrv"
}

func (mm *DnsReloadSrv) Run() {
	ticker := time.NewTicker(time.Second *  time.Duration(cfg.GetDNS().DnsConfUpdateTime))
	go func() {
		for {
			select {
			case <-ticker.C:
				//logger.Infof("【DnsReloadSrv】start ConfCheck %v\n", time.Now())
				dns.ConfCheck()
			}
		}
	}()

	<-mm.ctx.Done()
	ticker.Stop()
	logger.Infof("【DnsReloadSrv】stop daemon: %s", mm.Name())
	return
}


