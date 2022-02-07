package v1

import (
	"ezcoin.cc/ezcoin-go/server/app/api"
	"github.com/gin-gonic/gin"
)

type ScheduleOrderRouter struct{}

func (c *ScheduleOrderRouter) InitScheduleOrderRouter(Router *gin.RouterGroup) {
	router := Router.Group("schedule_orders")
	api := api.ApiGroupApp.ApiV1Group.ScheduleOrderApi
	{
		router.GET("", api.GetScheduleOrders)
		router.POST("", api.CreateScheduleOrders)
	}
}
