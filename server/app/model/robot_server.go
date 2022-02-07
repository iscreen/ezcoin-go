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

type RobotServer struct {
	EZ_MODEL
	Name string `json:"name" gorm:"comment:放貸主機名稱"`
	IP   string `json:"ip" gorm:"comment:放貸主機IP"`
}

func (rs RobotServer) List(db *gorm.DB, param request.TableQuery) (*response.PageResult, error) {
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

func (rs RobotServer) Create(db *gorm.DB) (*RobotServer, error) {
	if err := db.Create(&rs).Error; err != nil {
		return nil, err
	}
	return &rs, nil
}

func (rs RobotServer) Update(db *gorm.DB) (*RobotServer, error) {
	if err := db.Where("id = ?", rs.ID).Updates(&rs).Error; err != nil {
		return nil, err
	}
	return &rs, nil
}

func (rs RobotServer) Get(db *gorm.DB) (RobotServer, error) {
	var robotServer RobotServer
	err := db.Where("id = ?", rs.ID).First(&robotServer).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return robotServer, err
	}
	return robotServer, nil
}

func (rs RobotServer) Delete(db *gorm.DB) error {
	if err := db.Model(&rs).Update("deleted_at", time.Now()).Error; err != nil {
		return err
	}

	return nil
}

func (rs RobotServer) GetAny(db *gorm.DB) (RobotServer, error) {
	var robotServer RobotServer
	err := db.Where("id = ?", rs.ID).Order("RAND()").First(&robotServer).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return robotServer, err
	}
	return robotServer, nil
}
