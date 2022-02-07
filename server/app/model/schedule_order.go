package model

import (
	"fmt"
	"math"
	"time"

	"ezcoin.cc/ezcoin-go/server/app/model/common/request"
	"ezcoin.cc/ezcoin-go/server/app/model/common/response"
	"ezcoin.cc/ezcoin-go/server/pkg/app"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type ScheduleOrder struct {
	EZ_MODEL
	UserId            uint32          `json:"userId"`
	User              User            `json:"user"`
	FromCurrencyId    uint32          `json:"fromCurrencyId"`
	FromCurrency      Currency        `json:"fromCurrency"`
	ToCurrencyId      uint32          `json:"toCurrencyId"`
	ToCurrency        Currency        `json:"toCurrency"`
	CancelFromFunding bool            `json:"cancelFromFunding" gorm:"type:tinyint,comment:取消放貸"`
	TransferToFunding bool            `json:"transferToFunding" gorm:"type:tinyint,comment:取消放貸"`
	OrderId           uint            `json:"orderId" gorm:"<-:update;type:bigint(20),default:null,comment:Bitfinex訂單編號"`
	Status            string          `json:"status" gorm:"<-:update"`
	Meta              string          `json:"meta" gorm:"<-:update;type:longtext,comment:取消放貸"`
	FastTrade         bool            `json:"fastTrade" gorm:"type:tinyint,comment:快速交易"`
	Amount            decimal.Decimal `json:"amount" gorm:"type:decimal(20,10),comment:金額"`
	Price             decimal.Decimal `json:"price" gorm:"type:decimal(20,10),comment:加密貨幣價格"`
	OrderStatus       string          `json:"orderStatus" gorm:"<-:update;type:varchar(255),comment:訂單狀態"`
	OrderType         string          `json:"orderType" gorm:"<-:update;type:varchar(255),comment:訂單類型"`
	OrderSymbol       string          `json:"orderSymbol" gorm:"<-:update;type:varchar(255),comment:訂單幣別"`
	ProcessMessage    string          `json:"processMessage" gorm:"<-:update;type:text,comment:處理訊息"`
}

func (so ScheduleOrder) List(db *gorm.DB, param *request.TableQuery) (*response.PageResult, error) {
	var list []*ScheduleOrder
	var result response.PageResult
	var count int64
	var query *gorm.DB

	query = db.Model(ScheduleOrder{}).Preload("FromCurrency").Preload("ToCurrency")

	// Filter
	for k, v := range param.Filter {
		query = query.Where(fmt.Sprintf("%s = ?", k), v)
	}
	query = query.Count(&count)

	if query.Error != nil {
		return &result, query.Error
	}
	query = query.Limit(param.PageSize)
	if param.Sort != "" && param.Order != "" {
		query = query.Order(fmt.Sprintf("%s %s", param.Sort, param.Order))
	}
	query = query.Offset(app.GetPageOffset(param.Page, param.PageSize))
	err := query.Find(&list).Error
	if err != nil {
		return &result, err
	}
	result.List = list
	result.Meta = response.PageMeta{
		Total:     count,
		TotalPage: int(math.Ceil(float64(count) / float64(param.PageSize))),
		Page:      param.Page,
		PageSize:  param.PageSize,
	}

	return &result, err
}

func (o ScheduleOrder) Create(db *gorm.DB) (*ScheduleOrder, error) {
	if err := db.Create(&o).Error; err != nil {
		return nil, err
	}
	return &o, nil
}

func (o ScheduleOrder) Update(db *gorm.DB) (*ScheduleOrder, error) {
	if err := db.Where("id = ?", o.ID).Updates(&o).Error; err != nil {
		return nil, err
	}
	return &o, nil
}

func (o ScheduleOrder) Get(db *gorm.DB) (ScheduleOrder, error) {
	var order ScheduleOrder
	err := db.Where("id = ? AND deleted_at IS NOT NULL", o.ID).
		Preload("FromCurrency").Preload("ToCurrency").First(&order).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return order, err
	}
	return order, nil
}

func (o ScheduleOrder) Delete(db *gorm.DB) error {
	if err := db.Model(&o).Where("user_id = ?", o.UserId).Update("deleted_at", time.Now()).Error; err != nil {
		return err
	}
	// if err := db.Delete(&o).Error; err != nil {
	// 	return err
	// }
	return nil
}

func (u *ScheduleOrder) AfterCreate(tx *gorm.DB) (err error) {
	// global.GVA_LOG.Debug("create job to submit order")
	return
}
