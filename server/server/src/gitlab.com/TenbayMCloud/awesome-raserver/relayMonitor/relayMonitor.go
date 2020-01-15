package relayMonitor

import (
	"context"
	logger "github.com/cihub/seelog"
	"gitlab.com/TenbayMCloud/awesome-raserver/common"
	cfg "gitlab.com/TenbayMCloud/awesome-raserver/config"
	"gitlab.com/TenbayMCloud/awesome-raserver/model"
	"time"
)

type Context = context.Context

type RelayMonitor struct {
	ctx Context
	addr   string
}

func NewRelayMonitor(ctx Context) *RelayMonitor {
	return &RelayMonitor{ctx: ctx}
}

func (rm *RelayMonitor) Initialize() error {


	return nil
}

func (rm *RelayMonitor) Name() string {
	return "relayMonitor"
}

func (rm *RelayMonitor) Run() {
	ticker := time.NewTicker(time.Second * time.Duration(cfg.GetCommon().HeartTimeRelay))
	go func() {
		for {
			select {
			case <-ticker.C:
				logger.Infof("【RelayMonitor】start relayCheck %v\n", time.Now())
				relayCheck()
			}
		}
	}()

	<-rm.ctx.Done()
	ticker.Stop()
	logger.Infof("stop daemon: %s", rm.Name())
	return
}

func relayCheck()  {
	relays,err := model.GetAllRelays();
	if err != nil {
		logger.Errorf("【RelayMonitor】sql exec error! error: %s\n", err.Error())
		return
	}

	for _,relay := range *relays {
		// 超时检测
		isDelete := handleTimeout(&relay)

		// 黑名单检测
		if isDelete != true && cfg.GetCommon().RelayForbidSwitch == 1 {
			handleForbid(&relay)
		}
	}
}

func handleTimeout(relay *model.Relay) bool {
	var interval int64 = time.Now().Unix() - relay.HeartTime		// 与上次心跳间隔时间

	if relay.HeartTime == 0 {
		logger.Errorf( "【RelayMonitor】heartTime is 0! select data error!\n")
		return false
	}
	
	switch {
	case interval >= 0 && interval < cfg.GetCommon().HTRelayTimeOut:	// 正常
		// 不做任何操作
	case interval >= cfg.GetCommon().HTRelayTimeOut && interval < cfg.GetCommon().HTRelayDelete:	//	断开
		// 更新 conn_state字段
		if relay.ConnState != 0 {
			model.UpdateReleyState(relay.Devid, 0)
		}
	case interval >= cfg.GetCommon().HTRelayDelete:		// 注销
		// 删除relay信息
		common.DelRelay(relay.Devid)
		// 删除dns映射
		if err := common.DeleteDnsByRelay(relay); err != nil {
			logger.Errorf("【RelayMonitor】 deleteDnsByRelay err: %v\n",err)
		}
		return true
	default:
		logger.Errorf("【RelayMonitor】Calc interval error! interval = %v\n", interval)
	}
	return false
}

func handleForbid(relay *model.Relay)  {

	var interval int64 = time.Now().Unix() - relay.FailedCountTime

	if relay.ForbidState == 1 {  		// 已经被封禁
		if interval > cfg.GetCommon().ForbidTime {
			// 封禁时间到， 重置统计信息
			model.UpdateForbid(relay.Devid, 0)
		}
	}else {		// 未被封禁
		switch  {
		case interval < cfg.GetCommon().ForbidCountTime:
			// 不检测
		case interval >= cfg.GetCommon().ForbidCountTime && interval < cfg.GetCommon().ForbidCountTime + 60:
			if relay.RelayAllocNum > 1 && float32(relay.ConnFaildNum/relay.RelayAllocNum) > cfg.GetCommon().ForbidFailedRatio {
				// 失败率 大于 配置的失败率， 设置黑名单
				model.UpdateForbid(relay.Devid, 1)
			}else {
				// 	失败率小与配置的失败率， 重置统计信息
				model.UpdateForbid(relay.Devid, 0)
			}
		default:
			// 超过统计时间， 重置统计信息
			model.UpdateForbid(relay.Devid, 0)
		}
	}
}

