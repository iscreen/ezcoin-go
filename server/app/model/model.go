package model

import (
	"time"

	"gorm.io/gorm"
)

type EZ_MODEL struct {
	// gorm.Model
	ID        uint32         `json:"id" gorm:"type:bigint(20);primarykey"` // Primary Key ID
	CreatedAt time.Time      `json:"createdAt" gorm:"autoCreateTime"`      // Create Time
	UpdatedAt time.Time      `json:"updatedAt" gorm:"autoUpdateTime"`      // Update Time
	DeletedAt gorm.DeletedAt `json:"deletedAt" gorm:"index"`               // Delete Time
}

func (model *EZ_MODEL) BeforeUpdate(tx *gorm.DB) (err error) {
	model.UpdatedAt = time.Now()
	return
}
