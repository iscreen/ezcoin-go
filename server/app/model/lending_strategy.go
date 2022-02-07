package model

import "gorm.io/gorm"

type LendingStrategy struct {
	EZ_MODEL
	Name       string `json:"name" gorm:"comment:放貸策略名稱"`
	Desciption string `json:"description" gorm:"type:text;comment:放貸策略說明"`
}

func (ls LendingStrategy) GetLendingStrategies(db *gorm.DB) ([]*LendingStrategy, error) {
	var lendingStrategies []*LendingStrategy
	if err := db.Find(&lendingStrategies).Error; err != nil {
		return lendingStrategies, err
	}
	return lendingStrategies, nil
}
