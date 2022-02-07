package v1

import (
	"fmt"

	"ezcoin.cc/ezcoin-go/server/app/model/common/response"
	"ezcoin.cc/ezcoin-go/server/app/service"
	"ezcoin.cc/ezcoin-go/server/global"
	"ezcoin.cc/ezcoin-go/server/pkg/app"
	"ezcoin.cc/ezcoin-go/server/pkg/convert"
	"ezcoin.cc/ezcoin-go/server/pkg/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type UserApi struct{}

// @Summary Get user last month margin funding payments
// @Tags User
// @Version 1.0
// @Accept  json
// @Produce json
// @Security ApiKeyAuth
// @Success  200  {object}   response.Response{data=map[string][][]interface{}}
// @Failure  404  {object}  response.Response
// @Failure  422  {object}  response.Response
// @Router /v1/v1/my/last_month_margin_funding_payments [get]
func (u *UserApi) LastMonthMarginFundingPayments(c *gin.Context) {
	svc := service.New(c.Request.Context())
	result, err := svc.LastMonthMarginFundingPayments(utils.GetUserID(c))
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(result, c)
}

// @Summary Get user last month margin funding payments summary
// @Tags User
// @Version 1.0
// @Accept  json
// @Produce json
// @Security ApiKeyAuth
// @Success  200  {object}   response.Response{data=map[string]float32}
// @Failure  404  {object}  response.Response
// @Failure  422  {object}  response.Response
// @Router /v1/v1/my/last_month_margin_funding_payments_summary [get]
func (u *UserApi) LastMonthMarginFundingPaymentsSummary(c *gin.Context) {
	svc := service.New(c.Request.Context())
	result, err := svc.LastMonthMarginFundingPaymentsSummary(utils.GetUserID(c))
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(result, c)
}

// @Summary Get user
// @Tags User
// @Version 1.0
// @Accept  json
// @Produce json
// @Security ApiKeyAuth
// @Success  200  {object}   response.Response{data=model.User}
// @Failure  404  {object}  response.Response
// @Failure  422  {object}  response.Response
// @Router /v1/my/last_month_margin_funding_payments_summary [get]
func (u *UserApi) GetUser(c *gin.Context) {
	svc := service.New(c.Request.Context())
	result, err := svc.GetUser(utils.GetUserID(c))
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(result, c)
}

// @Summary Update user
// @Tags User
// @Version 1.0
// @Accept  json
// @Produce json
// @Security ApiKeyAuth
// @Success  200  {object}   response.Response{data=model.User}
// @Failure  404  {object}  response.Response
// @Failure  422  {object}  response.Response
// @Router /v1/my/last_month_margin_funding_payments_summary [get]
func (u *UserApi) UpdateUser(c *gin.Context) {
	param := service.UpdateUser{
		Id: utils.GetUserID(c),
	}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.GVA_LOG.Error(fmt.Sprintf("app.BindAndValid errs: %v", errs))
		response.FailWithDetailed(errs.ErrorFields(), "message", c)
		return
	}
	global.GVA_LOG.Debug("update user")
	svc := service.New(c.Request.Context())
	user, err := svc.UpdateUser(&param)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(user, c)
}

// @Summary Get all user's robots
// @Tags User
// @Version  1.0
// @Accept   json
// @Produce  json
// @Security ApiKeyAuth
// @Success  200  {array}  response.Response{data=model.UserRobot}
// @Failure  404  {object}  response.Response
// @Failure  422  {object}  response.Response
// @Router /v1/my/robots [get]
func (u *UserApi) GetUserRobots(c *gin.Context) {
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

// @Summary Create user's robot
// @Tags User
// @Version  1.0
// @Accept   json
// @Produce  json
// @Security ApiKeyAuth
// @Success  200  {object}  response.Response{data=model.UserRobot}
// @Failure  404  {object}  response.Response
// @Failure  422  {object}  response.Response
// @Router /v1/my/robots [post]
func (u *UserApi) CreateUserRobot(c *gin.Context) {
	param := &service.CreateUserRobot{UserId: utils.GetUserID(c)}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		response.FailWithDetailed(errs.ErrorFields(), "fields error", c)
		return
	}

	svc := service.New(c.Request.Context())
	robots, err := svc.CreateUserRobot(param)
	if err != nil {
		global.GVA_LOG.Error("新建失敗!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(robots, c)
}

// @Summary Update user's robot
// @Tags User
// @Version  1.0
// @Accept   json
// @Produce  json
// @Security ApiKeyAuth
// @Success  200  {object}  response.Response{data=model.UserRobot}
// @Failure  404  {object}  response.Response
// @Failure  422  {object}  response.Response
// @Router /v1/my/robots/{id} [patch]
func (u *UserApi) UpdateUserRobot(c *gin.Context) {
	param := service.UpdateUserRobot{
		ID:     convert.StrTo(c.Param("id")).MustUint32(),
		UserId: utils.GetUserID(c),
	}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.GVA_LOG.Error(fmt.Sprintf("app.BindAndValid errs: %v", errs))
		response.FailWithDetailed(errs.ErrorFields(), "fields error", c)
		return
	}

	svc := service.New(c.Request.Context())
	robots, err := svc.UpdateUserRobot(&param)
	if err != nil {
		global.GVA_LOG.Error("更新失敗!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(robots, c)
}

// @Summary Get user robot state
// @Tags User
// @Version  1.0
// @Accept   json
// @Produce  json
// @Security ApiKeyAuth
// @Success  200  {object}  response.Response
// @Failure  404  {object}  response.Response
// @Failure  422  {object}  response.Response
// @Router /v1/my/robots/{id}/state [get]
func (u *UserApi) GetUserRobotState(c *gin.Context) {
	id := c.Param("id")
	svc := service.New(c.Request.Context())
	state, err := svc.GetUserRobotState(uint32(convert.StrTo(id).MustInt()))
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(state, c)
}
