package v1

import (
	"fmt"

	"ezcoin.cc/ezcoin-go/server/app/model/common/request"
	"ezcoin.cc/ezcoin-go/server/app/model/common/response"
	"ezcoin.cc/ezcoin-go/server/app/service"
	"ezcoin.cc/ezcoin-go/server/global"
	"ezcoin.cc/ezcoin-go/server/pkg/app"
	"ezcoin.cc/ezcoin-go/server/pkg/convert"
	"ezcoin.cc/ezcoin-go/server/pkg/utils"
	"github.com/gin-gonic/gin"
)

type ScheduleOrderApi struct{}

// @Summary Get schedule orders
// @Tags ScheduleOrder
// @Version  1.0
// @Accept   json
// @Produce  json
// @Security ApiKeyAuth
// @Param page  body  int  false  "Page"
// @Param pageSize  body  int  false  "Page Size"
// @Param order  body  string  false  "Order"
// @Param sort  body  string  false  "Sort column"
// @Param filter  body  map[string]string  false  "Filter column and value"
// @Success  200  {array}  response.Response{data=model.ScheduleOrder}
// @Failure  404  {object}  response.Response
// @Failure  422  {object}  response.Response
// @Router /v1/my/schedule_orders [get]
func (so *ScheduleOrderApi) GetScheduleOrders(c *gin.Context) {
	var param request.TableQuery
	err := app.BindTableQuery(c, &param)
	if err != nil {
		global.GVA_LOG.Debug(fmt.Sprintf("err => %v", err))
	}
	svc := service.New(c.Request.Context())
	orders, err := svc.GetScheduleOrders(&param)
	if err != nil {
		response.FailWithMessage("取排程訂單失敗", c)
		return
	}
	response.OkWithData(orders, c)
}

// @Summary Create schedule orders
// @Tags ScheduleOrder
// @Version  1.0
// @Accept   json
// @Produce  json
// @Security ApiKeyAuth
// @Success  200  {Object}  response.Response{data=model.ScheduleOrder}
// @Failure  404  {object}  response.Response
// @Failure  422  {object}  response.Response
// @Router /v1/my/schedule_orders [get]
func (so *ScheduleOrderApi) CreateScheduleOrders(c *gin.Context) {
	param := service.CreateScheduleOrder{
		UserId: utils.GetUserID(c),
	}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		response.FailWithDetailed(errs.ErrorFields(), "fields error", c)
		return
	}
	svc := service.New(c.Request.Context())
	orders, err := svc.CreateScheduleOrder(&param)
	if err != nil {
		response.FailWithMessage("建立排程訂單失敗", c)
		return
	}
	response.OkWithData(orders, c)
}

// @Summary Create schedule orders
// @Tags ScheduleOrder
// @Version  1.0
// @Accept   json
// @Produce  json
// @Security ApiKeyAuth
// @Success  200  {Object}  response.Response{data=model.ScheduleOrder}
// @Failure  404  {object}  response.Response
// @Failure  422  {object}  response.Response
// @Router /v1/my/schedule_orders/{id} [delete]
func (so *ScheduleOrderApi) DeleteScheduleOrders(c *gin.Context) {
	svc := service.New(c.Request.Context())
	err := svc.DeleteScheduleOrder(&service.DeleteScheduleOrder{
		ID:     convert.StrTo(c.Param("id")).MustUint32(),
		UserId: utils.GetUserID(c),
	})
	if err != nil {
		response.FailWithMessage("刪除排程訂單失敗", c)
		return
	}
	response.OkWithMessage("成功刪除排程訂單", c)
}
