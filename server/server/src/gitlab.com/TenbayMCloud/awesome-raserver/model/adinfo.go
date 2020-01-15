package model

import (
	"fmt"
	"time"

	database "gitlab.com/TenbayMCloud/awesome-raserver/common/dbMgr"
)

// ORICO-NAS-WRC10 => vendor-product-type
type AdInfo struct {
	ID         int64  `json:"id" gorm:"column:id;primary_key;AUTO_INCREMENT"`
	DevVendor  string `json:"dev_vendor" gorm:"column:dev_vendor;size:64;default:null"`
	DevProduct string `json:"dev_product" gorm:"column:dev_product;size:64;default:null"`
	DevType    string `json:"dev_type" gorm:"column:dev_type;size:64;not null"` // 这个不为0
	DevSoftVer string `json:"dev_softver" gorm:"column:dev_softver;default:null"`
	//DevPlatform		string			`json:"dev_platform";gorm:"column:dev_platform;default:null"`
	AdURL          string `json:"ad_url" gorm:"column:ad_url;default:null"`
	Desc           string `json:"desc" gorm:"column:desc;default:null"`
	LastUpdateTime int64  `json:"last_update_time" gorm:"column:last_update_time;default:0"`
	CreateTime     int64  `json:"create_time" gorm:"column:create_time;default:0"`
	LastAccessTime int64  `json:"last_access_time" gorm:"column:last_access_time;default:0"`
	Ext            string `json:"ext" gorm:"column:ext;default:null"`
}

func (info *AdInfo) IsExist() bool {
	var adInfo AdInfo
	db := database.GetMysql()
	if err := db.Where(info).First(&adInfo).Error; err != nil {
		return false
	}
	return true
}

func (info *AdInfo) AddNoDup() error {
	if info.IsExist() {
		return fmt.Errorf("This AD info already exist")
	}

	info.CreateTime = time.Now().Unix()
	db := database.GetMysql()
	if err := db.Create(info).Error; err != nil {
		return err
	}
	return nil
}

func (info *AdInfo) Add() error {
	info.CreateTime = time.Now().Unix()
	db := database.GetMysql()
	if err := db.Create(info).Error; err != nil {
		return err
	}
	return nil
}

func (info *AdInfo) Del() error {
	db := database.GetMysql()
	if err := db.Delete(info).Error; err != nil {
		return err
	}
	return nil
}

func (info *AdInfo) Update() error {
	now := time.Now().Unix()
	info.LastUpdateTime = now
	if info.CreateTime == 0 {
		info.CreateTime = now
	}
	db := database.GetMysql()
	if err := db.Save(info).Error; err != nil {
		return err
	}
	return nil
}

func (info *AdInfo) Get() (*AdInfo, error) {
	var adInfo AdInfo
	db := database.GetMysql()
	if err := db.Where(info).First(&adInfo).Error; err != nil {
		return &adInfo, err
	}
	return &adInfo, nil
}

func GetAllAdInfo() (*[]AdInfo, error) {
	var adInfo []AdInfo
	db := database.GetMysql()
	if err := db.Find(&adInfo).Error; err != nil {
		return &adInfo, err
	}
	return &adInfo, nil
}
