package dao

import (
	"ezcoin.cc/ezcoin-go/server/app/model"
)

func (d *Dao) GetLendingStrategies() ([]*model.LendingStrategy, error) {
	ls := model.LendingStrategy{}
	return ls.GetLendingStrategies(d.engine)
}
