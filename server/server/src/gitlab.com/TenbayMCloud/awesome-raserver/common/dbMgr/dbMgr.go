package mysql

import (
	"fmt"
	log "github.com/cihub/seelog"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	cfg "gitlab.com/TenbayMCloud/awesome-raserver/config"
)

var mysqlDB *gorm.DB

func  MysqlConnectInit() error {
	url := fmt.Sprintf("%s:%s@(%s)/%s?charset=utf8&parseTime=True&loc=Local&multiStatements=true",
		cfg.GetMysql().User, cfg.GetMysql().Password, cfg.GetMysql().Address, cfg.GetMysql().DBName)

	log.Infof("usrl: %s\n", url)
	db, err := gorm.Open("mysql", url)
	if err != nil {
		log.Infof("Open error!error: %v\n", err)
		return err
	}
	
	//db = db.Set("gorm:table_options", "ENGING=InnoDB CHARSET=utf8mb4 auto_increment=1")
	db.LogMode(cfg.GetMysql().LogMode)
	db.DB().SetMaxIdleConns(cfg.GetMysql().MaxIdleConns)
	db.DB().SetMaxOpenConns(cfg.GetMysql().MaxOpenConns)
	mysqlDB = db
	return nil
}
func GetMysql() *gorm.DB {
	return mysqlDB
}

func MysqlClose() error {
	return mysqlDB.Close()
}
