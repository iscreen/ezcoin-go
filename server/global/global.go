package global

import (
	"ezcoin.cc/ezcoin-go/server/config"
	"github.com/go-redis/redis/v8"
	"github.com/songzhibin97/gkit/cache/local_cache"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"golang.org/x/sync/singleflight"
	"gorm.io/gorm"
)

const (
	ConfigEnv  = "GVA_CONFIG"
	ConfigFile = "config.yaml"
)

var (
	GVA_DB                  *gorm.DB
	GVA_CONFIG              config.Server
	GVA_VP                  *viper.Viper
	GVA_LOG                 *zap.Logger
	GVA_REDIS               *redis.Client
	GVA_Concurrency_Control = &singleflight.Group{}
	BlackCache              local_cache.Cache
)
