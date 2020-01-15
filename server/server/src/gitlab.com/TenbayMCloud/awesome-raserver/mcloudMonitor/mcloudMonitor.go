package mcloudMonitor

import (
	"context"
	logger "github.com/cihub/seelog"
	"gitlab.com/TenbayMCloud/awesome-raserver/common"
	cfg "gitlab.com/TenbayMCloud/awesome-raserver/config"
	"gitlab.com/TenbayMCloud/awesome-raserver/model"
	"time"
)

type Context = context.Context

type McloudMonitor struct {
	ctx Context
	addr   string
}

func NewmCloudMonitor(ctx Context) *McloudMonitor {
	return &McloudMonitor{ctx: ctx}
}

func (mm *McloudMonitor) Initialize() error {


	return nil
}

func (mm *McloudMonitor) Name() string {
	return "mcloudMonitor"
}

func (mm *McloudMonitor) Run() {
	ticker := time.NewTicker(time.Second *  time.Duration(cfg.GetCommon().HeartTimemcloud))
	go func() {
		for {
			select {
			case <-ticker.C:
				logger.Infof("【mcloudMonitor】start mcloudCheck %v\n", time.Now())
				mcloudCheck()
			}
		}
	}()

	<-mm.ctx.Done()
	ticker.Stop()
	logger.Infof("stop daemon: %s", mm.Name())
	return
}

func mcloudCheck()  {
	var mclouds *[]model.Mcloud
	var err error

	time.Sleep(time.Second *  time.Duration(cfg.GetCommon().HeartTimemcloud))
	if mclouds,err = model.GetAllMclouds(); err != nil {
		logger.Errorf("【mcloudMonitor】sql exec error! error: %s\n", err.Error())
		return
	}

	for _,mcloud := range *mclouds {
		checkTimeout(&mcloud)		// 超时检测
	}
	return
}

func checkTimeout(mcloud *model.Mcloud)  {
	var interval int64 = time.Now().Unix() - mcloud.HeartTime		// 与上次心跳间隔时间

	switch {
	case interval >= 0 && interval < cfg.GetCommon().HTMcloudTimeOut:	// 正常
		// 不做任何操作
	case interval >= cfg.GetCommon().HTMcloudTimeOut && interval < cfg.GetCommon().HTMcloudDelete:	//	断开
		// 更新 conn_state字段
		if mcloud.ConnState != 0 {
			model.UpdateMcloudState(mcloud.Devid, 0)
		}
	case interval >= cfg.GetCommon().HTMcloudDelete:		// 注销
		// 删除mcloud信息
		common.DelMcloud(mcloud.Devid, "")
	default:
		logger.Errorf("【mcloudMonitor】Calc interval error! interval = %v\n", interval)
	}
}


