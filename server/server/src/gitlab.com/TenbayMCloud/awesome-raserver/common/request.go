package common

import (
	logger "github.com/cihub/seelog"
	jsoniter "github.com/json-iterator/go"
	"github.com/valyala/fasthttp"
)

// 请求操作类型
const (
	OpNetJoin       = "net_join"       // 入网
	OpNetStop       = "net_stop"       // 关闭网络
	OpRequestDomain = "request_domain" // 请求mcloud域名
	OpAddMcloud     = "add_mcloud"     // 添加mcloud
	OpTunnelState   = "tunnel_state"   // 隧道异常
	OpHeartBeat     = "heart_beat"     // 心跳
)

// 客户端请求消息结构
type CliMsg struct {
	DevId string `json:"dev_id"` //	设备ID
	Data  myData `json:"data"`   //	请求数据
}

// 消息中的Data数据结构
type myData struct {
	DevMac         string      `json:"mac"`
	DevHostname    string      `json:"hostname"`
	Vendor         string      `json:"vendor"`
	SoftModel      string      `json:"soft_model"`
	SoftVersion    string      `json:"soft_version"`
	HardwareModel  string      `json:"hardware_model"`
	WanIp          string      `json:"wan_ip"`
	McloudMacList  []string    `json:"mclouds"`
	RelayCpu       RelayCpu    `json:"cpu"`
	RelayMemory    RelayMemory `json:"memory"`
	McloudNum      uint32      `json:"mcloud_num"`
	ConnRelayNum   int         `json:"conn_relay_num"`
	SuccessRelayIp []string    `json:"success_relay_ip"`
	FailedRelayIp  []string    `json:"failed_relay_ip"`
}

// relay的cpu信息
type RelayCpu struct {
	State string `json:"state"`
	Used  string `json:"used"`
}

// relay的内存信息
type RelayMemory struct {
	Used string `json:"used"`
	Free string `json:"free"`
}

type RelayExt struct {
	Cpu RelayCpu
	Mem RelayMemory
}

type DevInfo struct {
	DevVendor   string `json:"dev_vendor"`
	DevProduct  string `json:"dev_product"`
	DevType     string `json:"dev_type"`
	DevId       string `json:"dev_id"`
	DevPlatform string `json:"dev_platform"`
	dev_softver string `json:"dev_softver"`
	dev_desc    string `json:"dev_desc"`
}

func JsonUnmarshal(ctx *fasthttp.RequestCtx, msg *CliMsg) error {
	json_iterator := jsoniter.ConfigCompatibleWithStandardLibrary
	err := json_iterator.Unmarshal(ctx.PostBody(), msg) // 解析Json并赋值到结构体 msg 中
	if err != nil {
		logger.Errorf("json.Unmarshal error: %s\n", err.Error())
		return err
	}
	return nil
}
