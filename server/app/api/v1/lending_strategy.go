package v1

import (
	"ezcoin.cc/ezcoin-go/server/app/model/common/response"
	"ezcoin.cc/ezcoin-go/server/app/service"
	"ezcoin.cc/ezcoin-go/server/global"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type LendingStrategyApi struct{}

// @Summary Get Lending Strategies
// @Tags LendingStrategy
// @Version  1.0
// @Accept   json
// @Produce  json
// @Security ApiKeyAuth
// @Success  200  {object}  response.Response{data=model.LendingStrategy}
// @Failure  404  {object}  response.Response
// @Failure  422  {object}  response.Response
// @Router /v1/lending_strategies/{id} [get]
func (l *LendingStrategyApi) GetLendingStrategies(c *gin.Context) {
	svc := service.New(c.Request.Context())
	lendingStrategies, err := svc.GetLendingStrategies()
	if err != nil {
		global.GVA_LOG.Error("取放貸策略失敗!", zap.Error(err))
		response.FailWithMessage("取放貸策略失敗", c)
		return
	}
	response.OkWithData(lendingStrategies, c)
}
