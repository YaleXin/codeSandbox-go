package db

import (
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var dBClinet *gorm.DB

func init() {
	// 初始化客户端
	dsn := "root:123456@tcp(localhost:3307)/code_sandbox_data?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	dBClinet, err = gorm.Open(mysql.New(mysql.Config{
		DSN: dsn,
	}))
	if err != nil {

	}
	log.Infof("******* db init %v ******* ", dBClinet)

	// 初始化表

}
