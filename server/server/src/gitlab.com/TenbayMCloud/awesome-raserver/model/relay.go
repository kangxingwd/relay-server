package model

import (
	logger "github.com/cihub/seelog"
	"github.com/jinzhu/gorm"
	database "gitlab.com/TenbayMCloud/awesome-raserver/common/dbMgr"
	"time"
)

type Relay struct {
	Devid			string			`gorm:"column:dev_id;unique_index;size:128;not null;primary_key"`
	PublicIP		string			`gorm:"column:public_ip;size:32;default null"`
	WanIP			string			`gorm:"column:wan_ip;size:32;default null"`
	HeartTime		int64			`gorm:"column:heart_time;default:0"`
	ConnState		uint32			`gorm:"column:conn_state;default:0"`
	McloudNum		uint32			`gorm:"column:mcloud_num;default:0"`
	RelayAllocNum	uint32			`gorm:"column:relay_alloc_num;default:0"`
	ConnFaildNum	uint32			`gorm:"column:conn_faild_num;default:0"`
	FailedCountTime	int64			`gorm:"column:failed_count_time;default:0"`
	ForbidState		uint32			`gorm:"column:forbid_state;default:0"`
	OverLoad        uint32          `gorm:"column:overload;default:0"`
	Ext				string			`gorm:"column:ext;default:null" json:"Ext"`
}

func (relay *Relay) Add() error {
	db := database.GetMysql()
	if err := db.Create(relay).Error; err != nil {
		return err
	}
	return nil
}

func (relay *Relay) Del() error {
	db := database.GetMysql()
	if err := db.Delete(relay).Error; err != nil {
		return err
	}
	return nil
}

func (relay *Relay) Update() error {
	db := database.GetMysql()
	if err := db.Save(relay).Error; err != nil {
		return err
	}
	return nil
}

func IsExistRelay(devId string) bool {
	var relay Relay
	db := database.GetMysql()
	if err := db.Where(&Relay{Devid:devId}).First(&relay).Error; err != nil {
		return false
	}
	return true
}

func GetRelay(devId string) (*Relay, error) {
	var relay Relay
	db := database.GetMysql()
	if err := db.Where(&Relay{Devid:devId}).First(&relay).Error; err != nil {
		return &relay, err
	}
	return &relay, nil
}

func GetRelayByIp(ip string) (*Relay, error) {
	var relay Relay
	db := database.GetMysql()
	if err := db.Where("public_ip = ?", ip).First(&relay).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			logger.Errorf("updateRelayFaileNum update error! err: %v\n", err)
		}
		return &relay, err
	}
	return &relay, nil
}

func GetAllRelays() (*[]Relay ,error) {
	var relays []Relay

	db := database.GetMysql()
	if err := db.Find(&relays).Error; err != nil {
		return &relays,nil
	}
	return &relays,nil
}

func GetAllocRelayList(num int, excludeIpList []string, maxMcloudNum int) ([]string, error) {
	var err error
	var ipList []string
	db := database.GetMysql()
	if len(excludeIpList) == 0{
		err = db.Table("relays").
			Where("conn_state = 1 AND mcloud_num < ? AND forbid_state = 0 AND overload = 0",
			maxMcloudNum).
			Limit(num).Pluck("public_ip",&ipList).Error
	}else {
		err = db.Table("relays").
			Where("conn_state = 1 AND mcloud_num < ? AND forbid_state = 0 AND overload = 0",
			maxMcloudNum).
			Not("public_ip",excludeIpList).
			Limit(num).Pluck("public_ip",&ipList).Error
	}
	if err != nil {
		return ipList, err
	}
	return ipList, nil
}

func UpdateRelayAllocNum(ipList []string) error {
	db := database.GetMysql()
	for _,ip := range ipList {
		var relay Relay
		if err := db.Where("public_ip = ?", ip).First(&relay).Error; err != nil {
			logger.Errorf("update relay alloc num select error! ip: %v, err: %s\n", ip, err)
			continue
		}
		relay.RelayAllocNum = relay.RelayAllocNum + 1
		if err := db.Save(&relay).Error; err != nil {
			logger.Errorf("update relay alloc num save error!ip: %v, err: %s\n", ip, err)
		}
	}
	return nil
}

func UpdateHeartStateIP(devId string, ip string, ipType uint32) error  {
	db := database.GetMysql()

	var mcloud Mcloud
	
	//if err := db.First(&mcloud, devId).Error; err != nil {
	if err := db.Where("dev_id = ?", devId).First(&mcloud).Error; err != nil {
		logger.Errorf("updateHeartStateIP sql error!devId: %v, err: %v\n", devId, err)
		return err
	}
	mcloud.McloudIP = ip
	mcloud.IsPublicIP = ipType
	mcloud.HeartTime = time.Now().Unix()
	mcloud.ConnState = 1
	if err := db.Save(&mcloud).Error; err!= nil {
		logger.Errorf("updateHeartStateIP sql error!devId: %v, err: %v\n", devId, err)
		return err
	}

	//mcloud.Devid = devId
	//err := db.Model(&mcloud).Updates(map[string]interface{}{
	//	"mcloud_ip": 		ip,
	//	"is_public_ip": 	ipType,
	//	"heart_time": 		time.Now().Unix(),
	//	"conn_state": 		1}).
	//	Error
	//if err != nil {
	//	logger.Errorf("updateHeartStateIP sql error!devId: %v, err: %v\n", devId, err)
	//	return err
	//}
	return nil
}

func UpdateReleyState(devId string, newState uint32) error {
	db := database.GetMysql()
	err := db.Model(&Relay{}).Where("dev_id = ?", devId).
		Update("conn_state", newState).Error
	if err != nil {
		logger.Errorf("updateReleyState sql failed! dev_id = %s\n", devId)
		return err
	}
	logger.Infof("update relay set conn_state = %d where dev_id = %s\n", newState, devId)
	return nil
}

func UpdateForbid(devId string, forbidState uint32)  {
	db := database.GetMysql()
	err := db.Model(&Relay{}).Where("dev_id = ?", devId).
		Updates(map[string]interface{}{"relay_alloc_num": 0, "conn_faild_num": 0,
			"forbid_state": forbidState, "failed_count_time":time.Now().Unix()}).Error
	if err != nil {
		logger.Errorf("setForbid sql failed! dev_id = %s\n", devId)
		return
	}
	logger.Infof("setForbid dev_id = %s, forbidState = %d\n", devId, forbidState)
	return
}

