package model

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type LedgerCategory int

const (
	Unknow        LedgerCategory = iota
	Exchange      LedgerCategory = 5
	MarginPayment LedgerCategory = 28
	Settlement    LedgerCategory = 31
)

func (l LedgerCategory) String() string {
	return [...]string{"Exchange", "MarginPayment", "Settlement"}[l]
}

type Ledger struct {
	ID          uint32    `gorm:"type:bigint(20);primarykey"` // Primary Key ID
	UserId      uint32    `json:"user_id" gorm:"comment:用戶編號"`
	User        User      //`gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Category    uint      `json:"category"`
	Amount      float32   `json:"amount" gorm:"comment:金額"`
	Balance     float32   `json:"balance" gorm:"comment:金額"`
	MTS         time.Time `json:"mts" gorm:"comment:時間"`
	Currency    string    `json:"currency" gorm:"comment:幣別"`
	Description string    `json:"description" gorm:"comment:幣別"`
}

type LedgerMarginPayment struct {
	MTS      string  `json:"mts"`
	Currency string  `json:"current"`
	Amount   float32 `json:"amount" gorm:"comment:金額"`
}

//@function: LastMonthMarginFundingPayments
//@description: 取得指定用戶的最近一個月的利息收益
//@param: userId uint
//@return: map[string][][]interface{}, error
func (l Ledger) LastMonthMarginFundingPayments(db *gorm.DB) (map[string][][]interface{}, error) {
	// SELECT concat(DATE_FORMAT(mts, '%Y-%m-%d'), "_", currency), category, currency, sum(amount)
	// FROM `ledgers`
	// where user_id = 1 and mts between '2021-12-19' and '2022-01-19' and category=28
	// GROUP by concat(DATE_FORMAT(mts, '%Y-%m-%d'), "_", currency), category, currency
	var currencies []string
	var marginPayments []LedgerMarginPayment
	startDate := time.Now().AddDate(0, -1, 0)
	endDate := time.Now()
	category := fmt.Sprintf("%d", MarginPayment)

	err := db.Model(&Ledger{}).
		Select([]string{"DATE_FORMAT(mts, '%Y-%m-%d') mts", "currency", "sum(amount) amount"}).
		Where("user_id = ? AND category = ? AND mts >= ?", l.UserId, category, startDate.Format("2006-01-02")).
		Group("DATE_FORMAT(mts, '%Y-%m-%d'), currency").
		Find(&marginPayments).Error

	if len(marginPayments) == 0 && err != gorm.ErrRecordNotFound {
		return nil, errors.New("no margin payments")
	}

	result := map[string][][]interface{}{}
	currencyKeys := map[string]bool{}
	payments := map[string]LedgerMarginPayment{}
	// reorganize result and get unique currencies slice
	for _, marginPayment := range marginPayments {
		payments[marginPayment.MTS+"_"+marginPayment.Currency] = marginPayment
		if _, value := currencyKeys[marginPayment.Currency]; !value {
			currencyKeys[marginPayment.Currency] = true
		}
	}

	// convert currency map to slice
	for key := range currencyKeys {
		currencies = append(currencies, key)
	}

	// convert result sequence by date
	for _, currency := range currencies {
		currencyValue := [][]interface{}{}
		for d := startDate; !d.After(endDate); d = d.AddDate(0, 0, 1) {
			dateStr := d.Format("2006-01-02")
			if payment, ok := payments[dateStr+"_"+currency]; ok {
				currencyValue = append(currencyValue, []interface{}{currency, payment.Amount, dateStr})
			} else {
				currencyValue = append(currencyValue, []interface{}{currency, 0, dateStr})
			}
		}
		result[currency] = currencyValue
	}

	return result, nil
}

//@function: LastMonthMarginFundingPaymentsSummary
//@description: 取得指定用戶的最近一個月的利息收益摘要
//@param: userId uint
//@return: map[string][][]interface{}, error
func (l *Ledger) LastMonthMarginFundingPaymentsSummary(db *gorm.DB) (map[string]float32, error) {
	// SELECT currency, sum(amount)
	// FROM `ledgers`
	// where user_id = 1 and mts between '2021-12-19' and '2022-01-19' and category=28
	// GROUP by currency
	var marginPayments []Ledger
	startDate := time.Now().AddDate(0, -1, 0)
	category := fmt.Sprintf("%d", MarginPayment)

	db.Model(&Ledger{}).
		Select([]string{"currency", "sum(amount) amount"}).
		Where("user_id = ? AND category = ? AND mts >= ?", l.UserId, category, startDate.Format("2006-01-02")).
		Group("currency").
		Find(&marginPayments)

	if len(marginPayments) == 0 {
		return nil, errors.New("no margin payments")
	}

	result := map[string]float32{}
	for _, v := range marginPayments {
		result[v.Currency] = v.Amount
	}

	return result, nil
}
