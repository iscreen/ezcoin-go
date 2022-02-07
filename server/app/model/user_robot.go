package model

import (
	"context"
	"fmt"
	"math"
	"time"

	"ezcoin.cc/ezcoin-go/server/app/model/common/request"
	"ezcoin.cc/ezcoin-go/server/app/model/common/response"
	pb "ezcoin.cc/ezcoin-go/server/app/protos"
	"ezcoin.cc/ezcoin-go/server/global"
	"ezcoin.cc/ezcoin-go/server/pkg/app"
	"github.com/shopspring/decimal"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

type UserRobot struct {
	EZ_MODEL
	UserId            uint32          `json:"userId" validate:"required" gorm:"comment:用戶編號"`
	User              User            `json:"user" gorm:"foreignKey:UserId;references:ID;"`
	CurrencyId        uint32          `json:"currencyId" validate:"required" gorm:"comment:貨幣編號"`
	Currency          Currency        `json:"currency"`
	ApiKey            string          `json:"apiKey" gorm:"comment:Bitfinex API Key"`
	ApiSecret         string          `json:"apiSecret" gorm:"comment:Bitfinex API Secret"`
	Activated         bool            `json:"activated" gorm:"type:boolean;comment:啟用;default:false"`
	Hidden            bool            `json:"hidden" gorm:"type:boolean;comment:隱藏放款;default:false"`
	ReservedAmount    decimal.Decimal `json:"reservedAmount" gorm:"type:decimal(15,4);comment:保留放款金額;default:0"`
	MaxAmount         decimal.Decimal `json:"maxAmount" gorm:"type:decimal(10,4);comment:最大放款金額;default:500.00"`
	Period            uint            `json:"period" gorm:"comment:放款期間;default:3"`
	FixRatePeriod     uint            `json:"fixRatePeriod" gorm:"comment:固定利率放款期間;default:10"`
	Intervals         uint            `json:"intervals" gorm:"comment:放貸交易週期;default:15"`
	NumLastIntervals  uint            `json:"numLast_intervals" gorm:"comment:最近幾期放貸交易;default:3"`
	MinRate           decimal.Decimal `json:"minRate" gorm:"type:decimal(15,12);comment:最小利率;default:0.0005"`
	FixRate           decimal.Decimal `json:"fixRate" gorm:"type:decimal(15,12);comment:固定利率;default:0.0005"`
	RandomRangeLow    float32         `json:"randomRangeLow" gorm:"default:-80"`
	RandomRangeHigh   float32         `json:"randomRangeHigh" gorm:"default:-30"`
	RangeAmount       decimal.Decimal `json:"rangeAmount" gorm:"type:decimal(15,4);comment:區間金額;default:20000"`
	RangeNum          uint            `json:"rangeNum" gorm:"type:decimal(15,12);comment:區間放貸拆分數;default:10"`
	RangeLowRate      decimal.Decimal `json:"rangeLowRate" gorm:"type:decimal(15,12);comment:區間放貸最低利率;default:0.0008"`
	RangeHighRate     decimal.Decimal `json:"rangeHighRate" gorm:"type:decimal(15,12);comment:區間放貸最高利率;default:0.0009"`
	RangeRatePeriod   uint            `json:"rangeRatePeriod" gorm:"comment:區間放貸放款期間;default:0"`
	LendingStrategyId uint32          `json:"lendingStrategyId" validate:"required" gorm:"comment:策略編號"`
	LendingStrategy   LendingStrategy `json:"lendingStrategy" gorm:"foreignKey:LendingStrategyId;references:ID;"`
	RobotServerId     uint32          `json:"robotServerId" gorm:"comment:放貸主機編號"`
	RobotServer       RobotServer     `json:"robotServer"`

	CurrencyIdPreviousChange   []uint32 `json:"-" gorm:"-"`
	ActivatedPreviouslyChanged bool     `json:"-" gorm:"-"`
	ApiKeyPreviouslyChanged    bool     `json:"-" gorm:"-"`
	ApiSecretPreviouslyChanged bool     `json:"-" gorm:"-"`
	RobotServerIP              string   `json:"-" gorm:"-"`
	Username                   string   `json:"-" gorm:"-"`
	CurrencySymbolName         string   `json:"-" gorm:"-"`
}

func (ur UserRobot) List(db *gorm.DB, param request.TableQuery) (*response.PageResult, error) {
	var list []*RobotServer
	var result response.PageResult
	var count int64
	var query *gorm.DB

	query = db.Model(User{})

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

func (ur UserRobot) Create(db *gorm.DB) (*UserRobot, error) {
	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&ur).Error; err != nil {
			tx.Rollback()
			return err
		}
		user := User{
			EZ_MODEL: EZ_MODEL{
				ID: ur.UserId,
			},
		}
		if err := tx.Model(&user).Where("id = ?", ur.UserId).Update("user_robots_count", gorm.Expr("user_robots_count + ?", 1)).Error; err != nil {
			tx.Rollback()
			return err
		}

		return nil
	})
	return &ur, err
}

func (ur UserRobot) Update(db *gorm.DB) (*UserRobot, error) {
	global.GVA_LOG.Debug(fmt.Sprintf("activated: %v", ur.Activated))
	if err := db.Model(&ur).Where("id = ?", ur.ID).Updates(map[string]interface{}{
		"currency_id":         ur.CurrencyId,
		"api_key":             ur.ApiKey,
		"api_secret":          ur.ApiSecret,
		"hidden":              ur.Hidden,
		"activated":           ur.Activated,
		"reserved_amount":     ur.ReservedAmount,
		"period":              ur.Period,
		"fix_rate_period":     ur.FixRatePeriod,
		"intervals":           ur.Intervals,
		"num_last_intervals":  ur.NumLastIntervals,
		"min_rate":            ur.MinRate,
		"fix_rate":            ur.FixRate,
		"random_range_low":    ur.RandomRangeLow,
		"random_range_high":   ur.RandomRangeHigh,
		"range_amount":        ur.RangeAmount,
		"range_num":           ur.RangeNum,
		"range_low_rate":      ur.RangeLowRate,
		"range_high_rate":     ur.RangeHighRate,
		"range_rate_period":   ur.RangeRatePeriod,
		"lending_strategy_id": ur.LendingStrategyId,
	}).Error; err != nil {
		return nil, err
	}

	if err := db.Where("id = ?", ur.ID).Preload("User").Preload("Currency").Preload("LendingStrategy").Find(&ur).Error; err != nil {
		return nil, err
	}

	return &ur, nil
}

func (ur UserRobot) Get(db *gorm.DB) (UserRobot, error) {
	var robot UserRobot
	err := db.Where("id = ?", ur.ID).Preload("Currency").Preload("LendingStrategy").Preload("User").Preload("RobotServer").First(&robot).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return robot, err
	}
	return robot, nil
}

func (ur UserRobot) Delete(db *gorm.DB) error {
	if err := db.Delete(&ur).Error; err != nil {
		return err
	}
	return nil
}

func (ur UserRobot) GetRobotState(db *gorm.DB) (*pb.StatusReply, error) {
	if err := db.Where("id = ?", ur.ID).Preload("Currency").Preload("User").Preload("RobotServer").Find(&ur).Error; err != nil {
		return nil, err
	}
	ur.RobotServerIP = ur.RobotServer.IP
	return ur.robotState()
}

func (ur UserRobot) createRobot(tx *gorm.DB) error {
	user, _ := User{
		EZ_MODEL: EZ_MODEL{
			ID: ur.UserId,
		},
	}.Get(global.GVA_DB)

	currency, _ := Currency{
		EZ_MODEL: EZ_MODEL{
			ID: ur.CurrencyId,
		},
	}.Get(global.GVA_DB)
	stub, err := ur.grpcStub(global.GVA_DB)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := stub.CreateRobot(ctx,
		&pb.RobotRequest{
			Name:     user.Username,
			Currency: currency.SymbolName,
		},
	)

	if err != nil {
		global.GVA_LOG.Error(fmt.Sprintf("could not create robot: %v", err))
		return err
	}
	global.GVA_LOG.Info(fmt.Sprintf("CreateRobot: code: %d, message: %s", r.Code, r.Message))
	return nil
}

func (ur UserRobot) robotState() (*pb.StatusReply, error) {
	user, _ := User{
		EZ_MODEL: EZ_MODEL{
			ID: ur.UserId,
		},
	}.Get(global.GVA_DB)

	currency, _ := Currency{
		EZ_MODEL: EZ_MODEL{
			ID: ur.CurrencyId,
		},
	}.Get(global.GVA_DB)

	stub, err := ur.grpcStub(global.GVA_DB)
	if err != nil {
		return nil, err
	}
	global.GVA_LOG.Debug(fmt.Sprintf("name: %s, currency: %s", user.Username, currency.SymbolName))
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := stub.RobotStatus(ctx,
		&pb.RobotRequest{
			Name:     user.Username,
			Currency: currency.SymbolName,
		},
	)

	if err != nil {
		global.GVA_LOG.Error(fmt.Sprintf("could not check robot status: %v", err))
		return nil, err
	}

	global.GVA_LOG.Info(fmt.Sprintf("RobotState: code: %d, state: %s, message: %s", r.Code, r.State, r.Message))
	return r, nil
}

func (ur UserRobot) startRobot() (*pb.StatusReply, error) {

	user, _ := User{
		EZ_MODEL: EZ_MODEL{
			ID: ur.UserId,
		},
	}.Get(global.GVA_DB)

	currency, _ := Currency{
		EZ_MODEL: EZ_MODEL{
			ID: ur.CurrencyId,
		},
	}.Get(global.GVA_DB)

	stub, err := ur.grpcStub(global.GVA_DB)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := stub.StartRobot(ctx,
		&pb.RobotRequest{
			Name:     user.Username,
			Currency: currency.SymbolName,
		},
	)

	if err != nil {
		global.GVA_LOG.Error(fmt.Sprintf("startRobot: %v", err))
		return nil, err
	}
	global.GVA_LOG.Info(fmt.Sprintf("startRobot: code: %d, message: %s", r.Code, r.Message))
	return r, nil
}

func (ur UserRobot) stopRobot() (*pb.StatusReply, error) {
	user, _ := User{
		EZ_MODEL: EZ_MODEL{
			ID: ur.UserId,
		},
	}.Get(global.GVA_DB)

	currency, _ := Currency{
		EZ_MODEL: EZ_MODEL{
			ID: ur.CurrencyId,
		},
	}.Get(global.GVA_DB)
	stub, err := ur.grpcStub(global.GVA_DB)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	global.GVA_LOG.Debug(fmt.Sprintf("name: %s, currency: %s", user.Username, currency.SymbolName))
	r, err := stub.StopRobot(ctx,
		&pb.RobotRequest{
			Name:     user.Username,
			Currency: currency.SymbolName,
		},
	)

	if err != nil {
		global.GVA_LOG.Error(fmt.Sprintf("stopRobot: %v", err))
		return nil, err
	}
	global.GVA_LOG.Info(fmt.Sprintf("stopRobot: code: %d, message: %s", r.Code, r.Message))
	return r, nil
}

func (ur UserRobot) restartRobot() (*pb.StatusReply, error) {
	user, _ := User{
		EZ_MODEL: EZ_MODEL{
			ID: ur.UserId,
		},
	}.Get(global.GVA_DB)

	currency, _ := Currency{
		EZ_MODEL: EZ_MODEL{
			ID: ur.CurrencyId,
		},
	}.Get(global.GVA_DB)

	stub, err := ur.grpcStub(global.GVA_DB)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	global.GVA_LOG.Debug(fmt.Sprintf("name: %s, currency: %s", user.Username, currency.SymbolName))
	r, err := stub.RestartRobot(ctx,
		&pb.RobotRequest{
			Name:     user.Username,
			Currency: currency.SymbolName,
		},
	)

	if err != nil {
		global.GVA_LOG.Error(fmt.Sprintf("restartRobot: %v", err))
		return nil, err
	}
	global.GVA_LOG.Info(fmt.Sprintf("restartRobot: code: %d, message: %s", r.Code, r.Message))
	return r, nil
}

func (ur UserRobot) migrateRobot(currencies []uint32) error {
	user, err := User{
		EZ_MODEL: EZ_MODEL{
			ID: ur.UserId,
		},
	}.Get(global.GVA_DB)
	if err != nil {
		return err
	}

	fromCurrency, err := Currency{
		EZ_MODEL: EZ_MODEL{
			ID: currencies[0],
		},
	}.Get(global.GVA_DB)

	if err != nil {
		return err
	}

	toCurrency, err := Currency{
		EZ_MODEL: EZ_MODEL{
			ID: currencies[1],
		},
	}.Get(global.GVA_DB)

	if err != nil {
		return err
	}
	stub, err := ur.grpcStub(global.GVA_DB)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := stub.MigrateRobot(ctx,
		&pb.RobotMigrateRequest{
			Name:         user.Username,
			FromCurrency: fromCurrency.SymbolName,
			ToCurrency:   toCurrency.SymbolName,
		},
	)

	if err != nil {
		global.GVA_LOG.Error(fmt.Sprintf("migrateRobot: %v", err))
		return err
	}
	global.GVA_LOG.Info(fmt.Sprintf("migrateRobot: code: %d, message: %s", r.Code, r.Message))

	if ur.Activated {
		// restart robot
		go ur.startRobot()
	} else {
		// stop robot
		go ur.stopRobot()
	}
	return nil
}

func (ur UserRobot) grpcStub(db *gorm.DB) (pb.EZCoinRobotClient, error) {
	robotServer, err := RobotServer{
		EZ_MODEL: EZ_MODEL{
			ID: ur.RobotServerId,
		},
	}.Get(db)

	if err != nil {
		return nil, err
	}
	global.GVA_LOG.Debug(fmt.Sprintf("&&&&&robot server ip: %s", robotServer.IP))
	conn, err := grpc.Dial(robotServer.IP, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	// defer conn.Close()
	return pb.NewEZCoinRobotClient(conn), nil
}

func (ur *UserRobot) BeforeCreate(tx *gorm.DB) (err error) {
	var robotServer RobotServer

	server, err := robotServer.GetAny(tx)
	if err != nil {
		return err
	}
	ur.RobotServerId = server.ID
	return
}

func (ur *UserRobot) AfterCreate(tx *gorm.DB) (err error) {
	// create supervisor robot
	go ur.createRobot(tx)
	return
}

func (ur *UserRobot) AfterFind(tx *gorm.DB) (err error) {
	ur.MinRate = ur.MinRate.Mul(decimal.NewFromFloat32(100))
	ur.FixRate = ur.FixRate.Mul(decimal.NewFromFloat32(100))
	ur.RangeLowRate = ur.RangeLowRate.Mul(decimal.NewFromFloat32(100))
	ur.RangeHighRate = ur.RangeHighRate.Mul(decimal.NewFromFloat32(100))
	return
}

func (ur *UserRobot) BeforeSave(tx *gorm.DB) (err error) {
	ur.MinRate = ur.MinRate.Div(decimal.NewFromFloat32(100))
	ur.FixRate = ur.FixRate.Div(decimal.NewFromFloat32(100))
	ur.RangeLowRate = ur.RangeLowRate.Div(decimal.NewFromFloat32(100))
	ur.RangeHighRate = ur.RangeHighRate.Div(decimal.NewFromFloat32(100))
	return
}

func (ur *UserRobot) BeforeUpdate(tx *gorm.DB) (err error) {
	return
}

func (ur *UserRobot) AfterUpdate(tx *gorm.DB) (err error) {
	// create supervisor robot
	if len(ur.CurrencyIdPreviousChange) != 0 {
		go ur.migrateRobot(ur.CurrencyIdPreviousChange)
	} else if ur.ActivatedPreviouslyChanged {
		if ur.Activated {
			go ur.restartRobot()
		} else {
			go ur.stopRobot()
		}
	} else {
		if !ur.ApiKeyPreviouslyChanged && !ur.ApiSecretPreviouslyChanged {
			return
		}
		if ur.Activated {
			go ur.restartRobot()
		} else {
			go ur.stopRobot()
		}
	}
	return
}

func (ur *UserRobot) AfterDelete(tx *gorm.DB) (err error) {
	//TODO: call rpc to delete supervisor service
	return
}
