package model

import (
	"fmt"

	database "gitlab.com/TenbayMCloud/awesome-raserver/common/dbMgr"
)

// ORICO-NAS-WRC10 => vendor-product-type
type AdRsrcInfo struct {
	Path       string `json:"path" gorm:"column:path;primary_key;not null"`
	DevVendor  string `json:"dev_vendor" gorm:"column:dev_vendor;size:64;default:null"`
	DevProduct string `json:"dev_product" gorm:"column:dev_product;size:64;default:null"`
	DevType    string `json:"dev_type" gorm:"column:dev_type;size:64;not null"` // 这个不为0
	AdURL      string `json:"ad_url" gorm:"column:ad_url;default:null"`
	RsrcURL    string `json:"rsrc_url" gorm:"column:rsrc_url;default:null"`
	Desc       string `json:"desc" gorm:"column:desc;default:null"`
	CanDel     int    `json:"can_del" gorm:"column:can_del;default:0"` // 0不能删除，这里是在自动生成广告信息的时候使用
	Ext        string `json:"ext" gorm:"column:ext;default:null"`
}

func (pinfo *AdRsrcInfo) IsExist() bool {
	var info AdRsrcInfo
	db := database.GetMysql()
	if err := db.Where(pinfo).First(&info).Error; err != nil {
		return false
	}
	return true
}

func (pinfo *AdRsrcInfo) AddNoDup() error {
	if pinfo.IsExist() {
		return fmt.Errorf("This AD Rsrc info already exist")
	}

	db := database.GetMysql()
	if err := db.Create(pinfo).Error; err != nil {
		return err
	}
	return nil
}

func (pinfo *AdRsrcInfo) Add() error {
	db := database.GetMysql()
	if err := db.Create(pinfo).Error; err != nil {
		return err
	}
	return nil
}

func (pinfo *AdRsrcInfo) DelWhere(wStr string) error {
	db := database.GetMysql()
	if err := db.Where(wStr).Delete(AdRsrcInfo{}).Error; err != nil {
		return err
	}
	return nil
}

// Del 主键为空时，删除所有
func (pinfo *AdRsrcInfo) Del() error {
	db := database.GetMysql()
	if err := db.Delete(pinfo).Error; err != nil {
		return err
	}
	return nil
}

func (pinfo *AdRsrcInfo) Update() error {
	db := database.GetMysql()
	if err := db.Save(pinfo).Error; err != nil {
		return err
	}
	return nil
}

func (pinfo *AdRsrcInfo) Get() (*[]AdRsrcInfo, error) {
	var info []AdRsrcInfo
	db := database.GetMysql()
	if err := db.Where(pinfo).Find(&info).Error; err != nil {
		return &info, err
	}
	return &info, nil
}
