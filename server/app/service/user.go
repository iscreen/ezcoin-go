package service

import (
	"ezcoin.cc/ezcoin-go/server/app/dao"
	"ezcoin.cc/ezcoin-go/server/app/model"
)

type Login struct {
	Login     string `form:"login" json:"login" binding:"required" validate:"required"`
	Password  string `form:"password" json:"password" binding:"required" validate:"required"`
	ClientIP  string `json:"-"`
	UserAgent string `json:"-"`
}

type LoginResponse struct {
	User      model.User `json:"user"`
	Token     string     `json:"token"`
	ExpiresAt int64      `json:"expiresAt"`
}

type Register struct {
	Email                string `form:"email" json:"email" binding:"required" validate:"required,email"`
	Username             string `form:"username" json:"username" binding:"required" validate:"required"`
	Password             string `form:"password" json:"password" binding:"required" validate:"required"`
	PasswordConfirmation string `form:"passwordConfirmation" json:"passwordConfirmation" binding:"required,eqfield=Password"`
}

type ResendConfirmation struct {
	Email string `form:"email" json:"email" binding:"required" validate:"required,email"`
}

type Confirmation struct {
	ConfirmationToken string `form:"confirmationToken" json:"confirmationToken" binding:"required" validate:"required"`
	Password          string `form:"password" json:"password" binding:"required" validate:"required"`
}

type ForgetPassword struct {
	Email string `form:"email" json:"email" binding:"required" validate:"required,email"`
}

type ResetPassword struct {
	ResetPasswordToken string `form:"resetPasswordToken" json:"resetPasswordToken" binding:"required" validate:"required"`
	Password           string `form:"password" json:"password" binding:"required" validate:"required"`
}

type UpdateUser struct {
	Id        uint32 `json:"id" binding:"required"`
	ApiKey    string `json:"apiKey" binding:"required"`
	ApiSecret string `json:"apiSecret" binding:"required"`
}

//@function: Register
//@description: User Register
//@param: u model.User
//@return: userInter model.User, err error
func (svc *Service) Register(param *Register) (user *model.User, err error) {
	user, err = svc.dao.CreateUser(param.Username, param.Email, param.Password)
	if err != nil {
		return nil, err
	}
	return user, err
}

//@function: Login
//@description: User Login
//@param: u model.User
//@return: userInter model.User, err error
func (svc *Service) Login(param *Login) (user *model.User, err error) {
	user, err = svc.dao.Login(&dao.Login{
		Login:     param.Login,
		Password:  param.Password,
		ClientIP:  param.ClientIP,
		UserAgent: param.UserAgent,
	})
	if err != nil {
		return nil, err
	}
	return user, err
}

//@function: ResendConfirmation
//@description: User send confirmation instruction
//@param: param *ResendConfirmation
//@return: bool, error
func (svc *Service) ResendConfirmation(param *ResendConfirmation) (bool, error) {
	return svc.dao.ResendConfirmation(param.Email)
}

//@function: Confirmation
//@description: User confirmation account via token
//@param: param *Confirmation
//@return: bool, error
func (svc *Service) Confirmation(token string) (bool, error) {
	return svc.dao.Confirmation(token)
}

//@function: ForgetPassword
//@description: User send forget password instruction
//@param: param *ForgetPassword
//@return: bool, error
func (svc *Service) ForgetPassword(param *ForgetPassword) (bool, error) {
	return svc.dao.SendForgetPassword(param.Email)
}

//@function: ResetPassword
//@description: User reset password via token
//@param: param *ResetPassword
//@return: bool, error
func (svc *Service) ResetPassword(param *ResetPassword) (bool, error) {
	return svc.dao.ResetPassword(
		&dao.ResetPassword{
			ResetPasswordToken: param.ResetPasswordToken,
			Password:           param.Password,
		},
	)
}

//@function: GetUser
//@description: Get user
//@param: userId uint
//@return: userInter *model.User, err error
func (svc *Service) GetUser(userId uint32) (user model.User, err error) {
	return svc.dao.GetUser(userId)
}

//@function: UpdateUser
//@description: Update user profile
//@param: userId uint
//@return: model.User, error
func (svc *Service) UpdateUser(param *UpdateUser) (*model.User, error) {
	return svc.dao.UpdateProfile(&dao.UpdateUser{
		Id:        param.Id,
		ApiKey:    param.ApiKey,
		ApiSecret: param.ApiSecret,
	})
}

func (svc *Service) LastMonthMarginFundingPaymentsSummary(userId uint32) (map[string]float32, error) {
	return svc.dao.LastMonthMarginFundingPaymentsSummary(userId)
}

func (svc *Service) LastMonthMarginFundingPayments(userId uint32) (map[string][][]interface{}, error) {
	return svc.dao.LastMonthMarginFundingPayments(userId)
}

//@function: GetUserRobots
//@description: Get user's robots
//@param: userId
//@return: []*model.UserRobot, err error
func (svc *Service) GetUserRobots(userId uint32) ([]*model.UserRobot, error) {
	return svc.dao.GetUserRobots(userId)
}

//@function: CreateUserRobot
//@description: Create user's robot
//@param: CreateUserRobot
//@return: *model.UserRobot, error
func (svc *Service) CreateUserRobot(param *CreateUserRobot) (*model.UserRobot, error) {

	robot, err := svc.dao.CreateUserRobot(&dao.UserRobot{
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
		RandomRangeLow:    param.RandomRange.RandomRangeLow,
		RandomRangeHigh:   param.RandomRange.RandomRangeHigh,
		RangeAmount:       param.RangeAmount,
		RangeNum:          param.RangeNum,
		RangeLowRate:      param.RangeRate.RangeLowRate,
		RangeHighRate:     param.RangeRate.RangeHighRate,
		RangeRatePeriod:   param.RangeRatePeriod,
		LendingStrategyId: param.LendingStrategyId,
	})

	return robot, err
}

//@function: UpdateUserRobot
//@description: Update user's robot
//@param: UpdateUserRobot
//@return: *model.UserRobot, error
func (svc *Service) UpdateUserRobot(param *UpdateUserRobot) (*model.UserRobot, error) {
	robot, err := svc.dao.UpdateUserRobot(&dao.UserRobot{
		ID:                param.ID,
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
		RandomRangeLow:    param.RandomRange.RandomRangeLow,
		RandomRangeHigh:   param.RandomRange.RandomRangeHigh,
		RangeAmount:       param.RangeAmount,
		RangeNum:          param.RangeNum,
		RangeLowRate:      param.RangeRate.RangeLowRate,
		RangeHighRate:     param.RangeRate.RangeHighRate,
		RangeRatePeriod:   param.RangeRatePeriod,
		LendingStrategyId: param.LendingStrategyId,
	})

	return robot, err
}

//@function: GetUserRobotState
//@description: Get robot state
//@param: id uint32
//@return: int32, error
func (svc *Service) GetUserRobotState(id uint32) (int32, error) {
	return svc.dao.GetUSerRobotState(id)
}
