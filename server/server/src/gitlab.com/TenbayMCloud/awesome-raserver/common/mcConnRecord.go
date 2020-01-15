package common

import (
	logger "github.com/cihub/seelog"
	cfg "gitlab.com/TenbayMCloud/awesome-raserver/config"
	"sync"
	"time"
)

// 记录mcloud连接的relay,成功的和失败的relay 都记录
type relay struct {
	ip        string
	time      int64
	connState bool // 连接状态
}

type relayMap struct {
	Mutex sync.RWMutex
	Map   map[string][]relay
}

// 存储mcloud连接的relay信息   key: mcloud_id  value: relay 列表
var RelayMap = NewRelayMap()

// 自定义构造函数
func NewRelayMap() *relayMap {
	connRelayMap := make(map[string][]relay)
	return &relayMap{Map: connRelayMap}
}


func (relayMap *relayMap) IpExist(devId string, ip string, connState bool) bool {
	relayMap.Mutex.RLock()
	defer relayMap.Mutex.RUnlock()

	if _, ok := relayMap.Map[devId]; ok {
		ipList := relayMap.Map[devId]
		for _, v := range ipList {
			if v.connState == connState && v.ip == ip {
				return true
			}
		}
	}
	return false
}

func ipExist(devId string, ip string) (ret bool, index int) {
	if _, ok := RelayMap.Map[devId]; ok {
		ipList := RelayMap.Map[devId]
		for index, v := range ipList {
			if v.ip == ip {
				return true, index
			}
		}
	}
	return false, 0
}

func (relayMap *relayMap) Add(devId string, ip []string, connState bool) {
	relayMap.Mutex.Lock()
	defer relayMap.Mutex.Unlock()

	var ipList []relay
	for _, v := range ip {
		if ret, index := ipExist(devId, v); ret == false { // 不存在 添加
			ipList = append(ipList, relay{v, time.Now().Unix(), connState})
		} else { // 若存在, 更新状态
			relayMap.Map[devId][index].connState = connState
		}
	}

	if _, ok := relayMap.Map[devId]; ok {
		relayMap.Map[devId] = append(relayMap.Map[devId], ipList...)
	} else {
		relayMap.Map[devId] = ipList
	}
}

func (relayMap *relayMap) Del(devId string) {
	relayMap.Mutex.Lock()
	defer relayMap.Mutex.Unlock()
	delete(relayMap.Map, devId)
}

func (relayMap *relayMap) GetRelayList(devId string) []string {
	relayMap.Mutex.Lock()
	defer relayMap.Mutex.Unlock()

	var ip []string
	var relayListTmp []relay
	if _, ok := relayMap.Map[devId]; ok {
		ipList := relayMap.Map[devId]
		if len(ipList) == 0 {
			return ip
		}
		for _, v := range ipList {
			// 删除超过5分钟的失败记录, 保存成功记录和5分钟内的失败记录
			if (v.connState == false && time.Now().Unix()-v.time < cfg.GetCommon().DeleteRalayRecordTime) ||
				v.connState == true {
				relayListTmp = append(relayListTmp, v)
				ip = append(ip, v.ip)
			}
		}
		relayMap.Map[devId] = relayListTmp
	}
	return ip
}


func PrintExcludeIpList(devid string) {
	var successIp []string
	var failedIp []string

	for _, v := range RelayMap.Map[devid] {
		if v.connState == true {
			successIp = append(successIp, v.ip)
		} else {
			failedIp = append(failedIp, v.ip)
		}
	}
	logger.Infof("mcloud_id: %s, success_relay_ip: %v", devid, successIp)
	logger.Infof("mcloud_id: %s, failed_relay_ip: %v\n", devid, failedIp)
}

func DebugPrintExcludeIpList(devid string) {
	var successIp []string
	var failedIp []string

	for _, v := range RelayMap.Map[devid] {
		if v.connState == true {
			successIp = append(successIp, v.ip)
		} else {
			failedIp = append(failedIp, v.ip)
		}
	}
	logger.Debugf("mcloud_id: %s, success_relay_ip: %v", devid, successIp)
	logger.Debugf("mcloud_id: %s, failed_relay_ip: %v\n", devid, failedIp)
}

func GetSuccessAndFailedList(devid string) ([]string, []string) {
	var successIp []string
	var failedIp []string

	for _, v := range RelayMap.Map[devid] {
		if v.connState == true {
			successIp = append(successIp, v.ip)
		} else {
			failedIp = append(failedIp, v.ip)
		}
	}
	return successIp,failedIp
}
