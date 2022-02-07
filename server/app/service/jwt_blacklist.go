package service

import (
	"context"
	"time"

	"ezcoin.cc/ezcoin-go/server/app/model"
	"ezcoin.cc/ezcoin-go/server/global"
	"github.com/dgrijalva/jwt-go"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
)

type CustomClaims struct {
	BaseClaims
	BufferTime int64
	jwt.StandardClaims
}

type BaseClaims struct {
	UUID        uuid.UUID
	ID          uint32 `json:"id"`
	Username    string `json:"username"`
	Nickname    string `json:"nickname"`
	AuthorityId string `json:"authorityId"`
	RobotLimit  uint   `json:"robotLimit"`
}

// type JwtService struct{}

//@function: JwtInBlacklist
//@description: 判断JWT是否在黑名单内部
//@param: jwt string
//@return: bool
func (svc *Service) JwtInBlacklist(jwtStr string) (err error) {
	// var blackJWT model.JwtBlacklist
	// blackJWT.Jwt = jwtStr
	// err = global.GVA_DB.Engine().Create(&blackJWT).Error
	// if err != nil {
	// 	return
	// }
	// global.BlackCache.SetDefault(blackJWT.Jwt, struct{}{})
	return
}

//@function: IsBlacklist
//@description: 判斷 JWT 是否在黑名單內
//@param: jwt string
//@return: bool
func (svc *Service) IsBlacklist(jwt string) bool {
	_, ok := global.BlackCache.Get(jwt)
	return ok
	// err := global.GVA_DB.Where("jwt = ?", jwt).First(&system.JwtBlacklist{}).Error
	// isNotFound := errors.Is(err, gorm.ErrRecordNotFound)
	// return !isNotFound
}

//@function: GetRedisJWT
//@description: 从redis取jwt
//@param: userName string
//@return: redisJWT string, err error

func (svc *Service) GetRedisJWT(userName string) (redisJWT string, err error) {
	redisJWT, err = global.GVA_REDIS.Get(context.Background(), userName).Result()
	return redisJWT, err
}

//@function: SetRedisJWT
//@description: jwt存入redis并设置过期时间
//@param: jwt string, userName string
//@return: err error
func (svc *Service) SetRedisJWT(jwt string, userName string) (err error) {
	// 此处过期时间等于jwt过期时间
	timer := time.Duration(global.GVA_CONFIG.JWT.ExpiresTime) * time.Second
	err = global.GVA_REDIS.Set(context.Background(), userName, jwt, timer).Err()
	return err
}

// jwt 黑名單，加入 BlackCache 中
func LoadAll() {
	var data []string
	err := global.GVA_DB.Model(&model.JwtBlacklist{}).Select("jwt").Find(&data).Error
	if err != nil {
		global.GVA_LOG.Error("截入DB jwt 黑名單失敗！", zap.Error(err))
		return
	}
	for i := 0; i < len(data); i++ {
		global.BlackCache.SetDefault(data[i], struct{}{})
	}
}
