package logic

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"ydd_finance/app/identity/rpc/internal/svc"
	"ydd_finance/app/identity/rpc/pb"
	"ydd_finance/common/globalvar"
	"ydd_finance/common/perr"
)
var ErrClearTokenError = perr.NewErrMsg("退出token失败")
type ClearTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewClearTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ClearTokenLogic {
	return &ClearTokenLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  清除token，只针对用户服务开放访问
func (l *ClearTokenLogic) ClearToken(in *pb.ClearTokenReq) (*pb.ClearTokenResp, error) {
	// todo: add your logic here and delete this line
	userTokenKey := fmt.Sprintf(globalvar.CacheUserTokenKey, in.UserId)
	// redis 删除
	if _, err := l.svcCtx.RedisClient.Del(userTokenKey); err != nil {
		return nil, errors.Wrapf(ErrClearTokenError, "userId:%d,err:%v", in.UserId, err)
	}
	return &pb.ClearTokenResp{
		Ok: true,
	}, nil

}
