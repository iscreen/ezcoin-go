package main

import (
	"ezcoin.cc/ezcoin-go/server/core"
	"ezcoin.cc/ezcoin-go/server/global"
	"ezcoin.cc/ezcoin-go/server/initialize"
)

// @title           EZCoin API
// @version         1.0
// @description     This is a sample server celler server.
// @termsOfService  https://ezcoin.cc/terms/

// @contact.name   ServerMIS
// @contact.url    https://ezcoin.cc/support
// @contact.email  support@ezcoin.cc

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8888
// @BasePath  /api

// @securityDefinitions.apiKey  ApiKeyAuth
// @in header
// @name x-token

func main() {
	global.GVA_VP = core.Viper()
	global.GVA_LOG = core.Zap()
	global.GVA_DB = initialize.Gorm()

	if global.GVA_DB != nil {
		initialize.RegisterTables(global.GVA_DB) // 初始化Table
		// 程式結束前關閉資料庫連結
		db, _ := global.GVA_DB.DB()
		defer db.Close()
	}
	core.RunServer()
}
