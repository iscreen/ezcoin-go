package v1

import (
	"ezcoin.cc/ezcoin-go/server/app/api"
	"github.com/gin-gonic/gin"
)

type AuthRouter struct{}

func (a *AuthRouter) InitAuthRouter(Router *gin.RouterGroup) (R gin.IRouter) {
	Router.POST("/login", api.ApiGroupApp.ApiV1Group.Login)
	Router.POST("/sign_up", api.ApiGroupApp.ApiV1Group.Register)
	Router.POST("/password", api.ApiGroupApp.ApiV1Group.SendForgetPassword)
	Router.PATCH("/password", api.ApiGroupApp.ApiV1Group.ResetPassword)
	Router.POST("/confirmation", api.ApiGroupApp.ApiV1Group.ResendConfirmation)
	Router.PATCH("/confirmation", api.ApiGroupApp.ApiV1Group.Confirmation)
	return Router
}
