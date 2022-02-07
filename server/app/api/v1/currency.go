package v1

import (
	"fmt"

	"ezcoin.cc/ezcoin-go/server/app/model"
	"ezcoin.cc/ezcoin-go/server/app/model/common/response"
	"ezcoin.cc/ezcoin-go/server/app/service"
	"ezcoin.cc/ezcoin-go/server/global"
	"ezcoin.cc/ezcoin-go/server/pkg/app"
	"ezcoin.cc/ezcoin-go/server/pkg/convert"
	"ezcoin.cc/ezcoin-go/server/pkg/errcode"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type CurrencyApi struct{}

// @Summary Get all currencies
// @Tags Currency
// @Version  1.0
// @Accept   json
// @Produce  json
// @Security ApiKeyAuth
// @Success  200  {array}  response.Response{data=model.Currency}
// @Failure  404  {object}  response.Response
// @Failure  422  {object}  response.Response
// @Router /v1/currencies [get]
func (curr *CurrencyApi) GetCurrencies(c *gin.Context) {
	filter := c.DefaultQuery("filter[kind]", "")
	svc := service.New(c.Request.Context())
	currencies, err := svc.GetCurrencies(filter)
	if err != nil {
		global.GVA_LOG.Error("取幣種失敗!", zap.Error(err))
		response.FailWithMessage("取幣種失敗", c)
		return
	}

	response.OkWithData(currencies, c)
}

// @Summary Get Currency
// @Tags Currency
// @Version  1.0
// @Accept   json
// @Produce  json
// @Security ApiKeyAuth
// @Param    id   path  int  true  "Currency ID"
// @Success  200  {object}  response.Response{data=model.Currency}
// @Failure  404  {object}  response.Response
// @Failure  422  {object}  response.Response
// @Router /v1/currencies/{id} [get]
func (curr *CurrencyApi) GetCurrency(c *gin.Context) {
	var currency model.Currency
	id := convert.StrTo(c.Param("id")).MustInt()
	global.GVA_LOG.Debug(fmt.Sprintf("id: %d", id))
	svc := service.New(c.Request.Context())
	currency, err := svc.GetCurrency(id)
	if err != nil {
		global.GVA_LOG.Error("取幣種失敗!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithData(currency, c)
}

// @Summary Create Currency
// @Tags Currency
// @Version  1.0
// @Accept   json
// @Produce  json
// @Security ApiKeyAuth
// @Param    currency  body  service.CreateCurrency  true  "Currency"
// @Success  200  {object}  response.Response{data=model.Currency}
// @Failure  404  {object}  response.Response
// @Failure  422  {object}  response.Response
// @Router /v1/currencies [post]
func (curr *CurrencyApi) CreateCurrency(c *gin.Context) {
	var currReq service.CreateCurrency

	valid, errs := app.BindAndValid(c, &currReq)
	if !valid {
		global.GVA_LOG.Error(fmt.Sprintf("app.BindAndValid errs: %v", errs))
		response.FailWithError(errcode.InvalidParams.WithDetails(errs.Errors()...), c)
		return
	}
	svc := service.New(c.Request.Context())
	currency := &model.Currency{Kind: currReq.Kind, Name: currReq.Name, SymbolName: currReq.Name, MinFundingAmount: currReq.MinFundingAmount}
	err := svc.CreateCurrency(currency)
	if err != nil {
		global.GVA_LOG.Error("新增幣種失敗!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithData(currency, c)
}
