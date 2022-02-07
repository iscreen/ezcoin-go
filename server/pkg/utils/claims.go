package utils

import (
	"ezcoin.cc/ezcoin-go/server/app/service"
	"ezcoin.cc/ezcoin-go/server/global"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

func GetClaims(c *gin.Context) (*service.CustomClaims, error) {
	token := c.Request.Header.Get(global.GVA_CONFIG.JWT.TokenName)
	j := NewJWT()
	claims, err := j.ParseToken(token)
	if err != nil {
		global.GVA_LOG.Error("從 Gin 的 Context 中獲取 jwt 解析失敗，請檢查請求表頭中是否存在 x-token 且 claims 是否符合結構")
	}
	return claims, err
}

func GetUserID(c *gin.Context) uint32 {
	if claims, exists := c.Get("claims"); !exists {
		if cl, err := GetClaims(c); err != nil {
			return 0
		} else {
			return cl.ID
		}
	} else {
		waitUse := claims.(*service.CustomClaims)
		return waitUse.ID
	}
}

// 從 Gin 的 Context 中取得 jwt 解析出來的用戶 UUID
func GetUserUUID(c *gin.Context) uuid.UUID {
	if claims, exists := c.Get("claims"); !exists {
		if cl, err := GetClaims(c); err != nil {
			return uuid.UUID{}
		} else {
			return cl.UUID
		}
	} else {
		waitUse := claims.(*service.CustomClaims)
		return waitUse.UUID
	}
}

// 從 Gin 的 Context 中取得 jwt 解析出來的用戶角色 ID
func GetUserAuthorityId(c *gin.Context) string {
	if claims, exists := c.Get("claims"); !exists {
		if cl, err := GetClaims(c); err != nil {
			return ""
		} else {
			return cl.AuthorityId
		}
	} else {
		waitUse := claims.(*service.CustomClaims)
		return waitUse.AuthorityId
	}
}

// 從 Gin 的 Context 中取得 jwt 解析出來的用戶
func GetUserInfo(c *gin.Context) *service.CustomClaims {
	if claims, exists := c.Get("claims"); !exists {
		if cl, err := GetClaims(c); err != nil {
			return nil
		} else {
			return cl
		}
	} else {
		waitUse := claims.(*service.CustomClaims)
		return waitUse
	}
}
