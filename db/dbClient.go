package db

import (
	"codeSandbox/model"
	"codeSandbox/utils"
	"fmt"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"time"
)

var dBClinet *gorm.DB

func init() {
	// 初始化客户端

	log.Info("init database...")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		utils.Config.Server.Database.User, utils.Config.Server.Database.PassWord, utils.Config.Server.Database.Host, utils.Config.Server.Database.Port, utils.Config.Server.Database.Name)
	var err error
	dBClinet, err = gorm.Open(mysql.New(mysql.Config{
		DSN:                       dsn,
		DefaultStringSize:         256,  // string 类型字段的默认长度
		DisableDatetimePrecision:  true, // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true, // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true, // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: true, // 根据当前 MySQL 版本自动配置
	}), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		DisableForeignKeyConstraintWhenMigrating: true, // 数据迁移时不生成外键
	})
	if err != nil {
		log.Panic(fmt.Sprintf("database connect fail,%s", err))
	}
	sqlDB, err := dBClinet.DB()
	if err != nil {
		log.Panic(fmt.Sprintf("database init config fail,%s", err))
	}
	// TODO 添加所有的结构体
	err = dBClinet.AutoMigrate(&model.User{})
	err = dBClinet.AutoMigrate(&model.KeyPair{})
	err = dBClinet.AutoMigrate(&model.Execution{})

	if err != nil {
		log.Panic(fmt.Sprintf("database migrate fail,%s", err))
	}
	// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxIdleConns(10)
	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(100)
	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(10 * time.Second)
	//Db.Create(&User{Username: "root", Password: "123456", Role: 0})
	log.Info("init database success!")

}
