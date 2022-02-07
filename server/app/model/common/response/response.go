package response

import (
	"net/http"

	"ezcoin.cc/ezcoin-go/server/pkg/errcode"
	"github.com/gin-gonic/gin"
)

type Response struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

type ResponseError struct {
	Code   int         `json:"code"`
	Errors interface{} `json:"errors"`
	Msg    string      `json:"msg"`
}

const (
	ERROR   = 7
	SUCCESS = 0
)

func Result(status int, code int, data interface{}, msg string, c *gin.Context) {
	if code == 0 {
		c.JSON(status, Response{
			code,
			data,
			msg,
		})
	} else {
		c.JSON(status, ResponseError{
			code,
			data,
			msg,
		})
	}
}

func Ok(c *gin.Context) {
	Result(http.StatusOK, SUCCESS, map[string]interface{}{}, "Success", c)
}

func OkWithMessage(message string, c *gin.Context) {
	Result(http.StatusOK, SUCCESS, map[string]interface{}{}, message, c)
}

func OkWithData(data interface{}, c *gin.Context) {
	Result(http.StatusOK, SUCCESS, data, "Success", c)
}

func OkWithDetailed(data interface{}, message string, c *gin.Context) {
	Result(http.StatusOK, SUCCESS, data, message, c)
}

func Fail(c *gin.Context) {
	Result(http.StatusUnprocessableEntity, ERROR, map[string]interface{}{}, "Failed", c)
}

func FailWithError(err *errcode.Error, c *gin.Context) {
	Result(http.StatusUnprocessableEntity, err.Code(), err.Details(), err.Msg(), c)
}

func FailWithMessage(message string, c *gin.Context) {
	Result(http.StatusUnprocessableEntity, ERROR, map[string]interface{}{}, message, c)
}

func FailWithDetailed(data interface{}, message string, c *gin.Context) {
	Result(http.StatusUnprocessableEntity, ERROR, data, message, c)
}
