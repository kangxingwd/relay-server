package common

import (
	logger "github.com/cihub/seelog"
	"github.com/json-iterator/go"
	"github.com/valyala/fasthttp"
)

// 返回码
const (
	RetSuccess = 0 // 返回成功

	// 提交格式错误
	RetUrlErr           = 1 // url参数错误，字段拼写错误或者多余字段，或者少字段
	RetJsonFormatErr    = 2 // Json格式错误，字段错误，提交的字段拼写错误或者多余或者少字段
	RetDevTypeErr       = 3 // 设备类型错误
	RetRequestMethodErr = 4 // 请求格式错误 需要post
	RetRequestOptionErr = 5 // 请求操作类型错误

	// 业务错误类型
	RetPasswdAuthFailed = 101 // 密码认证失败
	RetTokenAuthFailed  = 102 // token认证失败，先进行密码认证
	RetTokenExpire      = 103 // token参数失效，参数过期
	RetNoPublicIp       = 104 // 没有公网IP
	RetIsPublicIp       = 105 // 是公网IP
	RetRepetJoin        = 106 // 重复提交入网
	RetRequestMacIsNull = 107 // 请求的域名为空
	RetNoDomainInfo     = 108 // 没有对应的域名信息
	RetNoRelayInfo      = 109 // 没有relay信息，请重新入网
	RetNoMcloudInfo     = 110 // 没有mcloud信息，请重新入网
	RetAddDnsFailed     = 111 // 添加dns映射失败
	RetNoAvailableRelay = 112 // 没有可用relay
	RetIpChange         = 113 // ip发生变化
	RetCpuMemFormatErr  = 114 // cpu或mem格式错误

	// 服务端错误
	RetServerErr = 201 // 服务器内部出错
)

// 错误信息
var ErrInfoMap = map[int]string{
	RetSuccess:          "Success!",
	RetUrlErr:           "Url args error !",
	RetJsonFormatErr:    "Json format error or miss parameter!",
	RetDevTypeErr:       "Dev type error!",
	RetRequestMethodErr: "Request method error!",
	RetRequestOptionErr: "Request Option error!",
	RetPasswdAuthFailed: "Passwd auth failed!",
	RetTokenAuthFailed:  "Token auth failed, verify the password first!",
	RetTokenExpire:      "Token expire!",
	RetNoPublicIp:       "No public IP!",
	RetIsPublicIp:       "Is public IP!",
	RetRepetJoin:        "Repet commit net join!",
	RetServerErr:        "Server interior error!",
	RetRequestMacIsNull: "Request domain is null!",
	RetNoDomainInfo:     "No domain info about of these mcloud mac!",
	RetNoRelayInfo:      "No relay info, please request netjoin!",
	RetNoMcloudInfo:     "No mcloud info, please request netjoin!",
	RetAddDnsFailed:     "Add dns failed!",
	RetNoAvailableRelay: "No available Relay!",
	RetIpChange:         "Ip change!",
	RetCpuMemFormatErr:  "Cpu or mem format error!",
}

// 返回信息结构
type RetInfo struct {
	Result  int     `json:"result"`   // 返回结果
	ErrInfo string  `json:"err_info"` // 错误信息
	Data    retData `json:"data"`     // 返回数据
}

// 返回的数据
type retData struct {
	Mclouds    []Mcloud `json:"mclouds"`
	IpType     uint32   `json:"ip_type"` // ip类型， 0： 私网， 1：公网
	RelayIp    []string `json:"relay_ip"`
	Token      string   `json:"token"`
	ExpireTime string   `json:"expire_time"`
}

// mcloud的mac和域名映射
type Mcloud struct {
	Mac      string `json:"mac"`
	Hostname string `json:"hostname"`
}

// 返回Json格式的回应
func RetJson(ctx *fasthttp.RequestCtx, retInfo *RetInfo) {
	retInfo.ErrInfo = ErrInfoMap[retInfo.Result]

	json_iterator := jsoniter.ConfigCompatibleWithStandardLibrary
	b, err := json_iterator.Marshal(*retInfo)
	if err != nil {
		logger.Errorf("json Marshal err: %s", err.Error())
		ctx.Error("InternalServerError", fasthttp.StatusInternalServerError)
		return
	}
	ctx.Write(b)
}