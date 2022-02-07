package dao

import "github.com/shopspring/decimal"

type UserRobot struct {
	ID                uint32          `json:"id"`
	UserId            uint32          `json:"userId"`
	CurrencyId        uint32          `json:"currencyId"`
	ApiKey            string          `json:"apiKey"`
	ApiSecret         string          `json:"apiSecret"`
	Activated         bool            `json:"activated"`
	Hidden            bool            `json:"hidden"`
	ReservedAmount    decimal.Decimal `json:"reservedAmount"`
	MaxAmount         decimal.Decimal `json:"maxAmount"`
	Period            uint            `json:"period"`
	FixRatePeriod     uint            `json:"fixRatePeriod"`
	Intervals         uint            `json:"intervals"`
	NumLastIntervals  uint            `json:"numLast_intervals"`
	MinRate           decimal.Decimal `json:"minRate"`
	FixRate           decimal.Decimal `json:"fixRate"`
	RandomRangeLow    float32         `json:"randomRangeLow"`
	RandomRangeHigh   float32         `json:"randomRangeHigh"`
	RangeAmount       decimal.Decimal `json:"rangeAmount"`
	RangeNum          uint            `json:"rangeNum"`
	RangeLowRate      decimal.Decimal `json:"rangeLowRate"`
	RangeHighRate     decimal.Decimal `json:"rangeHighRate"`
	RangeRatePeriod   uint            `json:"rangeRatePeriod"`
	LendingStrategyId uint32          `json:"lendingStrategyId"`
	RobotServerId     uint32          `json:"robotServerId"`
}
