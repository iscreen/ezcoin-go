package app

import (
	"ezcoin.cc/ezcoin-go/server/global"
	"ezcoin.cc/ezcoin-go/server/pkg/convert"
	"github.com/gin-gonic/gin"
)

func GetPage(c *gin.Context) int {
	page := convert.StrTo(c.Query("page")).MustInt()
	if page <= 0 {
		return 1
	}

	return page
}

func GetPageSize(c *gin.Context) int {
	pageSize := convert.StrTo(c.Query("page_size")).MustInt()
	if pageSize <= 0 {
		return global.GVA_CONFIG.System.DefaultPageSize
	}
	if pageSize > global.GVA_CONFIG.System.MaxPageSize {
		return global.GVA_CONFIG.System.MaxPageSize
	}

	return pageSize
}

func GetPageOffset(page, pageSize int) int {
	result := 0
	if page > 0 {
		result = (page - 1) * pageSize
	}

	return result
}
