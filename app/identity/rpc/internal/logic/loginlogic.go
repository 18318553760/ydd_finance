package logic

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"ydd_finance/common/perr"
	"ydd_finance/app/identity/model"
	"ydd_finance/app/identity/rpc/internal/svc"
	"ydd_finance/app/identity/rpc/pb"
	"github.com/zeromicro/go-zero/core/logx"
)

var ErrUserNoExistsError = perr.NewErrMsg("账号不存在")
var ErrUsernamePwdError = perr.NewErrMsg("账号或密码不正确")
type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  当前这步需要分离可以另外起服务，为了测试暂时先放这里
func (l *LoginLogic) Login(in *pb.ValidateLoginReq) (*pb.ValidateLoginResp, error) {
	// todo: add your logic here and delete this line
	// 数据库操作
	var userId int64
	var err error
	userId, err = l.loginByUsername(in.Username, in.Pass)
	if err != nil {
		return nil, err
	}
	//2、生成token
	ll := NewGenerateTokenLogic(l.ctx, l.svcCtx)
	resp, err := ll.GenerateToken(&pb.GenerateTokenReq{
		UserId: userId,
	})
	if err != nil {
		return nil, errors.Wrapf(ErrGenerateTokenError, "IdentityRpc.GenerateToken userId : %d", userId)
	}
	return &pb.ValidateLoginResp{
		AccessToken:  resp.AccessToken,
		AccessExpire: resp.AccessExpire,
		RefreshAfter: resp.RefreshAfter,
	}, nil

}
// 用户登录
func (l *LoginLogic) loginByUsername(username, pass string) (int64, error) {
	user, err := l.svcCtx.UserModel.FindOneByUsername(username)
	if err != nil && err != model.ErrNotFound {
		return 0, errors.Wrapf(perr.NewErrCode(perr.DB_ERROR), "根据手机号查询用户信息失败，username:%s,err:%v", username, err)
	}
	fmt.Println(user)
	if user == nil {
		return 0, errors.Wrapf(ErrUserNoExistsError, "username:%s", username)
	}
	if !(pass == user.Password) {
		return 0, errors.Wrap(ErrUsernamePwdError, "密码匹配出错")
	}
	return user.Id, nil
}
