package verify

import (
	"context"
	"fmt"
	"ydd_finance/app/identity/rpc/identity"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
	"ydd_finance/app/identity/api/internal/svc"
	"ydd_finance/app/identity/api/internal/types"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) LoginLogic {
	return LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req types.VerifyLoginReq) (resp *types.VerifyLoginResp, err error) {
	// todo: add your logic here and delete this line
	loginResp, err := l.svcCtx.IdentityRpc.Login(l.ctx, &identity.ValidateLoginReq{
		Username:req.Username,
		Pass:req.Pass,
	})
	fmt.Println(err)
	if err != nil {
		return nil, err
	}
	resp = &types.VerifyLoginResp{}
	_ = copier.Copy(resp, loginResp)
	return resp, nil

}
