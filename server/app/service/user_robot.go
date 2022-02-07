package service

import (
	"errors"

	"ezcoin.cc/ezcoin-go/server/app/model"
	"github.com/shopspring/decimal"
)

type RandomRange struct {
	RandomRangeLow  float32 `json:"randomRangeLow" binding:"lt=0,ltfield=RandomRangeHigh"`
	RandomRangeHigh float32 `json:"randomRangeHigh" binding:"lt=0"`
}

type RangeAnnualRate struct {
	RangeAnnualHighRate float32 `json:"rangeAnnualHighRate"`
	RangeAnnualLowRate  float32 `json:"rangeAnnualLowRate"`
}

type RangeRate struct {
	RangeHighRate decimal.Decimal `json:"rangeHighRate"`
	RangeLowRate  decimal.Decimal `json:"rangeLowRate"`
}

type CreateUserRobot struct {
	CurrencyId        uint32          `json:"currencyId" binding:"required" validate:"required"`
	UserId            uint32          `json:"userId" binding:"required" validate:"required"`
	ApiKey            string          `json:"apiKey" binding:"required" validate:"required"`
	ApiSecret         string          `json:"apiSecret" binding:"required" validate:"required"`
	Activated         bool            `json:"activated" validate:"required"`
	Hidden            bool            `json:"hidden"`
	ReservedAmount    decimal.Decimal `json:"reservedAmount" binding:"gte=0"`
	MaxAmount         decimal.Decimal `json:"maxAmount" binding:"gte=0"`
	LendingStrategyId uint32          `json:"lendingStrategyId" binding:"required" validate:"required"`
	Period            uint            `json:"period" binding:"gte=0"`
	FixRatePeriod     uint            `json:"fixRatePeriod" binding:"gte=0"`
	Intervals         uint            `json:"intervals" binding:"gte=0"`
	NumLastIntervals  uint            `json:"numLastIntervals"`
	MinRate           decimal.Decimal `json:"minRate" binding:"required,gt=0"`
	FixRate           decimal.Decimal `json:"fixRate" binding:"required,gt=0"`
	RandomRange       RandomRange     `json:"randomRanges" binding:"required"`
	RangeAmount       decimal.Decimal `json:"rangeAmount"`
	RangeNum          uint            `json:"rangeNum"`
	RangeAnnualRate   RangeAnnualRate `json:"rangeAnnualRates"`
	RangeRate         RangeRate       `json:"rangeRates"`
	RangeRatePeriod   uint            `json:"rangeRatePeriod"`
}

type UpdateUserRobot struct {
	ID                uint32          `form:"id" binding:"required,gte=1"`
	UserId            uint32          `json:"userId" binding:"required" validate:"required"`
	CurrencyId        uint32          `json:"currencyId" binding:"required" validate:"required"`
	ApiKey            string          `json:"apiKey" binding:"required" validate:"required"`
	ApiSecret         string          `json:"apiSecret" binding:"required" validate:"required"`
	Activated         bool            `json:"activated" validate:"required"`
	Hidden            bool            `json:"hidden"`
	ReservedAmount    decimal.Decimal `json:"reservedAmount" binding:"gte=0"`
	MaxAmount         decimal.Decimal `json:"maxAmount" binding:"gte=0"`
	LendingStrategyId uint32          `json:"lendingStrategyId" binding:"required" validate:"required"`
	Period            uint            `json:"period" binding:"gte=0"`
	FixRatePeriod     uint            `json:"fixRatePeriod" binding:"gte=0"`
	Intervals         uint            `json:"intervals" binding:"gte=0"`
	NumLastIntervals  uint            `json:"numLastIntervals"`
	MinRate           decimal.Decimal `json:"minRate" binding:"required,gt=0"`
	FixRate           decimal.Decimal `json:"fixRate" binding:"required,gt=0"`
	RandomRange       RandomRange     `json:"randomRanges" binding:"required"`
	RangeAmount       decimal.Decimal `json:"rangeAmount"`
	RangeNum          uint            `json:"rangeNum"`
	RangeAnnualRate   RangeAnnualRate `json:"rangeAnnualRates"`
	RangeRate         RangeRate       `json:"rangeRates"`
	RangeRatePeriod   uint            `json:"rangeRatePeriod"`
}

//@function: GetUserRobots
//@description: Get user's robots
//@param: userId
//@return: []model.UserRobot, err error
func (svc *Service) GetRobots(userId uint) ([]model.UserRobot, error) {
	// var robots []model.UserRobot
	// err := global.GVA_DB.Engine().Where("user_id = ?", userId).Preload("Currency").Preload("RobotServer").Find(&robots).Error
	// if err != nil {
	// 	return robots, err
	// }

	// return robots, err
	return nil, errors.New("")
}

//@function: GetUserRobot
//@description: Get user's robot
//@param: userId
//@return: model.UserRobot, err error
func (svc *Service) GetUserRobot(userId uint, id int) (model.UserRobot, error) {
	// var robot model.UserRobot
	// err := svc.dao.Where("id = ? AND user_id = ?", id, userId).Preload("Currency").Preload("RobotServer").Find(&robot).Error
	// if err != nil {
	// 	return robot, err
	// }
	// robot.MinRate = robot.MinRate.Mul(decimal.NewFromInt(100))
	// return robot, err
	return model.UserRobot{}, nil
}

//@function: UserRobotExist
//@description: check user's robot exist
//@param: userId uint, id int, currencyId int
//@return: err error
func (svc *Service) UserRobotExist(userId uint, id int, currencyId int) error {
	// var robot model.UserRobot
	// // Check currency exists?
	// if !errors.Is(global.GVA_DB.Engine().
	// 	Where("currency_id = ? AND user_id = ? AND id != ?", currencyId, userId, id).Find(&robot).Error, gorm.ErrRecordNotFound) {
	// 	return errors.New("幣別已存在")
	// }
	// return nil
	return nil
}
