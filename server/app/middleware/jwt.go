package middleware

import (
	"strconv"
	"strings"
	"time"

	"ezcoin.cc/ezcoin-go/server/app/model/common/response"
	"ezcoin.cc/ezcoin-go/server/app/service"
	"ezcoin.cc/ezcoin-go/server/global"
	"ezcoin.cc/ezcoin-go/server/pkg/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get(global.GVA_CONFIG.JWT.TokenName)
		tokens := strings.Split(token, " ")
		if len(tokens) != 2 {
			response.FailWithDetailed(gin.H{"reload": true}, "未登入或非正常訪問", c)
			c.Abort()
			return
		}
		token = tokens[len(tokens)-1]
		// global.GVA_LOG.Debug(token)
		if token == "" {
			response.FailWithDetailed(gin.H{"reload": true}, "未登入或非正常訪問", c)
			c.Abort()
			return
		}

		svc := service.New(c.Request.Context())
		if svc.IsBlacklist(token) {
			response.FailWithDetailed(gin.H{"reload": true}, "您的帳戶異地登入或Token失效", c)
			c.Abort()
			return
		}
		j := utils.NewJWT()
		claims, err := j.ParseToken(token)
		if err != nil {
			if err == utils.TokenExpired {
				response.FailWithDetailed(gin.H{"reload": true}, "授權取過期", c)
				c.Abort()
				return
			}
			response.FailWithDetailed(gin.H{"reload": true}, err.Error(), c)
			c.Abort()
			return
		}
		if claims.ExpiresAt-time.Now().Unix() < claims.BufferTime {
			claims.ExpiresAt = time.Now().Unix() + global.GVA_CONFIG.JWT.ExpiresTime
			newToken, _ := j.CreateTokenByOldToken(token, *claims)
			newClaims, _ := j.ParseToken(newToken)
			c.Header("new-token", newToken)
			c.Header("new-expires-at", strconv.FormatInt(newClaims.ExpiresAt, 10))
			if global.GVA_CONFIG.System.UseMultipoint {
				RedisJwtToken, err := svc.GetRedisJWT(newClaims.Username)
				if err != nil {
					global.GVA_LOG.Error("get redis jwt failed", zap.Error(err))
				} else { // 當之前的取成功時才進入黑名單
					_ = svc.JwtInBlacklist(RedisJwtToken)
				}
				// 記錄當前的活躍狀態
				_ = svc.SetRedisJWT(newToken, newClaims.Username)
			}
		}
		c.Set("claims", claims)
		c.Next()
	}
}
