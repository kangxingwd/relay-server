package handler

import (
	logger "github.com/cihub/seelog"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/valyala/fasthttp"
	"gitlab.com/TenbayMCloud/awesome-raserver/common"
	"gitlab.com/TenbayMCloud/awesome-raserver/common/dns"
	cfg "gitlab.com/TenbayMCloud/awesome-raserver/config"
	"gitlab.com/TenbayMCloud/awesome-raserver/model"
	"time"
)

// 入网
func McloudNetJoin(ctx  *fasthttp.RequestCtx, msg  *common.CliMsg)  {
	var ret common.RetInfo
	var err error

	// 请求参数检查
	if msg.DevId == "" || msg.Data.Vendor == "" || msg.Data.DevHostname == "" ||
		msg.Data.DevMac == ""|| msg.Data.WanIp == "" || common.IsValidIP(msg.Data.WanIp) == false {
		ret.Result = common.RetJsonFormatErr			// 请求参数格式错误或缺少参数
		common.RetJson(ctx, &ret)
		return
	}

	// 检查是否已经入网
	if model.IsExistMcloud(msg.DevId) == true {
		ret.Result = common.RetRepetJoin			// 重复入网
		common.RetJson(ctx, &ret)
		return
	}

	// 向dev表添加设备信息
	if err = common.DevAdd(msg, common.McloudFlag); err != nil {
		ret.Result = common.RetServerErr			//	服务器内部错误
		common.RetJson(ctx, &ret)
		logger.Errorf("DevAdd error: %s \n", err.Error())
		return
	}

	// 判断是否有公网
	ipType := common.IsPublicNetwork(ctx, msg)

	var ipList []string
	if ipType == common.IsPublicIp {		//	mcloud是公网ip
		ipList = append(ipList, msg.Data.WanIp)
		if err := dns.DnsAdd(msg.Data.DevHostname, ipList); err != nil {  // 添加映射
			ret.Result = common.RetServerErr			//	服务器内部错误
			common.RetJson(ctx, &ret)
			logger.Errorf("mcloudAdd error: %s\n", err.Error())
			return
		}
	}else {					//	mcloud是私网ip
		// 分配relay
		if ipList, err = allocRelay(msg.DevId, cfg.GetCommon().AllocRelayNum); err != nil {
			ret.Result = common.RetServerErr			//	服务器内部错误
			common.RetJson(ctx, &ret)
			logger.Errorf("allocRelay error: %s\n", err.Error())
			return
		}
		if len(ipList) == 0 {
			ret.Result = common.RetNoAvailableRelay			//	没有可用relay
			common.RetJson(ctx, &ret)
			return
		}
		ret.Data.RelayIp = ipList
	}

	// 添加mlcoud信息到mlcoud表
	if err = mcloudAdd(msg, ipType); err != nil {
		ret.Result = common.RetServerErr			//	服务器内部错误
		common.RetJson(ctx, &ret)
		logger.Errorf("mcloudAdd error: %s\n", err.Error())
		return
	}

	// 入网成功
	ret.Result = common.RetSuccess
	ret.Data.IpType = ipType
	common.RetJson(ctx, &ret)
}

// 关闭网络
func McloudNetStop(ctx  *fasthttp.RequestCtx, msg  *common.CliMsg)  {
	var ret common.RetInfo

	// 请求参数检查
	if msg.DevId == "" || msg.Data.DevHostname == "" {
		ret.Result = common.RetJsonFormatErr			// 请求参数格式错误或缺少参数
		common.RetJson(ctx, &ret)
		return
	}

	// 检查是否入网
	if model.IsExistMcloud(msg.DevId) == false {
		ret.Result = common.RetNoMcloudInfo			// 没有mcloud相关信息
		common.RetJson(ctx, &ret)
		return
	}

	// 删除mcloud信息
	if err := common.DelMcloud(msg.DevId, msg.Data.DevHostname); err != nil {
		ret.Result = common.RetServerErr			// 服务器内部出错
		common.RetJson(ctx, &ret)
		return
	}

	common.RelayMap.Del(msg.DevId);

	// 关网成功
	ret.Result = common.RetSuccess
	common.RetJson(ctx, &ret)
}

// 隧道状态上报
func McloudTunnelState(ctx  *fasthttp.RequestCtx, msg  *common.CliMsg)  {
	var ret common.RetInfo
	var err error
	var ipList []string

	// 请求参数检查
	if msg.DevId == "" || !( msg.Data.ConnRelayNum >=0 && msg.Data.ConnRelayNum <= cfg.GetCommon().AllocRelayNum  )||
		(len(msg.Data.FailedRelayIp) != 0 && common.IsValidIP(msg.Data.FailedRelayIp) == false ) ||
		(len(msg.Data.SuccessRelayIp) != 0 && common.IsValidIP(msg.Data.SuccessRelayIp) == false) ||
		msg.Data.DevHostname == "" || msg.Data.WanIp == "" {
		ret.Result = common.RetJsonFormatErr			// 请求参数格式错误或缺少参数
		common.RetJson(ctx, &ret)
		return
	}

	// 检查是否入网
	var mcloud *model.Mcloud
	if mcloud,err = model.GetMcloud(msg.DevId); err != nil {
		ret.Result = common.RetNoMcloudInfo			// 没有mcloud相关信息
		common.RetJson(ctx, &ret)
		return
	}

	// 判断是否有公网
	ret.Data.IpType = common.IsPublicNetwork(ctx, msg)
	if ret.Data.IpType == common.IsPublicIp {
		ret.Result = common.RetIsPublicIp			// 是公网ip
		common.RetJson(ctx, &ret)
		return
	}

	if msg.Data.WanIp != mcloud.McloudIP {
		ret.Result = common.RetIpChange			// ip变化, 由心跳处理
		common.RetJson(ctx, &ret)
		return
	}

	logger.Debugf("DevId: %v, ConnRelayNum = %v, SuccessRelayIp = %v, FailedRelayIp = %v\n",
		msg.DevId, msg.Data.ConnRelayNum, msg.Data.SuccessRelayIp, msg.Data.FailedRelayIp)

	var successList []string
	successList,_ = common.GetSuccessAndFailedList(msg.DevId)

	//  有连接成功的relay
	if len(msg.Data.SuccessRelayIp) != 0 || len(successList) != 0{
		var needDelList []string
		var needAddList []string

		// 有上报成功的就对比上报的与之前缓存的，没有上报成功的就直接删除缓存的
		if len(msg.Data.SuccessRelayIp) != 0 {
			needDelList, needAddList = getDiffList(successList, msg.Data.SuccessRelayIp)
		}else {
			needDelList = successList;
		}

		// 需要删除的成功记录
		if len(needDelList) != 0 {
			//  删除rcmap表映射
			if err:= model.DelRcmap(msg.DevId, needDelList); err!= nil {
				logger.Errorf("McloudTunnelState-rcMapDel error: %v", err)
				goto SERVER_ERROR
			}

			// 添加 失败记录到fieldMap
			common.RelayMap.Add(msg.DevId, needDelList, false)

			// 删除DNS映射
			if err := dns.DnsDelByDomainAndIp(msg.Data.DevHostname, needDelList); err != nil {
				logger.Errorf("McloudTunnelState-DnsDelByDomainAndIp error: %s", err.Error())
				goto SERVER_ERROR
			}
		}

		//	需要添加的成功记录
		if len(needAddList) != 0 {
			// 添加到RMmap表
			if err := model.AddRcmap(msg.DevId, needAddList); err != nil {
				logger.Errorf("McloudTunnelState-rmMapAdd err! err: %v\n", err)
				goto SERVER_ERROR
			}

			// 添加DNS映射
			if err := dns.DnsAdd(msg.Data.DevHostname, needAddList); err != nil {
				logger.Errorf("McloudTunnelState-DnsAdd error: %s", err.Error())
				goto SERVER_ERROR
			}
			// 添加 记录到fieldMap, 保证已经成功的不会在分配
			common.RelayMap.Add(msg.DevId, needAddList, true)
		}
	}

	//  有连接失败的relay
	if len(msg.Data.FailedRelayIp) != 0 {
		var newFalseRelayList []string
		newFalseRelayList = filterRelayList(msg.DevId, msg.Data.FailedRelayIp, false)

		// 更新relay表  conn_faild_num
		if err:= updateRelayFaileNum(msg.DevId ,newFalseRelayList); err!= nil {
			logger.Errorf("McloudTunnelState-updateRelayFaileNum error: %v", err)
			goto SERVER_ERROR
		}

		//  删除rcmap表映射
		if err:= model.DelRcmap(msg.DevId, newFalseRelayList); err!= nil {
			logger.Errorf("McloudTunnelState-rcMapDel error: %v", err)
			goto SERVER_ERROR
		}

		// 添加 失败记录到fieldMap
		common.RelayMap.Add(msg.DevId, newFalseRelayList, false)
		
		// 删除DNS映射
		if err := dns.DnsDelByDomainAndIp(msg.Data.DevHostname, newFalseRelayList); err != nil {
			logger.Errorf("McloudTunnelState-DnsDelByDomainAndIp error: %s", err.Error())
			goto SERVER_ERROR
		}
	}

	// common.DebugPrintExcludeIpList(msg.DevId)
	// 分配relay
	if ipList, err = allocRelay(msg.DevId, cfg.GetCommon().AllocRelayNum - msg.Data.ConnRelayNum); err != nil {
		logger.Errorf("allocRelay error: %s\n", err.Error())
		goto SERVER_ERROR
	}
	if len(ipList) == 0 {
		ret.Result = common.RetNoAvailableRelay			//	没有可用relay
		common.RetJson(ctx, &ret)
		return
	}
	ret.Data.RelayIp = ipList

	ret.Result = common.RetSuccess
	common.RetJson(ctx, &ret)
	return

SERVER_ERROR:
	ret.Result = common.RetServerErr		//	服务器内部错误
	common.RetJson(ctx, &ret)
}

// mcloud心跳
func McloudHeartBeat(ctx  *fasthttp.RequestCtx, msg  *common.CliMsg)  {
	var ret common.RetInfo
	var err error

	// 请求参数检查
	if msg.DevId == "" || msg.Data.DevHostname == "" || msg.Data.WanIp == "" ||
		common.IsValidIP(msg.Data.WanIp) == false {
		ret.Result = common.RetJsonFormatErr			// 请求参数格式错误或缺少参数
		common.RetJson(ctx, &ret)
		return
	}

	// 检查是否入网
	var mcloud *model.Mcloud
	if mcloud,err = model.GetMcloud(msg.DevId); err != nil {
		ret.Result = common.RetNoMcloudInfo			// 没有mcloud相关信息
		common.RetJson(ctx, &ret)
		return
	}

	// 判断是否有公网
	ipType := common.IsPublicNetwork(ctx, msg)
	ret.Data.IpType = ipType

	// ip变化
	if msg.Data.WanIp != mcloud.McloudIP {
		if mcloud.IsPublicIP == common.IsPublicIp {
			if ipType == common.IsPublicIp {		// 公网 - 公网
				// 更新DNS映射
				if err := dns.DnsDelByDomain(msg.Data.DevHostname); err != nil {
					logger.Errorf("DnsDelByDomain error! err: %v\n",err)
					goto SERVER_ERROR
				}
				if err := dns.DnsAdd(msg.Data.DevHostname, []string{msg.Data.WanIp}); err != nil {
					logger.Errorf("DnsDelByDomain error! err: %v\n",err)
					goto SERVER_ERROR
				}
			}else {							// 公网 - 私网
				// 删除mcloud信息
				if err := common.DelMcloud(msg.DevId, msg.Data.DevHostname); err != nil {
					logger.Errorf("delMcloud error! err: %v\n",err)
					goto SERVER_ERROR
				}
				ret.Result = common.RetNoMcloudInfo			// 没有mcloud信息，请重新入网
				common.RetJson(ctx, &ret)
				return
			}
		}else {
			if ipType == common.IsPublicIp {		// 私网 - 公网
				// 更新DNS映射
				if err := dns.DnsDelByDomain(msg.Data.DevHostname); err != nil {
					logger.Errorf("DnsDelByDomain error! err: %v\n",err)
					goto SERVER_ERROR
				}
				if err := dns.DnsAdd(msg.Data.DevHostname, []string{msg.Data.WanIp}); err != nil {
					logger.Errorf("DnsDelByDomain error! err: %v\n",err)
					goto SERVER_ERROR
				}

				// 删除 rcMap表中的映射
				if err := model.DelRcmapByMcloud(msg.DevId); err != nil {
					logger.Errorf("rcMapDel error! err: %v\n",err)
					goto SERVER_ERROR
				}
			}else {									// 私网 - 私网
				// 不做处理， 后面会更新ip
			}
		}
	}

	// 更新心跳时间、连接状态和IP
	if err := model.UpdateHeartStateIP(msg.DevId, msg.Data.WanIp, ipType); err != nil {
		logger.Errorf("updateHeartStateIP error! err: %v\n", err)
		goto SERVER_ERROR
	}

	ret.Result = common.RetSuccess
	common.RetJson(ctx, &ret)
	return

SERVER_ERROR:
	ret.Result = common.RetServerErr		//	服务器内部错误
	common.RetJson(ctx, &ret)
}

// 功能：	添加mcloud信息到mcloud表
// 输入		msg:	请求的JSON信息
//			ipType:	ip类型 0私网,1公网
// 返回： 	error:	错误信息
func mcloudAdd(msg  *common.CliMsg, ipType uint32)  (error) {
	mcloud := model.Mcloud{
		Devid:		msg.DevId,
		McloudIP:	msg.Data.WanIp,
		IsPublicIP:	ipType,
		HeartTime:	time.Now().Unix(),
		ConnState:	1,
		Ext:		"",
	}

	if err := mcloud.Add(); err != nil {
		logger.Errorf("mcloudAdd: insert mcloud! err: %v\n", err)
		return err
	}
	logger.Infof("mcloudAdd: %+v\n", mcloud)
	return nil
}

func allocRelay(devId string, num int) ([]string, error) {
	var ipList []string
	var err error

	excludeIpList := common.RelayMap.GetRelayList(devId)		// 获取例外relay IP列表
	logger.Infof("-- mcloudid %s: excludeIpList: %v", devId, excludeIpList)
	common.PrintExcludeIpList(devId)

	// 分配relay
	if ipList, err = model.GetAllocRelayList(num, excludeIpList, cfg.GetCommon().MaxMcloudNum); err != nil {
		logger.Errorf("allocRelay sql error! err: %v\n", err)
		return ipList,err
	}
	logger.Infof("mcloud_id: %s, allocRelay: %v\n", devId, ipList)

	// 更新relay 分配个数
	if err := model.UpdateRelayAllocNum(ipList); err != nil {
		logger.Error("update relay alloc num error! err: %s\n", err)
	}
	return ipList, nil
}

func updateRelayFaileNum(mcloudId string, relayIpList []string) error  {
	for _,ip := range relayIpList {
		if common.RelayMap.IpExist(mcloudId, ip, false) == true {
			continue
		}

		relay, err := model.GetRelayByIp(ip);
		if err != nil {
			if err != gorm.ErrRecordNotFound {
				logger.Errorf("updateRelayFaileNum update error! err: %v\n", err)
			}
			continue
		}

		relay.ConnFaildNum = relay.ConnFaildNum + 1
		logger.Infof("relay.ConnFaildNum = %v\n", relay.ConnFaildNum)
		if err := (*relay).Update(); err != nil {
			logger.Errorf("updateRelayFaileNum update error! err: %v\n", err)
			continue
		}
	}
	return nil
}

// 过滤 map已经存在的 mcloud连接的relay信息
func filterRelayList(mcloudId string, relayIpList []string, connState bool) []string {
	var newRelayList []string
	for _,ip := range relayIpList {
		if common.RelayMap.IpExist(mcloudId, ip, connState) == true {
			continue
		}
		newRelayList = append(newRelayList, ip)
	}
	logger.Infof("connState: %v, newRelayList: %v\n", connState, newRelayList)
	return newRelayList
}

// 比较两个slice，获取各自私有的数据
func getDiffList(list1 []string, list2 []string) ([]string, []string) {
	var list1HaveData []string
	var list2HaveData []string
	var shareData []string
	var flag int32

	for _, v1 := range list1 {
		for _, v2 := range list2 {
			if v1 == v2 {
				flag = 1
				break
			}
		}
		if flag == 0 {
			list1HaveData = append(list1HaveData, v1)
		}else {
			shareData = append(shareData, v1)
			flag = 0
		}
	}

	for _, v1 := range list2 {
		for _, v2 := range shareData {
			if v1 == v2 {
				flag = 1
				break
			}
		}
		if flag == 0 {
			list2HaveData = append(list2HaveData, v1)
		}else {
			flag = 0
		}
	}

	return list1HaveData,list2HaveData
}