package v1

import (
	"ezcoin.cc/ezcoin-go/server/app/api"
	"github.com/gin-gonic/gin"
)

type LendingStrategyRouter struct{}

func (c *LendingStrategyRouter) InitLendingStrategyRouter(Router *gin.RouterGroup) {
	router := Router.Group("lending_strategies")
	api := api.ApiGroupApp.ApiV1Group.LendingStrategyApi
	{
		router.GET("", api.GetLendingStrategies)
	}
}
