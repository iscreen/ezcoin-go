package dao

import (
	"fmt"

	"ezcoin.cc/ezcoin-go/server/app/model"
	"ezcoin.cc/ezcoin-go/server/global"
)

type UpdateUser struct {
	Id        uint32 `json:"id"`
	ApiKey    string `json:"apiKey"`
	ApiSecret string `json:"apiSecret"`
}

type Login struct {
	Login     string `json:"login"`
	Password  string `json:"password"`
	ClientIP  string `json:"-"`
	UserAgent string `json:"-"`
}

type Confirmation struct {
	ConfirmationToken string `json:"confirmationToken"`
	Password          string `json:"password"`
}

type ResetPassword struct {
	ResetPasswordToken string `json:"resetPasswordToken"`
	Password           string `json:"password"`
}

func (d *Dao) Login(param *Login) (*model.User, error) {
	auth := model.User{
		Username:        param.Login,
		Email:           param.Login,
		Password:        param.Password,
		CurrentSignInIp: param.ClientIP,
		Machine:         param.UserAgent,
	}
	return auth.Login(d.engine)
}

func (d *Dao) GetUser(id uint32) (model.User, error) {
	user := model.User{
		EZ_MODEL: model.EZ_MODEL{
			ID: id,
		},
	}
	return user.Get(d.engine)
}

func (d *Dao) ResendConfirmation(email string) (bool, error) {
	user := model.User{
		Email: email,
	}
	return user.ResendConfirmation(d.engine)
}

func (d *Dao) Confirmation(token string) (bool, error) {
	user := model.User{
		ConfirmationRawToken: token,
	}
	return user.Confirm(d.engine)
}

func (d *Dao) SendForgetPassword(email string) (bool, error) {
	user := model.User{
		Email: email,
	}
	return user.SendForgetPassword(d.engine)
}

func (d *Dao) ResetPassword(param *ResetPassword) (bool, error) {
	user := model.User{
		ResetPasswordRawToken: param.ResetPasswordToken,
		Password:              param.Password,
	}
	return user.ResetPassword(d.engine)
}

func (d *Dao) CreateUser(username, email, password string) (*model.User, error) {
	user := model.User{Username: username, Email: email, Password: password}
	return user.Create(d.engine)
}

func (d *Dao) UpdateProfile(param *UpdateUser) (*model.User, error) {
	user := model.User{
		EZ_MODEL: model.EZ_MODEL{
			ID: param.Id,
		},
		ApiKey:    param.ApiKey,
		ApiSecret: param.ApiSecret,
	}
	return user.UpdateProfile(d.engine)
}

func (d *Dao) LastMonthMarginFundingPaymentsSummary(userId uint32) (map[string]float32, error) {
	l := model.Ledger{ID: userId}
	return l.LastMonthMarginFundingPaymentsSummary(d.engine)
}

func (d *Dao) LastMonthMarginFundingPayments(userId uint32) (map[string][][]interface{}, error) {
	l := model.Ledger{ID: userId}
	return l.LastMonthMarginFundingPayments(d.engine)
}

func (d *Dao) GetUserRobots(userId uint32) ([]*model.UserRobot, error) {
	u := model.User{
		EZ_MODEL: model.EZ_MODEL{
			ID: userId,
		},
	}
	return u.Robots(d.engine)
}

func (d *Dao) CreateUserRobot(param *UserRobot) (*model.UserRobot, error) {
	userRobot := model.UserRobot{
		UserId:            param.UserId,
		CurrencyId:        param.CurrencyId,
		ApiKey:            param.ApiKey,
		ApiSecret:         param.ApiSecret,
		Activated:         param.Activated,
		Hidden:            param.Hidden,
		ReservedAmount:    param.ReservedAmount,
		MaxAmount:         param.MaxAmount,
		Period:            param.Period,
		FixRatePeriod:     param.FixRatePeriod,
		Intervals:         param.Intervals,
		NumLastIntervals:  param.NumLastIntervals,
		MinRate:           param.MinRate,
		FixRate:           param.FixRate,
		RandomRangeLow:    param.RandomRangeLow,
		RandomRangeHigh:   param.RandomRangeHigh,
		RangeAmount:       param.RangeAmount,
		RangeNum:          param.RangeNum,
		RangeLowRate:      param.RangeLowRate,
		RangeHighRate:     param.RangeHighRate,
		RangeRatePeriod:   param.RangeRatePeriod,
		LendingStrategyId: param.LendingStrategyId,
		RobotServerId:     param.RobotServerId,
	}

	return userRobot.Create(d.engine)
}

func (d *Dao) UpdateUserRobot(param *UserRobot) (*model.UserRobot, error) {
	global.GVA_LOG.Debug(fmt.Sprintf("robot id: %d, activate: %v", param.ID, param.Activated))
	robot, err := model.UserRobot{
		EZ_MODEL: model.EZ_MODEL{
			ID: param.ID,
		},
	}.Get(d.engine)

	if err != nil {
		return nil, err
	}

	userRobot := model.UserRobot{
		EZ_MODEL: model.EZ_MODEL{
			ID: param.ID,
		},
		UserId:                     param.UserId,
		CurrencyId:                 param.CurrencyId,
		ApiKey:                     param.ApiKey,
		ApiSecret:                  param.ApiSecret,
		Activated:                  param.Activated,
		Hidden:                     param.Hidden,
		ReservedAmount:             param.ReservedAmount,
		MaxAmount:                  param.MaxAmount,
		Period:                     param.Period,
		FixRatePeriod:              param.FixRatePeriod,
		Intervals:                  param.Intervals,
		NumLastIntervals:           param.NumLastIntervals,
		MinRate:                    param.MinRate,
		FixRate:                    param.FixRate,
		RandomRangeLow:             param.RandomRangeLow,
		RandomRangeHigh:            param.RandomRangeHigh,
		RangeAmount:                param.RangeAmount,
		RangeNum:                   param.RangeNum,
		RangeLowRate:               param.RangeLowRate,
		RangeHighRate:              param.RangeHighRate,
		RangeRatePeriod:            param.RangeRatePeriod,
		LendingStrategyId:          param.LendingStrategyId,
		RobotServerId:              robot.RobotServerId,
		ActivatedPreviouslyChanged: (robot.Activated != param.Activated),
		ApiKeyPreviouslyChanged:    (robot.ApiKey != param.ApiKey),
		ApiSecretPreviouslyChanged: (robot.ApiSecret != param.ApiSecret),
	}
	// Currency Id has changed
	if robot.CurrencyId != param.CurrencyId {
		userRobot.CurrencyIdPreviousChange = []uint32{robot.CurrencyId, param.CurrencyId}
	}

	return userRobot.Update(d.engine)
}

func (d *Dao) GetUSerRobotState(id uint32) (int32, error) {
	robot := model.UserRobot{
		EZ_MODEL: model.EZ_MODEL{
			ID: id,
		},
	}

	stateReply, err := robot.GetRobotState(d.engine)
	if err != nil {
		return 0, err
	}

	if stateReply.Code == 0 && stateReply.State == "RUNNING" {
		return 0, nil
	}
	return 1, nil
}
