package core

import (
	"flag"
	"fmt"
	"os"
	"time"

	"ezcoin.cc/ezcoin-go/server/global"
	"github.com/fsnotify/fsnotify"
	"github.com/songzhibin97/gkit/cache/local_cache"
	"github.com/spf13/viper"
)

func Viper(path ...string) *viper.Viper {
	var config string
	if len(path) == 0 {
		flag.StringVar(&config, "c", "", "choose config file.")
		flag.Parse()
		if config == "" { // 優先等級：命令列 > 環境變數 > 預設值
			if configEnv := os.Getenv(global.ConfigEnv); configEnv == "" {
				config = global.ConfigFile
				fmt.Printf("您正在使用設定的預設值，設定檔的路徑為 %v\n", global.ConfigFile)
			} else {
				config = configEnv
				fmt.Printf("您正在使用 GVA_CONFIG 環境變數，設定檔的路徑為 %v\n", config)
			}
		} else {
			fmt.Printf("您正在使用命令列的 -c 的值，設定檔的路徑為 %v\n", config)
		}
	} else {
		config = path[0]
		fmt.Printf("你正在使用 func Viper() 傳遞的值，設定檔的路徑為 %v\n", config)
	}

	v := viper.New()
	v.SetConfigFile(config)
	v.SetConfigType("yaml")
	err := v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err))
	}
	v.WatchConfig()

	v.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("config file changed:", e.Name)
		if err := v.Unmarshal(&global.GVA_CONFIG); err != nil {
			fmt.Println(err)
		}
	})
	if err := v.Unmarshal(&global.GVA_CONFIG); err != nil {
		fmt.Println(err)
	}
	// // root 适配性
	// // 根据root位置去找到对应迁移位置,保证root路径有效
	// global.GVA_CONFIG.AutoCode.Root, _ = filepath.Abs("..")
	global.BlackCache = local_cache.NewCache(
		local_cache.SetDefaultExpire(time.Second * time.Duration(global.GVA_CONFIG.JWT.ExpiresTime)),
	)
	return v
}
