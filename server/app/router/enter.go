package router

import (
	v1 "ezcoin.cc/ezcoin-go/server/app/router/v1"
	v2 "ezcoin.cc/ezcoin-go/server/app/router/v2"
)

type RouterGroup struct {
	V1RouterGroup v1.RouterGroup
	V2RouterGroup v2.RouterGroup
}

var RouterGroupApp = new(RouterGroup)
