package dao

import (
	"ezcoin.cc/ezcoin-go/server/app/model"
	"ezcoin.cc/ezcoin-go/server/app/model/common/request"
	"ezcoin.cc/ezcoin-go/server/app/model/common/response"
	"github.com/shopspring/decimal"
)

type CreateScheduleOrder struct {
	UserId            uint32          `json:"userId"`
	FromCurrencyId    uint32          `json:"fromCurrencyId"`
	ToCurrencyId      uint32          `json:"toCurrencyId"`
	CancelFromFunding bool            `json:"cancelFromFunding"`
	TransferToFunding bool            `json:"transferToFunding"`
	FastTrade         bool            `json:"fastTrade"`
	Amount            decimal.Decimal `json:"amount"`
	Price             decimal.Decimal `json:"price"`
}

type UpdateScheduleOrder struct {
	ID                uint32          `json:"id"`
	UserId            uint32          `json:"userId"`
	FromCurrencyId    uint32          `json:"fromCurrencyId"`
	ToCurrencyId      uint32          `json:"toCurrencyId"`
	CancelFromFunding bool            `json:"cancelFromFunding"`
	TransferToFunding bool            `json:"transferToFunding"`
	FastTrade         bool            `json:"fastTrade"`
	Amount            decimal.Decimal `json:"amount"`
	Price             decimal.Decimal `json:"price"`
}

type DeleteScheduleOrder struct {
	ID     uint32 `json:"id" `
	UserId uint32 `json:"userId" `
}

func (d *Dao) GetScheduleOrders(param *request.TableQuery) (*response.PageResult, error) {
	so := model.ScheduleOrder{}
	return so.List(d.engine, param)
}

func (d *Dao) CreateScheduleOrders(param *CreateScheduleOrder) (*model.ScheduleOrder, error) {
	so := model.ScheduleOrder{
		UserId:            param.UserId,
		FromCurrencyId:    param.FromCurrencyId,
		ToCurrencyId:      param.ToCurrencyId,
		CancelFromFunding: param.CancelFromFunding,
		TransferToFunding: param.TransferToFunding,
		FastTrade:         param.FastTrade,
		Amount:            param.Amount,
		Price:             param.Price,
	}

	return so.Create(d.engine)
}

func (d *Dao) UpdateScheduleOrders(param *UpdateScheduleOrder) (*model.ScheduleOrder, error) {
	so := model.ScheduleOrder{
		UserId:            param.UserId,
		FromCurrencyId:    param.FromCurrencyId,
		ToCurrencyId:      param.ToCurrencyId,
		CancelFromFunding: param.CancelFromFunding,
		TransferToFunding: param.TransferToFunding,
		FastTrade:         param.FastTrade,
		Amount:            param.Amount,
		Price:             param.Price,
	}

	return so.Update(d.engine)
}

func (d *Dao) DeleteScheduleOrders(param *DeleteScheduleOrder) error {
	so := model.ScheduleOrder{
		EZ_MODEL: model.EZ_MODEL{
			ID: param.ID,
		},
		UserId: param.UserId,
	}

	return so.Delete(d.engine)
}
