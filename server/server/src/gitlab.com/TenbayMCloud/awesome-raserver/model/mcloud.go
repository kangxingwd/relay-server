package model

import (
	logger "github.com/cihub/seelog"
	database "gitlab.com/TenbayMCloud/awesome-raserver/common/dbMgr"
)

type Mcloud struct {
	Devid			string			`gorm:"column:dev_id;unique_index;size:128;not null;primary_key"`
	McloudIP		string			`gorm:"column:mcloud_ip;default:null"`
	IsPublicIP		uint32			`gorm:"column:is_public_ip;default:0"`
	HeartTime		int64			`gorm:"column:heart_time;default:0"`
	ConnState		uint32			`gorm:"column:conn_state;default:0"`
	Ext				string			`gorm:"column:ext;default:null"`
}

func (mcloud *Mcloud) Add() error {
	db := database.GetMysql()
	if err := db.Create(mcloud).Error; err != nil {
		return err
	}
	return nil
}

func (mcloud *Mcloud) Del() error {
	db := database.GetMysql()
	if err := db.Delete(mcloud).Error; err != nil {
		return err
	}
	return nil
}

func (mcloud *Mcloud) Update() error {
	db := database.GetMysql()
	if err := db.Save(mcloud).Error; err != nil {
		return err
	}
	return nil
}

func IsExistMcloud(devId string) bool {
	var mcloud Mcloud
	db := database.GetMysql()
	if err := db.Where(&Mcloud{Devid:devId}).First(&mcloud).Error; err != nil {
		return false
	}
	return true
}

func GetMcloud(devId string) (*Mcloud, error) {
	var mcloud Mcloud
	db := database.GetMysql()
	if err := db.Where(&Mcloud{Devid:devId}).First(&mcloud).Error; err != nil {
		return &mcloud, err
	}
	return &mcloud, nil
}

func GetAllMclouds() (*[]Mcloud ,error) {
	var mclouds []Mcloud

	db := database.GetMysql()
	if err := db.Find(&mclouds).Error; err != nil {
		return &mclouds,nil
	}
	return &mclouds,nil
}

func UpdateMcloudState(devId string, newState uint32) error {
	db := database.GetMysql()
	err := db.Model(&Mcloud{}).Where("dev_id = ?", devId).
		Update("conn_state", newState).Error
	if err != nil {
		logger.Errorf("updateReleyState sql failed! dev_id = %s\n", devId)
		return err
	}
	logger.Infof("update relay set conn_state = %d where dev_id = %s\n", newState, devId)
	return nil
}

func GetAllPublicMclouds() (*[]Mcloud ,error) {
	var mclouds []Mcloud

	db := database.GetMysql()
	if err := db.Where("is_public_ip = ?", 1).Find(&mclouds).Error; err != nil {
		return &mclouds,nil
	}
	return &mclouds,nil
}