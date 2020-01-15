package model

import (
	logger "github.com/cihub/seelog"
	"github.com/jinzhu/gorm"
	database "gitlab.com/TenbayMCloud/awesome-raserver/common/dbMgr"
)

type Rcmap	struct {
	RelayID			string				`gorm:"column:relay_id;not null;size:128;primary_key"`
	McloudID		string				`gorm:"column:mcloud_id;not null;size:128;primary_key"`
}

func (rcMap *Rcmap) Add() error {
	db := database.GetMysql()
	if err := db.Create(rcMap).Error; err != nil {
		return err
	}
	return nil
}

func (rcMap *Rcmap) del() error {
	db := database.GetMysql()
	if err := db.Delete(rcMap).Error; err != nil {
		return err
	}
	return nil
}

func GetAllRcmaps() (*[]Rcmap ,error) {
	var rcMaps []Rcmap

	db := database.GetMysql()
	if err := db.Find(&rcMaps).Error; err != nil {
		return &rcMaps,nil
	}
	return &rcMaps,nil
}

func DelRcmapByRelay(devId string) error  {
	db := database.GetMysql()
	if err := db.Where("relay_id=?",devId).Delete(&Rcmap{}).Error; err != nil {
		return err
	}
	return nil
}

func DelRcmapByMcloud(devId string) error  {
	db := database.GetMysql()
	if err := db.Where("mcloud_id=?",devId).Delete(&Rcmap{}).Error; err != nil {
		return err
	}
	return nil
}

func AddRcmap(mcloudId string, relayIpList []string) error  {
	rcMap := Rcmap{
		McloudID:	mcloudId,
	}

	db := database.GetMysql()
	for _,ip := range relayIpList {
		var relay Relay
		if err := db.Where("public_ip = ?", ip).First(&relay).Error; err != nil {
			logger.Errorf("add rcmap error! mcloudid = %v, relay_ip = %v, select relay err: %v\n",
				mcloudId, ip, err)
			continue
		}

		rcMap.RelayID = relay.Devid
		if err := db.Create(&rcMap).Error; err != nil {
			//logger.Errorf("rcMapAdd error! err: %v\n", err)
		}
	}
	return nil
}

func DelRcmap(mcloudId string, relayIpList []string) error  {
	rcMap := Rcmap{
		McloudID:	mcloudId,
	}

	db := database.GetMysql()
	for _,ip := range relayIpList {
		var relay Relay
		if err := db.Where("public_ip = ?", ip).First(&relay).Error; err != nil {
			if err != gorm.ErrRecordNotFound {
				logger.Errorf("del rcmap error! mcloudid = %v, relay_ip = %v, select relay err: %v\n",
					mcloudId, ip, err)
			}
			continue
		}

		rcMap.RelayID = relay.Devid
		if err := db.Delete(&rcMap).Error; err != nil {
			logger.Errorf("rmMapDel error! err: %v\n", err)
			return err
		}
	}
	return nil
}


