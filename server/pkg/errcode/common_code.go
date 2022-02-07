package errcode

var (
	Success                   = NewError(0, "成功")
	ServerError               = NewError(10000000, "服務內部錯誤")
	InvalidParams             = NewError(10000001, "參數錯誤")
	NotFound                  = NewError(10000002, "找不到")
	UnauthorizedAuthNotExist  = NewError(10000003, "認證失敗，找不到對應的 AppKey 程 AppSecert")
	UnauthorizedTokenError    = NewError(10000004, "認證失敗，Token 錯誤")
	UnauthorizedTokenTimeout  = NewError(10000005, "認證失敗，Token 超時")
	UnauthorizedTokenGenerate = NewError(10000006, "認證失敗，Token 生成失敗")
	TooManyRequests           = NewError(10000007, "請求過多")
)
