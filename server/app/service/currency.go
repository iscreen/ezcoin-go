package service

import (
	"ezcoin.cc/ezcoin-go/server/app/model"
)

type CreateCurrency struct {
	Kind             uint    `form:"kind" json:"kind" binding:"required" validate:"required,gte=0,lte=2" minimum:"0" maximum:"2" default:"0"`
	Name             string  `form:"name" json:"name" binding:"required" validate:"required"`
	SymbolName       string  `form:"symbolName" json:"symbolName" binding:"required" validate:"required"`
	MinFundingAmount float32 `form:"minFundingAmount" json:"minFundingAmount" binding:"required" validate:"required"`
}

//@function: GetCurrencies
//@description: 取得所有幣種
//@param: kind string
//@return: []model.Currency, error
func (svc *Service) GetCurrencies(kind string) ([]*model.Currency, error) {
	var currencies []*model.Currency
	var err error

	if kind == "" {
		currencies, err = svc.dao.GetCurrencies()
		if err != nil {
			return currencies, err
		}
	} else {
		var kindInt uint
		if kind != "" {
			if kind == "funding" {
				kindInt = uint(model.FundingCurrency)
			} else {
				kindInt = uint(model.TradeCurrency)
			}
		}
		currencies, err = svc.dao.GetCurrencyByKind(kindInt)
		if err != nil {
			return currencies, err
		}
	}

	return currencies, nil
}

//@function: GetCurrency
//@description: 取得指定幣種
//@param: id int
//@return: []model.Currency, error
func (svc *Service) GetCurrency(id int) (model.Currency, error) {
	var currency model.Currency
	// if errors.Is(global.GVA_DB.Engine().First(&currency, id).Error, gorm.ErrRecordNotFound) {
	// 	return currency, errors.New("幣種不存在")
	// }
	return currency, nil
}

//@function: CreateCurrency
//@description: 取得指定幣種
//@param: currency model.Currency
//@return: err error
func (svc *Service) CreateCurrency(currency *model.Currency) (err error) {
	// if (!errors.Is(global.GVA_DB.Engine().First(&model.Currency{}, "kind = ? AND (name = ? or symbol_name = ?)", currency.Kind, currency.Name, currency.SymbolName).Error, gorm.ErrRecordNotFound)) {
	// 	return errors.New("存在相同的類型與名稱")
	// }
	// err = global.GVA_DB.Engine().Create(&currency).Error
	return err
}
