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

type Plan struct {
	EZ_MODEL
	Name       string `json:"name" gorm:"comment:訂閱計畫名稱"`
	Desciption string `json:"description" gorm:"type:text,comment:訂閱計畫說明"`
	ExternalId string `json:"external_id" gorm:"comment:金流計畫編號"`
	Period     uint   `json:"period" gorm:"comment:訂閱計畫週期"`
	RobotLimit uint   `json:"robot_limit" gorm:"comment:Robot 限制"`
	Default    bool   `json:"default" gorm:"default:false"`
	Trial      bool   `json:"trial" gorm:"default:false"`
}

func (p Plan) List(db *gorm.DB, param request.TableQuery) (*response.PageResult, error) {
	var list []*Plan
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

func (p Plan) Create(db *gorm.DB) (*Plan, error) {
	if err := db.Create(&p).Error; err != nil {
		return nil, err
	}
	return &p, nil
}

func (p Plan) Update(db *gorm.DB) (*Plan, error) {
	if err := db.Where("id = ?", p.ID).Updates(&p).Error; err != nil {
		return nil, err
	}
	return &p, nil
}

func (p Plan) Get(db *gorm.DB) (Plan, error) {
	var plan Plan
	err := db.Where("id = ? AND deleted_at IS NOT NULL", p.ID).First(&plan).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return plan, err
	}
	return plan, nil
}

func (p Plan) Delete(db *gorm.DB) error {
	if err := db.Model(&p).Update("deleted_at", time.Now()).Error; err != nil {
		return err
	}

	return nil
}
