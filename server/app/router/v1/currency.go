package v1

import (
	"ezcoin.cc/ezcoin-go/server/app/api"
	"github.com/gin-gonic/gin"
)

type CurrencyRouter struct{}

func (c *CurrencyRouter) InitCurrencyRouter(Router *gin.RouterGroup) {
	currencyRouter := Router.Group("currencies")
	currencyApi := api.ApiGroupApp.ApiV1Group.CurrencyApi
	{
		currencyRouter.GET("", currencyApi.GetCurrencies)
		currencyRouter.GET(":id", currencyApi.GetCurrency)
		currencyRouter.POST("", currencyApi.CreateCurrency)
	}
}
