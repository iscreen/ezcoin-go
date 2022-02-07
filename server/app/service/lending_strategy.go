package service

import (
	"ezcoin.cc/ezcoin-go/server/app/model"
)

//@function: GetLendingStrategies
//@description: 取得所有策略
//@return: []*model.LendingStrategy, error
func (svc *Service) GetLendingStrategies() ([]*model.LendingStrategy, error) {
	return svc.dao.GetLendingStrategies()
}
