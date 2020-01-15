package model

import (
	logger "github.com/cihub/seelog"
	database "gitlab.com/TenbayMCloud/awesome-raserver/common/dbMgr"
)

// Migrate the schema of database if needed
func AutoMigrate() (err error) {
	db := database.GetMysql()

	// 这里建表
	if err := db.AutoMigrate(&Device{}, &Relay{}, &Mcloud{}, &Rcmap{}, &AdInfo{}, &AdRsrcInfo{}).Error; err != nil {
		logger.Errorf("AutoMigrate error: %v", err)
		return err
	}
	return nil
}

