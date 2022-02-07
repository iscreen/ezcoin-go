package initialize

import (
	"ezcoin.cc/ezcoin-go/server/app/middleware"
	"ezcoin.cc/ezcoin-go/server/docs"
	"ezcoin.cc/ezcoin-go/server/global"

	"ezcoin.cc/ezcoin-go/server/app/router"
	val "ezcoin.cc/ezcoin-go/server/pkg/validator"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"

	// "github.com/go-playground/validator/v10"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Routers() *gin.Engine {
	Router := gin.Default()

	// if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
	// 	v.RegisterValidation("bookabledate", valid.BookableDate)
	// }
	// var Validator *val.CustomValidator
	Validator := val.NewCustomValidator()
	Validator.Engine()
	binding.Validator = Validator

	// Router.Use(middleware.Cors())
	global.GVA_LOG.Info("router register success")
	Router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	docs.SwaggerInfo.BasePath = "/api"

	v1Router := router.RouterGroupApp.V1RouterGroup
	v2Router := router.RouterGroupApp.V2RouterGroup

	PublicApiGroup := Router.Group("/api")
	{
		// 健庫監測
		PublicApiGroup.GET("/health", func(c *gin.Context) {
			c.JSON(200, "ok")
		})
	}

	PrivateApiV1Group := PublicApiGroup.Group("/v1")
	PrivateApiV1Group.Use(middleware.JWTAuth(), middleware.Translations())
	{
		v1Router.InitCurrencyRouter(PrivateApiV1Group)
		v1Router.InitUserRouter(PrivateApiV1Group)
		v1Router.InitLendingStrategyRouter(PrivateApiV1Group)
		v1Router.InitScheduleOrderRouter(PrivateApiV1Group)
	}

	PublicApiV1Group := PublicApiGroup.Group("/v1")
	PublicApiGroup.Use(middleware.Translations())
	{
		router.RouterGroupApp.V1RouterGroup.InitAuthRouter(PublicApiV1Group)
	}

	PrivateApiV2Group := PublicApiGroup.Group("/v2")
	PrivateApiV2Group.Use(middleware.JWTAuth(), middleware.Translations())
	{
		v2Router.InitUserRouter(PrivateApiV2Group)
	}

	return Router
}
