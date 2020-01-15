package handler

import (
	logger "github.com/cihub/seelog"
	"github.com/json-iterator/go"
	"github.com/valyala/fasthttp"
	"gitlab.com/TenbayMCloud/awesome-raserver/common"
	"gitlab.com/TenbayMCloud/awesome-raserver/common/dns"
	cfg "gitlab.com/TenbayMCloud/awesome-raserver/config"
	"gitlab.com/TenbayMCloud/awesome-raserver/model"
	"strconv"
	"strings"
	"time"
)

// 入网
func RelayNetJoin(ctx  *fasthttp.RequestCtx, msg  *common.CliMsg)  {
	var ret common.RetInfo

	// 请求参数检查
	if  msg.DevId == "" || msg.Data.Vendor == "" || msg.Data.DevHostname == "" ||
		msg.Data.DevMac == ""	||  msg.Data.WanIp == "" || common.IsValidIP(msg.Data.WanIp) == false {
		ret.Result = common.RetJsonFormatErr			// 请求参数格式错误或缺少参数
		common.RetJson(ctx, &ret)
		return
	}

	// 检查是否入网
	if model.IsExistRelay(msg.DevId) == true {
		ret.Result = common.RetRepetJoin			//  重复提交入网
		common.RetJson(ctx, &ret)
		return
	}

	// 向dev表添加设备信息
	if err := common.DevAdd(msg, common.RelayFlag); err != nil {
		ret.Result = common.RetServerErr			//	服务器内部错误
		common.RetJson(ctx, &ret)
		logger.Errorf("DevAdd error: %s \n", err.Error())
		return
	}

	// 判断是否有公网
	ret.Data.IpType = common.IsPublicNetwork(ctx, msg)
	if ret.Data.IpType == common.IsNotPublicIp {
		ret.Result = common.RetNoPublicIp			//	没有公网IP
		common.RetJson(ctx, &ret)
		return
	}

	// 添加relay信息到relay表
	if err := relayAdd(msg); err != nil{
		ret.Result = common.RetServerErr			//	服务器内部错误
		common.RetJson(ctx, &ret)
		logger.Errorf("relayAdd error: %s\n", err.Error())
		return
	}

	// 入网成功
	ret.Result = common.RetSuccess
	common.RetJson(ctx, &ret)
}

// 关闭网络
func RelayNetStop(ctx  *fasthttp.RequestCtx, msg  *common.CliMsg)  {
	var ret common.RetInfo
	var err error

	// 请求参数检查
	if  msg.DevId == "" {
		ret.Result = common.RetJsonFormatErr			// 请求参数格式错误或缺少参数
		common.RetJson(ctx, &ret)
		return
	}

	// 关闭relay角色
	if err := common.DelRelayRole(msg.DevId); err != nil {
		ret.Result = common.RetServerErr			//	服务器内部错误
		common.RetJson(ctx, &ret)
		logger.Errorf("DelRelayRole error: %s \n", err.Error())
		return
	}

	// 检测是否有该relay信息
	var relay *model.Relay
	if relay,err = model.GetRelay(msg.DevId); relay.Devid == "" {
		ret.Result = common.RetNoRelayInfo			// 没有relay信息,
		common.RetJson(ctx, &ret)
		return
	}

	// 删除dns映射
	if err := common.DeleteDnsByRelay(relay); err != nil {
		logger.Errorf("relay netstop, deleteDnsByRelay err: %v\n",err)
	}

	// 删除 rcMap表中的映射
	if err := model.DelRcmapByRelay(msg.DevId); err != nil {
		logger.Errorf("rcMapDel error! err: %v\n",err)
	}

	// 删除relay信息
	if err = common.DelRelay(msg.DevId); err != nil{
		logger.Errorf("relayDel error: %v\n", err.Error())
		ret.Result = common.RetServerErr			//	服务器内部错误
		common.RetJson(ctx, &ret)
		return
	}

	// 关网成功
	ret.Result = common.RetSuccess
	common.RetJson(ctx, &ret)
}

// 请求域名
func RequestDomain(ctx  *fasthttp.RequestCtx, msg  *common.CliMsg)  {
	var ret common.RetInfo
	var err error

	// 请求参数检查
	if  msg.DevId == "" || len(msg.Data.McloudMacList) == 0 {
		ret.Result = common.RetJsonFormatErr			// 请求参数格式错误或缺少参数
		common.RetJson(ctx, &ret)
		return
	}

	// 检测是否有该relay信息
	if model.IsExistRelay(msg.DevId) == false {
		ret.Result = common.RetNoRelayInfo			// 没有relay信息,
		common.RetJson(ctx, &ret)
		return
	}

	var mclouds []common.Mcloud
	if mclouds,err = common.GetDomainMapList(msg.Data.McloudMacList); err != nil {
		logger.Errorf("GetDomain error: %v\n", err)
		ret.Result = common.RetServerErr			//	服务器内部错误
		common.RetJson(ctx, &ret)
		return
	}

	if len(mclouds) == 0 {
		ret.Result = common.RetNoDomainInfo
		common.RetJson(ctx, &ret)
		return
	}

	// 请求域名成功
	ret.Data.Mclouds = mclouds
	ret.Result = common.RetSuccess
	common.RetJson(ctx, &ret)
}

// 添加mcloud上报
func AddMcloud(ctx  *fasthttp.RequestCtx, msg  *common.CliMsg)  {

}

// relay心跳
func RelayHeartBeat(ctx  *fasthttp.RequestCtx, msg  *common.CliMsg)  {
	var ret common.RetInfo

	// 请求参数检查
	if  msg.DevId == "" || msg.Data.WanIp == "" || common.IsValidIP(msg.Data.WanIp) == false ||
		msg.Data.RelayCpu.State == "" || msg.Data.RelayCpu.Used == "" ||
		msg.Data.RelayMemory.Used == "" || msg.Data.RelayMemory.Free == "" {
		ret.Result = common.RetJsonFormatErr			// 请求参数格式错误或缺少参数
		common.RetJson(ctx, &ret)
		return
	}

	// 检测是否有该relay信息
	var relay *model.Relay
	if relay,_ = model.GetRelay(msg.DevId); relay.Devid == "" {
		ret.Result = common.RetNoRelayInfo			// 没有relay信息,
		common.RetJson(ctx, &ret)
		return
	}

	// 检查IP是否变化
	var ipChangeFlag = 0
	if msg.Data.WanIp != relay.WanIP {  // 变化
		ipChangeFlag = 1

		// 删除 rcMap表中的映射
		if err := model.DelRcmapByRelay(msg.DevId); err != nil {
			logger.Errorf("rcMapDel error! err: %v\n",err)
		}

		// 判断是否有公网
		ret.Data.IpType = common.IsPublicNetwork(ctx, msg)
		if ret.Data.IpType == common.IsNotPublicIp {
			// 私网
			if err := common.DelRelay(msg.DevId); err != nil {				// 删除relay信息
				logger.Errorf("relayDel error: %s\n", err.Error())
				ret.Result = common.RetServerErr
				return
			}

			/// 删除dns映射
			if err := dns.DnsDelByIp(relay.WanIP); err != nil {
				logger.Errorf(" relay ip change not public, DnsDelByIp err: %v\n",err)
			}

			ret.Result = common.RetNoPublicIp			// 	没有公网ip
			common.RetJson(ctx, &ret)
			return
		}

		// 删除dns映射,不删除本身的mcloud映射
		if err := common.DeleteDnsByRelay(relay); err != nil {
			logger.Errorf(" relay ip change public, deleteDnsByRelay err: %v\n",err)
		}
	}

	// 判断负载
	var overload uint32 = 0		// 设备负载
	result,err := isOverload(msg);
	if err != nil {
		logger.Errorf("isOverload error: %s\n", err.Error())
		ret.Result = common.RetCpuMemFormatErr
		common.RetJson(ctx, &ret)
		return
	}
	if result == true {
		overload = 1
	}

	// 更新relay
	if err := updateReley(msg, ipChangeFlag, overload); err != nil {
		logger.Errorf("updateReley error: %s\n", err.Error())
		ret.Result = common.RetServerErr
		common.RetJson(ctx, &ret)
		return
	}

	ret.Result = common.RetSuccess
	common.RetJson(ctx, &ret)
}

// 功能：	添加relay信息到relay表
// 输入		msg:	请求的JSON信息
// 返回： 	error:	错误信息
func relayAdd(msg  *common.CliMsg)  (error) {
	extInfo := common.RelayExt{
		Cpu:	msg.Data.RelayCpu,
		Mem:	msg.Data.RelayMemory,
	}

	json_iterator := jsoniter.ConfigCompatibleWithStandardLibrary
	ext, err := json_iterator.Marshal(extInfo)
	if err != nil {
		logger.Info("ext jsonMarshal error! err: %v\n", err.Error())
		ext = []byte{}
	}

	relay := &model.Relay{
		Devid:			msg.DevId,
		PublicIP:		msg.Data.WanIp,
		WanIP:			msg.Data.WanIp,
		HeartTime:		time.Now().Unix(),
		ConnState:		1,
		FailedCountTime:time.Now().Unix(),
		Ext:			string(ext),
	}
	if err := relay.Add(); err != nil {
		logger.Errorf("relay add failed! err: %v\n", err)
		return err
	}
	logger.Infof("add relay:\n %+v\n", *relay)
	return nil
}

// 功能：	判断设备是否负载
// 输入		msg:	请求的JSON信息
// 返回： 	bool:	是否负载
func isOverload(msg  *common.CliMsg) (bool, error) {
	var mem = msg.Data.RelayMemory
	var err error

	var cpuUsed string
	cpuUsed = strings.Split(msg.Data.RelayCpu.Used, ".")[0]
	relayCpuPercent,err := strconv.Atoi(cpuUsed)
	if err != nil {
		logger.Errorf("cpu used format error! need: \"90\" or \"90.1\", but given %s",
			msg.Data.RelayCpu.Used)
		return false, err
	}

	var relayUseMem int64
	if relayUseMem, err = strconv.ParseInt(mem.Used, 10, 64); err != nil {
		logger.Errorf("mem used string to int64 error! mem.Free = %s\n", mem.Free)
		return false, err
	}

	var relayFreeMem int64
	if relayFreeMem, err = strconv.ParseInt(mem.Free, 10, 64); err != nil {
		logger.Errorf("mem free string to int64 error! mem.Free = %s\n", mem.Free)
		return false, err
	}
	relayMemPercent := int(100 * relayUseMem/(relayUseMem + relayFreeMem))

	if relayCpuPercent >= cfg.GetCommon().MaxRelayCpuPercent ||
		relayMemPercent >= cfg.GetCommon().MaxRelayMemPercent {
			//logger.Errorf("relayCpuPercent = %v, relayMemPercent=%v, " +
			//	"cfg.GetCommon().MaxRelayCpuPercent = %v, cfg.GetCommon().MaxRelayMemPercent = %v", relayCpuPercent, relayMemPercent,
			//	cfg.GetCommon().MaxRelayCpuPercent, cfg.GetCommon().MaxRelayMemPercent)
		return true, nil
	}
	return false, nil
}

// 功能：	更新relay信息--
// 输入		msg:			请求的JSON信息
//			ipChangeFlag: 	ip变化标记(1 变化，0 没有变化)
//			overLoad:		负载 (1 负载, 0 没有负载)
// 返回： 	error:			错误信息
func updateReley(msg  *common.CliMsg, ipChangeFlag int, overLoad uint32) error {
	extInfo := common.RelayExt{
		Cpu:	msg.Data.RelayCpu,
		Mem:	msg.Data.RelayMemory,
	}

	json_iterator := jsoniter.ConfigCompatibleWithStandardLibrary
	ext, err := json_iterator.Marshal(extInfo)
	if err != nil {
		logger.Info("ext jsonMarshal error! err: %v\n", err.Error())
		ext = []byte{}
	}
	
	var relay *model.Relay
	if relay,err = model.GetRelay(msg.DevId); err != nil {
		logger.Errorf("updateReley failed! dev_id = %s, err: %v\n", msg.DevId, err)
		return err
	}

	relay.ConnState = 1
	relay.McloudNum = msg.Data.McloudNum
	relay.HeartTime = time.Now().Unix()
	relay.OverLoad = overLoad
	relay.Ext = string(ext)

	if ipChangeFlag == 1 { 			// ip 改变
		relay.PublicIP = msg.Data.WanIp
		relay.WanIP = msg.Data.WanIp
		relay.McloudNum = 0
		relay.RelayAllocNum = 0
		relay.ConnFaildNum = 0
		relay.FailedCountTime = 0
	}

	if err := relay.Update(); err!= nil {
		logger.Errorf("update relay sql failed! dev_id = %s, err: %v\n", msg.DevId, err)
		return err
	}
	return nil
}
