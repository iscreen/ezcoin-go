package model

import (
	"fmt"
	"math"
	"time"

	"ezcoin.cc/ezcoin-go/server/app/model/common/request"
	"ezcoin.cc/ezcoin-go/server/app/model/common/response"
	"ezcoin.cc/ezcoin-go/server/pkg/app"
	"gorm.io/gorm"
)

type CurrencyKind int

const (
	FundingCurrency CurrencyKind = iota
	TradeCurrency
)

func (k CurrencyKind) String() string {
	states := [...]string{"FundingCurrency", "TradeCurrency"}
	if len(states) < int(k) {
		return ""
	}

	return states[k]
}

type Currency struct {
	EZ_MODEL

	Kind             uint    `json:"kind" gorm:"comment:類型"`
	Name             string  `json:"name" gorm:"comment:名稱"`
	SymbolName       string  `json:"symbolName" gorm:"comment:Symbol名稱"`
	MinFundingAmount float32 `json:"minFundingAmount" gorm:"comment:最小放貸金額"`
}

func (c Currency) List(db *gorm.DB, param request.TableQuery) (*response.PageResult, error) {
	var list []*Currency
	var result response.PageResult
	var count int64
	var query *gorm.DB

	query = db.Model(Currency{})

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

func (c Currency) Create(db *gorm.DB) (*Currency, error) {

	if err := db.Create(&c).Error; err != nil {
		return nil, err
	}
	return &c, nil
}

func (c Currency) Update(db *gorm.DB) (*Currency, error) {
	if err := db.Where("id = ?", c.ID).Updates(&c).Error; err != nil {
		return nil, err
	}
	return &c, nil
}

func (c Currency) Get(db *gorm.DB) (Currency, error) {
	var currency Currency
	err := db.Where("id = ?", c.ID).First(&currency).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return currency, err
	}
	return currency, nil
}

func (c Currency) Delete(db *gorm.DB) error {
	if err := db.Model(&c).Update("deleted_at", time.Now()).Error; err != nil {
		return err
	}
	return nil
}
func (c Currency) GetCurrencies(db *gorm.DB) ([]*Currency, error) {
	var currencies []*Currency
	if err := db.Find(&currencies).Error; err != nil {
		return currencies, err
	}
	return currencies, nil
}

func (c Currency) GetCurrenciesByKind(db *gorm.DB) ([]*Currency, error) {
	var currencies []*Currency
	if err := db.Where("kind = ?", c.Kind).Find(&currencies).Error; err != nil {
		return currencies, err
	}
	return currencies, nil
}
