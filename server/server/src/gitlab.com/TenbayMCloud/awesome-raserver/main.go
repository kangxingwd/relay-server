package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	S "syscall"
	"time"

	log "github.com/cihub/seelog"
	"gitlab.com/TenbayMCloud/awesome-raserver/common"
	cfg "gitlab.com/TenbayMCloud/awesome-raserver/config"
	DnsReloadSrv "gitlab.com/TenbayMCloud/awesome-raserver/dnsReloadSrv"
	"gitlab.com/TenbayMCloud/awesome-raserver/mcloudMonitor"
	"gitlab.com/TenbayMCloud/awesome-raserver/relayMonitor"
	"gitlab.com/TenbayMCloud/awesome-raserver/websrv"
)

func init() {
	cfg.InitConfig()
	log.Infof("GetWebSvr	%+v\n", *cfg.GetWebSvr())
	log.Infof("GetDNS	%+v\n", *cfg.GetDNS())
	log.Infof("GetCommon	%+v\n", *cfg.GetCommon())
	log.Infof("GetMysql	%+v\n", *cfg.GetMysql())
	log.Infof("GetRedis	%+v\n", *cfg.GetRedis())
}

type Daemon = common.Daemon

func setupSignal(cancel context.CancelFunc) {
	sigChan := make(chan os.Signal)

	signal.Notify(sigChan, S.SIGINT, S.SIGTERM, S.SIGQUIT /*, S.SIGUSR1, S.SIGUSR2*/)

	go func() {
		for s := range sigChan {
			log.Infof("recv signal: %s", s)

			switch s {
			case S.SIGINT, S.SIGTERM, S.SIGQUIT:
				cancel()
				/*case S.SIGUSR1:
					cfg.HotReload()
				case S.SIGUSR2:*/
			default:
			}
		}
	}()
}

func mustInitDaemon(daemons []Daemon) error {
	for _, d := range daemons {
		err := d.Initialize()
		if err != nil {
			log.Errorf("Initialize error! err: %s", err.Error())
			return err
		}
	}
	return nil
}

func main() {
	log.Infof("-------------main")
	ctx, cancel := context.WithCancel(context.Background())
	setupSignal(cancel)

	daemons := []common.Daemon{
		websrv.NewWebSrv(ctx),
		relayMonitor.NewRelayMonitor(ctx),
		mcloudMonitor.NewmCloudMonitor(ctx),
		DnsReloadSrv.NewmDnsReloadSrv(ctx),
	}

	log.Infof("mustInitDaemon(daemons)")
	if err := mustInitDaemon(daemons); err != nil {
		return
	}

	var wg sync.WaitGroup
	for _, d := range daemons {
		wg.Add(1)

		d := d
		go func() {
			d.Run()
			log.Infof("daemon: %s exit", d.Name())
			wg.Done()
		}()
		time.Sleep(time.Second / 10)
	}

	<-ctx.Done()

	log.Infof("recv quit signal. wait")
	wg.Wait()
	log.Infof("all routines end. exit")
	log.Flush()
}
