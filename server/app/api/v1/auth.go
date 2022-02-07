package v1

import (
	"fmt"

	"ezcoin.cc/ezcoin-go/server/app/model"
	"ezcoin.cc/ezcoin-go/server/app/model/common/response"
	"ezcoin.cc/ezcoin-go/server/app/service"
	"ezcoin.cc/ezcoin-go/server/global"
	"ezcoin.cc/ezcoin-go/server/pkg/app"
	"ezcoin.cc/ezcoin-go/server/pkg/errcode"
	"ezcoin.cc/ezcoin-go/server/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"go.uber.org/zap"
)

type AuthApi struct{}

// @Summary User registration
// @Tags Auth
// @Version 1.0
// @Accept  json
// @Produce json
// @Param param  body  service.Register  true  "Register"
// @Success  200  {object}   response.Response{data=model.User}
// @Failure  404  {object}  response.Response
// @Failure  422  {object}  response.Response
// @Router /v1/sign_up [post]
func (u *AuthApi) Register(c *gin.Context) {
	var param service.Register
	svc := service.New(c.Request.Context())
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		global.GVA_LOG.Error(fmt.Sprintf("app.BindAndValid errs: %v", errs))
		response.FailWithError(errcode.InvalidParams.WithDetails(errs.Errors()...), c)
		return
	}

	if _, err := svc.Register(&param); err != nil {
		global.GVA_LOG.Error("註冊失敗！", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("已註冊成功，請收取確認信件開通您的帳號", c)
}

// @Summary Resend confirmation instruction
// @Tags Auth
// @Version 1.0
// @Accept  json
// @Produce json
// @Param param  body  service.ResendConfirmation  true  "ResendConfirmation"
// @Success  200  {object}   response.Response
// @Failure  404  {object}  response.Response
// @Failure  422  {object}  response.Response
// @Router /v1/confirmation [post]
func (u *AuthApi) ResendConfirmation(c *gin.Context) {
	var param service.ResendConfirmation
	svc := service.New(c.Request.Context())
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		response.FailWithError(errcode.InvalidParams.WithDetails(errs.Errors()...), c)
		return
	}

	if _, err := svc.ResendConfirmation(&param); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("您將在幾分鐘後收到一封電子郵件，內有開啟帳戶的步驟說明。", c)
}

// @Summary Confirm user account via token
// @Tags Auth
// @Version 1.0
// @Accept  json
// @Produce json
// @Param confirmation_token query string true "confirmation token"
// @Success  200  {object}   response.Response
// @Failure  404  {object}  response.Response
// @Failure  422  {object}  response.Response
// @Router /v1/reset_password [patch]
func (u *AuthApi) Confirmation(c *gin.Context) {
	svc := service.New(c.Request.Context())
	token := c.Query("confirmation_token")
	if token == "" {
		response.FailWithError(errcode.InvalidParams.WithDetails("Token invalid"), c)
	}

	if _, err := svc.Confirmation(token); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("帳號已開通，請重新登入。", c)
}

// @Summary Send forget password instruction
// @Tags Auth
// @Version 1.0
// @Accept  json
// @Produce json
// @Param param  body  service.ForgetPassword  true  "ForgetPassword"
// @Success  200  {object}   response.Response
// @Failure  404  {object}  response.Response
// @Failure  422  {object}  response.Response
// @Router /v1/password [post]
func (u *AuthApi) SendForgetPassword(c *gin.Context) {
	var param service.ForgetPassword
	svc := service.New(c.Request.Context())
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		response.FailWithError(errcode.InvalidParams.WithDetails(errs.Errors()...), c)
		return
	}

	if _, err := svc.ForgetPassword(&param); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("您將在幾分鐘後收到一封電子郵件，內有重新設定密碼的步驟說明。", c)
}

// @Summary Update user password via forget password token
// @Tags Auth
// @Version 1.0
// @Accept  json
// @Produce json
// @Param user  body  service.ResetPassword  true  "ResetPassword"
// @Success  200  {object}   response.Response
// @Failure  404  {object}  response.Response
// @Failure  422  {object}  response.Response
// @Router /v1/reset_password [patch]
func (u *AuthApi) ResetPassword(c *gin.Context) {
	var param service.ResetPassword
	svc := service.New(c.Request.Context())
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		response.FailWithError(errcode.InvalidParams.WithDetails(errs.Errors()...), c)
		return
	}

	if _, err := svc.ResetPassword(&param); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("密碼已重設，請重新登入。", c)
}

// @Summary Login
// @Tags Auth
// @Version 1.0
// @Accept  json
// @Produce json
// @Param user  body  service.Login  true  "Login"
// @Success  200  {object}   response.Response{data=model.User}
// @Failure  404  {object}  response.Response
// @Failure  422  {object}  response.Response
// @Router /v1/login [post]
func (a *AuthApi) Login(c *gin.Context) {
	param := service.Login{
		ClientIP:  c.ClientIP(),
		UserAgent: c.Request.UserAgent(),
	}
	svc := service.New(c.Request.Context())
	valid, errs := app.BindAndValid(c, &param)
	global.GVA_LOG.Debug(fmt.Sprintf("login: %s \n%s", param.Login, c.Request.UserAgent()))
	if !valid {
		response.FailWithError(errcode.InvalidParams.WithDetails(errs.Errors()...), c)
		return
	}

	if user, err := svc.Login(&param); err != nil {
		global.GVA_LOG.Error("登入失敗！用戶名不存在或者密碼錯誤！", zap.Error(err))
		response.FailWithMessage("用戶名不存在或者密碼錯誤", c)
	} else {
		a.tokenNext(c, &svc, *user)
	}
}

// 登錄後簽發 JWT
func (a *AuthApi) tokenNext(c *gin.Context, svc *service.Service, user model.User) {
	j := utils.JWT{SigningKey: []byte(global.GVA_CONFIG.JWT.SigningKey)}
	claims := j.CreateClaims(service.BaseClaims{
		ID:         user.ID,
		Username:   user.Username,
		RobotLimit: user.Plan.RobotLimit,
	})

	token, err := j.CreateToken(claims)
	if err != nil {
		global.GVA_LOG.Error("獲取 token 失敗!", zap.Error(err))

		response.FailWithError(errcode.UnauthorizedTokenGenerate, c)
		// app.Response.ToErrorResponse(errcode.UnauthorizedTokenGenerate.Error())
		// response.FailWithMessage("獲取 token 失敗!", c)
		return
	}

	if !global.GVA_CONFIG.System.UseMultipoint {
		response.OkWithDetailed(service.LoginResponse{
			User:      user,
			Token:     token,
			ExpiresAt: claims.StandardClaims.ExpiresAt * 1000,
		}, "登錄成功", c)
		return
	}

	if jwtStr, err := svc.GetRedisJWT(user.Username); err == redis.Nil {
		if err := svc.SetRedisJWT(token, user.Username); err != nil {
			global.GVA_LOG.Error("設置發錄狀態失敗!", zap.Error(err))
			response.FailWithMessage("設置發錄狀態失敗", c)
			return
		}
		response.OkWithDetailed(service.LoginResponse{
			User:      user,
			Token:     token,
			ExpiresAt: claims.StandardClaims.ExpiresAt * 1000,
		}, "登錄成功", c)
	} else if err != nil {
		global.GVA_LOG.Error("設置發錄狀態失敗!", zap.Error(err))
		response.FailWithMessage("設置發錄狀態失敗", c)
	} else {
		if err := svc.JwtInBlacklist(jwtStr); err != nil {
			response.FailWithMessage("jwt 作癈失敗", c)
			return
		}
		if err := svc.SetRedisJWT(token, user.Username); err != nil {
			response.FailWithMessage("設置發錄狀態失敗", c)
			return
		}
		response.OkWithDetailed(service.LoginResponse{
			User:      user,
			Token:     token,
			ExpiresAt: claims.StandardClaims.ExpiresAt * 1000,
		}, "登錄成功", c)
	}
}
