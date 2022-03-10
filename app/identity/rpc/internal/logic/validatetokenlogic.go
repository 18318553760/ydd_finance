package logic

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"ydd_finance/app/identity/rpc/internal/svc"
	"ydd_finance/app/identity/rpc/pb"
	"ydd_finance/common/globalvar"
	"ydd_finance/common/perr"

	"github.com/zeromicro/go-zero/core/logx"
)
var ValidateTokenError = perr.NewErrCode(perr.TOKEN_EXPIRE_ERROR)
type ValidateTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewValidateTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ValidateTokenLogic {
	return &ValidateTokenLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  validateToken ，只很对用户服务、授权服务api开放
func (l *ValidateTokenLogic) ValidateToken(in *pb.ValidateTokenReq) (*pb.ValidateTokenResp, error) {
	// todo: add your logic here and delete this line
	userTokenKey := fmt.Sprintf(globalvar.CacheUserTokenKey, in.UserId)
	dbToken, err := l.svcCtx.RedisClient.Get(userTokenKey)
	if err != nil {
		return nil, errors.Wrapf(ValidateTokenError, "ValidateToken RedisClient Get userId:%d ,err:%v", in.UserId, err)
	}
	if dbToken != in.Token {
		return nil, errors.Wrapf(ValidateTokenError, "ValidateToken is invalid  userId:%d , token:%s , dbToken:%s", in.UserId, in.Token, dbToken)
	}
	return &pb.ValidateTokenResp{
		Ok: true,
	}, nil

}
