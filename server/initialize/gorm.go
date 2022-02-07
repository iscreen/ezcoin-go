package initialize

import (
	"os"

	"ezcoin.cc/ezcoin-go/server/app/model"
	"ezcoin.cc/ezcoin-go/server/global"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func Gorm() *gorm.DB {
	return GormMysql()
}

func RegisterTables(db *gorm.DB) {
	// 暫時空著

	err := db.AutoMigrate(
		model.ScheduleOrder{},
	)
	if err != nil {
		global.GVA_LOG.Error("register table failed", zap.Error(err))
		os.Exit(0)
	}
	global.GVA_LOG.Info("register table success")
}
