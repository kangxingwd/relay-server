package model

import (
	database "gitlab.com/TenbayMCloud/awesome-raserver/common/dbMgr"
	"time"
)

type Device struct {
	Devid			string			`gorm:"column:dev_id;unique_index;size:128;not null;primary_key" `
	Mac				string			`gorm:"column:mac;size:32;not null"`
	HostName		string			`gorm:"column:hostname;size:128;default:null"`
	Vendor			string			`gorm:"column:vendor;size:64;default:null"`
	SoftModel		string			`gorm:"column:soft_model;size:64;default:null"`
	SoftVersion		string			`gorm:"column:soft_version;size:64;default:null"`
	HardwareModel	string			`gorm:"column:hardware_model;size:64;default:null"`
	JoinTime		time.Time
	Role			uint32			`gorm:"column:role;defalt:0"`
	Ext				string			`gorm:"column:ext;default:null" json:"Ext"`
}

func (dev *Device) Add() error {
	db := database.GetMysql()
	if err := db.Create(dev).Error; err != nil {
		return err
	}
	return nil
}

func (dev *Device) Update() error {
	db := database.GetMysql()
	if err := db.Save(dev).Error; err != nil {
		return err
	}
	return nil
}

func DelDev(devId string) error {
	db := database.GetMysql()
	if err := db.Where("dev_id=?", devId).Delete(&Device{}).Error; err != nil {
		return err
	}
	return nil
}

func (dev *Device) UpdateRole(devFlag uint32) error {
	db := database.GetMysql()
	role := dev.Role | devFlag
	if err := db.Model(dev).Update("role", role).Error; err != nil {
		return err
	}
	return nil
}

func (dev *Device) DelRelayRole() error {
	db := database.GetMysql()
	role := dev.Role & 2
	if err := db.Model(dev).Update("role", role).Error; err != nil {
		return err
	}
	return nil
}

func IsExistDev(devId string) bool {
	var dev Device
	db := database.GetMysql()
	if err := db.Where(&Device{Devid:devId}).First(&dev).Error; err != nil {
		return false
	}
	return true
}

func GetDev(devId string) (*Device, error) {
	var dev Device
	db := database.GetMysql()
	if err := db.Where(&Device{Devid:devId}).First(&dev).Error; err != nil {
		return &dev, err
	}
	return &dev, nil
}

func GetDomainByMacList(macList []string) (*map[string]string, error) {
	macHostnameMap := make(map[string]string)

	db := database.GetMysql()
	rows,err := db.Table("devices").Where("mac in (?)", macList).Select("mac, hostname").Rows()
	defer rows.Close()
	if err != nil {
		return &macHostnameMap, err
	}

	var mac,hostname string
	for rows.Next() {
		rows.Scan(&mac, &hostname) // error
		macHostnameMap[mac] = hostname
	}
	return &macHostnameMap,nil
}
