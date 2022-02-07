package validator

import (
	"ezcoin.cc/ezcoin-go/server/app/service"
	"github.com/bitfinexcom/bitfinex-api-go/pkg/models/common"
	"github.com/bitfinexcom/bitfinex-api-go/v2/rest"
	"github.com/go-playground/validator/v10"
)

func UpdateUserRobotStructLevelValidation(sl validator.StructLevel) {
	robot := sl.Current().Interface().(service.UpdateUserRobot)
	if robot.ApiKey != "" && robot.ApiSecret != "" {
		CheckBitfinexApiKey(sl, robot.ApiKey, robot.ApiSecret)
	}
}

func CreateUserRobotStructLevelValidation(sl validator.StructLevel) {
	robot := sl.Current().Interface().(service.UpdateUserRobot)
	if robot.ApiKey != "" && robot.ApiSecret != "" {
		CheckBitfinexApiKey(sl, robot.ApiKey, robot.ApiSecret)
	}
}

func CheckBitfinexApiKey(sl validator.StructLevel, apiKey, apiSecret string) {
	client := rest.NewClient().Credentials(apiKey, apiSecret)
	req, err := client.NewAuthenticatedRequest(common.PermissionRead, "permissions")
	if err != nil {
		sl.ReportError(nil, "apiKey", "", "apiKey", "")
	}
	_, err = client.Request(req)
	if err != nil {
		sl.ReportError(nil, "apiKey", "", "apiKey", "")
	}
}
