package validator

import (
	"ezcoin.cc/ezcoin-go/server/app/service"
	"ezcoin.cc/ezcoin-go/server/global"
	"github.com/go-playground/validator/v10"
)

func UserStructLevelValidation(sl validator.StructLevel) {
	global.GVA_LOG.Debug("entry UserStructLevelValidation")
	user := sl.Current().Interface().(service.UpdateUser)

	if user.ApiKey != "" && user.ApiSecret != "" {
		CheckBitfinexApiKey(sl, user.ApiKey, user.ApiSecret)
	}
}
