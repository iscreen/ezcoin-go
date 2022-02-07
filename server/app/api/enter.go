package api

import (
	v1 "ezcoin.cc/ezcoin-go/server/app/api/v1"
)

type ApiGroup struct {
	ApiV1Group v1.ApiV1Group
}

var ApiGroupApp = new(ApiGroup)
