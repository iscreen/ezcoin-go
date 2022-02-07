package service

import (
	"context"

	"ezcoin.cc/ezcoin-go/server/app/dao"
	"ezcoin.cc/ezcoin-go/server/global"
)

type Service struct {
	ctx context.Context
	dao *dao.Dao
}

func New(ctx context.Context) Service {
	svc := Service{ctx: ctx}
	svc.dao = dao.New(global.GVA_DB)
	return svc
}

func (svc *Service) Context() *context.Context {
	return &svc.ctx
}
