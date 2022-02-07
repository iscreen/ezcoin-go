package dao

import (
	"ezcoin.cc/ezcoin-go/server/app/model"
)

func (d *Dao) GetCurrencyByKind(kind uint) ([]*model.Currency, error) {
	currency := model.Currency{Kind: kind}
	return currency.GetCurrenciesByKind(d.engine)
}

func (d *Dao) GetCurrencies() ([]*model.Currency, error) {
	currency := model.Currency{}
	return currency.GetCurrencies(d.engine)
}
