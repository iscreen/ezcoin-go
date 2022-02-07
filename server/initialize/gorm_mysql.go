package initialize

import (
	"ezcoin.cc/ezcoin-go/server/global"
	"ezcoin.cc/ezcoin-go/server/initialize/internal"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func GormMysql() *gorm.DB {
	m := global.GVA_CONFIG.Mysql
	if m.Dbname == "" {
		return nil
	}
	mysqlConfig := mysql.Config{
		DSN:                       m.Dsn(), // DSN data source name
		DefaultStringSize:         255,     // string 类型字段的默认长度
		DisableDatetimePrecision:  true,    // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,    // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,    // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false,   // 根据版本自动配置
	}
	if db, err := gorm.Open(mysql.New(mysqlConfig), internal.Gorm.Config()); err != nil {
		return nil
	} else {
		sqlDB, _ := db.DB()
		sqlDB.SetMaxIdleConns(m.MaxIdleConns)
		sqlDB.SetMaxOpenConns(m.MaxOpenConns)
		return db
	}
}

// // GormMysqlByConfig 初始化Mysql数据库用过传入配置
// func GormMysqlByConfig(m config.DB) *gorm.DB {
// 	if m.Dbname == "" {
// 		return nil
// 	}
// 	mysqlConfig := mysql.Config{
// 		DSN:                       m.Dsn(), // DSN data source name
// 		DefaultStringSize:         191,     // string 类型字段的默认长度
// 		SkipInitializeWithVersion: false,   // 根据版本自动配置
// 	}
// 	if db, err := gorm.Open(mysql.New(mysqlConfig), internal.Gorm.Config()); err != nil {
// 		panic(err)
// 	} else {
// 		sqlDB, _ := db.DB()
// 		sqlDB.SetMaxIdleConns(m.MaxIdleConns)
// 		sqlDB.SetMaxOpenConns(m.MaxOpenConns)
// 		return db
// 	}
// }
