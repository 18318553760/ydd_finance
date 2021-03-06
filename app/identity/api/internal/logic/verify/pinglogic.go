package verify

import (
	"context"

	"ydd_finance/app/identity/api/internal/svc"
	"ydd_finance/app/identity/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type PingLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPingLogic(ctx context.Context, svcCtx *svc.ServiceContext) PingLogic {
	return PingLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PingLogic) Ping() (resp *types.VerifyPingResp, err error) {
	// todo: add your logic here and delete this line
	resp = &types.VerifyPingResp{
		Ok:"pong",
	}
	return
}
