package v1

import (
	"ezcoin.cc/ezcoin-go/server/app/api"
	"github.com/gin-gonic/gin"
)

type UserRouter struct{}

func (u *UserRouter) InitUserRouter(Router *gin.RouterGroup) {
	userRouter := Router.Group("my")
	userApi := api.ApiGroupApp.ApiV1Group.UserApi
	{
		userRouter.GET("", userApi.GetUser)
		userRouter.PATCH("", userApi.UpdateUser)
		userRouter.GET("last_month_margin_funding_payments", userApi.LastMonthMarginFundingPayments)
		userRouter.GET("last_month_margin_funding_payments_summary", userApi.LastMonthMarginFundingPaymentsSummary)
		userRouter.GET("robots", userApi.GetUserRobots)
	}

	schedulerOrderGroup := userRouter.Group("schedule_orders")
	scheduleOrderApi := api.ApiGroupApp.ApiV1Group.ScheduleOrderApi
	{
		schedulerOrderGroup.GET("", scheduleOrderApi.GetScheduleOrders)
		schedulerOrderGroup.POST("", scheduleOrderApi.CreateScheduleOrders)
		schedulerOrderGroup.DELETE(":id", scheduleOrderApi.DeleteScheduleOrders)
	}
}
