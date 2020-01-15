package common

import (
	"encoding/json"
	"net"
	"time"

	logger "github.com/cihub/seelog"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	jsoniter "github.com/json-iterator/go"
	"github.com/valyala/fasthttp"
	"gitlab.com/TenbayMCloud/awesome-raserver/common/dns"
	cfg "gitlab.com/TenbayMCloud/awesome-raserver/config"
	"gitlab.com/TenbayMCloud/awesome-raserver/model"
)

// 设备位标记, 与数据库对应
const (
	RelayFlag  = 1
	McloudFlag = 2
)

// ip类型
const (
	IsPublicIp    = 1
	IsNotPublicIp = 0
)

// 功能：	向dev表添加新设备信息
// 输入		msg：	请求的JSON信息
//			devFlag: 设备类型，即请求角色
// 返回： 	error:	错误信息
func DevAdd(msg *CliMsg, devFlag uint32) error {
	dev, _ := model.GetDev(msg.DevId)
	if dev.Devid != "" { // 存在, 更新role
		if err := dev.UpdateRole(devFlag); err != nil {
			logger.Errorf("DevAdd: update role! err: %v\n", err)
			return err
		}
		return nil
	}

	// 不存在， 创建dev
	dev = &model.Device{
		Devid:         msg.DevId,
		Mac:           msg.Data.DevMac,
		HostName:      msg.Data.DevHostname,
		Vendor:        msg.Data.Vendor,
		SoftModel:     msg.Data.SoftModel,
		SoftVersion:   msg.Data.SoftVersion,
		HardwareModel: msg.Data.HardwareModel,
		JoinTime:      time.Now(),
		Role:          devFlag,
		Ext:           "",
	}

	if err := dev.Add(); err != nil {
		logger.Errorf("DevAdd error! err: %v\n", err)
		return err
	}
	logger.Infof("DevAdd: %+v\n", *dev)
	return nil
}

// 功能：	从relay表删除relay信息
// 输入		devId:	设备ID
// 返回： 	error:	错误信息
func DelRelay(devId string) error {
	// 删除Relay信息
	var relay model.Relay
	relay.Devid = devId
	if err := relay.Del(); err != nil {
		logger.Errorf("del relay sql error! relayId:%v , err: %v\n", devId, err)
		return err
	}
	logger.Infof("del relay relayId: %v\n", devId)

	// 删除dev表信息
	if model.IsExistMcloud(devId) == false {
		if err := model.DelDev(devId); err != nil {
			logger.Errorf("del dev error! dev_id = %s\n", devId)
			return err
		}
		logger.Infof("del dev dev_id = %v\n", devId)
	}
	return nil
}

func DelRelayRole(devId string) error {
	dev, _ := model.GetDev(devId)
	if dev.Devid != "" { // 存在, 更新role
		if err := dev.DelRelayRole(); err != nil {
			logger.Errorf("DevAdd: DelRelayRole role! err: %v\n", err)
			return err
		}
		return nil
	}
	return nil
}

// 删除和relay建立连接的映射
func DeleteDnsByRelay(relay *model.Relay) (error) {
	dev, _ := model.GetDev(relay.Devid)
	if dev.Devid == "" && dev.HostName == "" { 		// 不存在
		return nil
	}

	// 删除时忽略自己的映射
	if err := dns.DnsDelByIpIgnoreDomain(relay.WanIP, dev.HostName); err != nil {
		logger.Errorf("relay netstop, DnsDelByIp err: %v\n",err)
	}
	return nil
}

func DelMcloud(devId string, domain string) error {
	// 删除mcloud信息
	var mcloud model.Mcloud
	mcloud.Devid = devId
	if err := mcloud.Del(); err != nil {
		logger.Errorf("del Mcloud sql error! mcloudId:%s , err: %v\n", devId, err)
		return err
	}
	logger.Infof("del mcloud mcloud_id = %v\n", devId)

	// 删除存储连接失败relay的map信息
	RelayMap.Del(devId)

	// 删除Rcmap表信息
	var rcMap model.Rcmap
	rcMap.McloudID = devId
	if err := model.DelRcmapByMcloud(devId); err != nil {
		logger.Errorf("del Rcmap sql error! mcloudId:%s , err: %v\n", devId, err)
		return err
	}

	// 删除DNS 映射
	if domain == "" {
		dev, err := model.GetDev(devId)
		if err != nil {
			logger.Errorf("select dev error!mcloudId:%s, err: %v\n", devId, err)
			return err
		}
		domain = dev.HostName
	}
	if err := dns.DnsDelByDomain(domain); err != nil {
		logger.Errorf("DnsDelByDomain error! mcloudId:%s ,domain: %s,  err: %v\n", devId, domain, err)
		return err
	}

	// 删除dev表信息
	if model.IsExistRelay(devId) == false {
		if err := model.DelDev(devId); err != nil {
			logger.Errorf("del dev error! dev_id = %s\n", devId)
			return err
		}
	}
	logger.Infof("del dev dev_id = %v\n", devId)
	return nil
}

// 获取mac和hostname的映射列表
func GetDomainMapList(macList []string) ([]Mcloud, error) {
	var mclouds []Mcloud

	macHostnameMapPtr, err := model.GetDomainByMacList(macList)
	if err != nil {
		logger.Error("GetDomainByMacList error! err: %v\n", err)
		return mclouds, err
	}

	var mcloud Mcloud
	for k, v := range *macHostnameMapPtr {
		mcloud.Mac = k
		mcloud.Hostname = v
		mclouds = append(mclouds, mcloud)
	}
	return mclouds, nil
}

// 功能：	判断IP是否合法
// 输入		ip：	需要判断的ip
// 返回： 	bool:	true 合法  false 不合法
func IsValidIP(ip interface{}) bool {
	switch ip.(type) {
	case string:
		oneIp, _ := ip.(string)
		if net.ParseIP(oneIp) == nil {
			return false
		}
	case []string:
		ips, _ := ip.([]string)
		for _, v := range ips {
			if net.ParseIP(v) == nil {
				return false
			}
		}
	default:
		return false
	}
	return true
}

// 功能：	判断设备是否有公网IP
// 输入		ctx：	请求结构
//			msg:	请求的JSON信息
// 返回： 	uint32:	IP类型  0 公网ip  1 私网ip
//			error:	错误信息
func IsPublicNetwork(ctx *fasthttp.RequestCtx, msg *CliMsg) uint32 {
	remoteAddr := ctx.Request.Header.Peek("X-Forwarded-For")
	if remoteAddr == nil {
		logger.Debugf("No proxy remote ip!(Request head have no \"X-Forwarded-For\"!)")
		return IsNotPublicIp
	}
	if string(remoteAddr) == msg.Data.WanIp {
		return IsPublicIp
	}
	return IsNotPublicIp
}

// 初始化dns配置
func InitDnsConfig() error {
	if err := dns.CpFile(cfg.GetDNS().NasDnsBkFileName, cfg.GetDNS().NasDnsFileName); err != nil {
		logger.Errorf("mv dns file failed!err: %v\n", err)
		return err
	}

	rcMap, err := model.GetAllRcmaps()
	if err != nil {
		logger.Errorf("initDnsConfig: error! err: %v\n", err)
		return err
	}
	// 初始化内网mcloud 的dns映射
	var relay *model.Relay
	var dev *model.Device
	// relay ip : mcloud mac
	for _, v := range *rcMap {
		logger.Infof("relay_id: %v, mcloud_id: %v", v.RelayID, v.McloudID)
		if relay, _ = model.GetRelay(v.RelayID); relay.Devid == "" {
			continue
		}

		if dev, _ = model.GetDev(v.McloudID); dev.Devid == "" {
			continue
		}

		logger.Infof("realy.PublicIP: %v, dev.HostName: %v", relay.PublicIP, dev.HostName)
		if err := dns.DnsAdd(dev.HostName, []string{relay.PublicIP}); err != nil {
			logger.Errorf("initDnsConfig: add rule failed! err: %v\n", err)
			return err
		}
	}

	// 初始化外网mcloud 的dns映射
	mclouds, err := model.GetAllPublicMclouds()
	if err != nil {
		logger.Errorf("initDnsConfig: GetAllPublicMclouds error! err: %v\n", err)
		return err
	}

	for _, v := range *mclouds {
		logger.Infof("mcloud_id: %v", v.Devid)

		if dev, _ = model.GetDev(v.Devid); dev.Devid == "" {
			continue
		}

		logger.Infof("mcloud.McloudIP: %v, mcloud.HostName: %v", v.McloudIP, dev.HostName)
		if err := dns.DnsAdd(dev.HostName, []string{v.McloudIP}); err != nil {
			logger.Errorf("initDnsConfig: add rule failed! err: %v\n", err)
			return err
		}
	}

	return nil
}

// JSONDecode 把json字符串([]byte)转换为Struct
func JSONDecode(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

// JSONEncode 把Struct转换为json字符串([]byte)
func JSONEncode(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

// ArgsGet 获取GET和POST参数
// 参数：
// 	gwant: 希望得到的URL参数列表
// 	pwant: 希望得到的请求体参数列表
// 返回: 一个string和[]byte的map
func ArgsGet(ctx *fasthttp.RequestCtx,
	gwant []string, pwant []string) map[string][]byte {

	var Args = make(map[string][]byte)

	var getArgs = ctx.QueryArgs()
	for _, gw := range gwant {
		//Args[gw] = string(getArgs.Peek(gw))
		var t = getArgs.Peek(gw)
		if len(t) != 0 {
			Args[gw] = t
		}
	}

	var postArgs = ctx.PostArgs()
	for _, pw := range pwant {
		//Args[pw] = string(postArgs.Peek(pw))
		var t = postArgs.Peek(pw)
		if len(t) != 0 {
			Args[pw] = t
		}
	}

	return Args
}

// ResponseInfo 响应的固定结构
type ResponseInfo struct {
	Status     int         `json:"status"`
	StatusDesc string      `json:"status_desc"`
	Data       interface{} `json:"data"`
}

// ResponseJSON 返回Json信息
// 注意：data如果是json，里面的Key应该要首字母大写！
func ResponseJSON(ctx *fasthttp.RequestCtx,
	status int, statusDesc string, data interface{}) error {

	respInfo := ResponseInfo{
		Status:     status,
		StatusDesc: statusDesc,
		Data:       data,
	}

	jsonIterator := jsoniter.ConfigCompatibleWithStandardLibrary
	b, err := jsonIterator.Marshal(respInfo)
	if err != nil {
		logger.Errorf("json Marshal err: %s", err.Error())
		ctx.Error("InternalServerError", fasthttp.StatusInternalServerError)
		return err
	}

	ctx.Write(b)
	return nil
}

// ResponseOkJSON -- OK
func ResponseOkJSON(ctx *fasthttp.RequestCtx, data interface{}) error {
	return ResponseJSON(ctx, 0, "OK", data)
}

// ResponseErrJSON -- ERROR
func ResponseErrJSON(ctx *fasthttp.RequestCtx, data interface{}) error {
	return ResponseJSON(ctx, 1, "ERROR", data)
}

const (
	TRUE  = 1
	FALSE = 0
)
