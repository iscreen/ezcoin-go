package service

import (
	"fmt"

	"ezcoin.cc/ezcoin-go/server/app/dao"
	"ezcoin.cc/ezcoin-go/server/app/model"
	"ezcoin.cc/ezcoin-go/server/app/model/common/request"
	"ezcoin.cc/ezcoin-go/server/app/model/common/response"
	"ezcoin.cc/ezcoin-go/server/global"
	"github.com/shopspring/decimal"
)

type CreateScheduleOrder struct {
	UserId            uint32          `json:"userId" form:"UserId" validate:"required"`
	FromCurrencyId    uint32          `json:"fromCurrencyId" form:"fromCurrencyId" validate:"required"`
	ToCurrencyId      uint32          `json:"toCurrencyId" form:"toCurrencyId" validate:"required"`
	CancelFromFunding bool            `json:"cancelFromFunding" form:"cancelFromFunding"`
	TransferToFunding bool            `json:"transferToFunding" form:"transferToFunding"`
	FastTrade         bool            `json:"fastTrade" form:"fastTrade"`
	Amount            decimal.Decimal `json:"amount" form:"amount" binding:"gte=0" validate:"required"`
	Price             decimal.Decimal `json:"price" form:"price" binding:"gte=0" validate:"required" `
}

type DeleteScheduleOrder struct {
	ID     uint32 `json:"id" `
	UserId uint32 `json:"userId" `
}

//@function: GetScheduleOrders
//@description: Get schedule orders
//@param: queryReq request.TableQuery
//@return: response.PageResult, error
func (svc *Service) GetScheduleOrders(param *request.TableQuery) (*response.PageResult, error) {
	return svc.dao.GetScheduleOrders(param)
}

func (svc *Service) CreateScheduleOrder(param *CreateScheduleOrder) (*model.ScheduleOrder, error) {
	return svc.dao.CreateScheduleOrders(&dao.CreateScheduleOrder{
		UserId:            param.UserId,
		FromCurrencyId:    param.FromCurrencyId,
		ToCurrencyId:      param.ToCurrencyId,
		CancelFromFunding: param.CancelFromFunding,
		FastTrade:         param.FastTrade,
		TransferToFunding: param.TransferToFunding,
		Amount:            param.Amount,
		Price:             param.Price,
	})
}

func (svc *Service) DeleteScheduleOrder(param *DeleteScheduleOrder) error {
	global.GVA_LOG.Debug(fmt.Sprintf("DeleteScheduleOrder: %v", param))
	return svc.dao.DeleteScheduleOrders(&dao.DeleteScheduleOrder{
		ID:     param.ID,
		UserId: param.UserId,
	})
}
