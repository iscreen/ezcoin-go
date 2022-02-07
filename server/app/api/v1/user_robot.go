package v1

import (
	"fmt"

	"ezcoin.cc/ezcoin-go/server/app/model/common/response"
	"ezcoin.cc/ezcoin-go/server/app/service"
	"ezcoin.cc/ezcoin-go/server/global"
	"ezcoin.cc/ezcoin-go/server/pkg/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type UserRobotApi struct{}

// @Summary Get all user's robots
// @Tags UserRobot
// @Version  1.0
// @Accept   json
// @Produce  json
// @Security ApiKeyAuth
// @Success  200  {array}  response.Response{data=model.UserRobot}
// @Failure  404  {object}  response.Response
// @Failure  422  {object}  response.Response
// @Router /v1/user_robots [get]
func (ur *UserRobotApi) GetRobots(c *gin.Context) {
	svc := service.New(c.Request.Context())
	robots, err := svc.GetUserRobots(utils.GetUserID(c))
	if err != nil {
		global.GVA_LOG.Error("取設定失敗!", zap.Error(err))
		response.FailWithMessage("取設定失敗", c)
		return
	}
	global.GVA_LOG.Debug(fmt.Sprintf("%v", robots))
	response.OkWithData(robots, c)
}
