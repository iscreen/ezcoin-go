package v2

import (
	"ezcoin.cc/ezcoin-go/server/app/api"
	"github.com/gin-gonic/gin"
)

type UserRouter struct{}

func (u *UserRouter) InitUserRouter(Router *gin.RouterGroup) {
	userRouter := Router.Group("my")
	userApi := api.ApiGroupApp.ApiV1Group.UserApi
	{
		userRouter.GET("robots", userApi.GetUserRobots)
		userRouter.PATCH("robots/:id", userApi.UpdateUserRobot)
		userRouter.GET("robots/:id/robot_state", userApi.GetUserRobotState)
		userRouter.POST("robots", userApi.CreateUserRobot)
	}

}
