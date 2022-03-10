package logic

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"ydd_finance/app/identity/rpc/internal/svc"
	"ydd_finance/app/identity/rpc/pb"
	"ydd_finance/common/ctxdata"
	"ydd_finance/common/globalvar"
	"ydd_finance/common/perr"
	"github.com/golang-jwt/jwt/v4"
	"github.com/zeromicro/go-zero/core/logx"
	"time"
)
var ErrGenerateTokenError = perr.NewErrMsg("生成token失败")

type GenerateTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGenerateTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GenerateTokenLogic {
	return &GenerateTokenLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

//  生成token，只针对用户服务开放访问
func (l *GenerateTokenLogic) GenerateToken(in *pb.GenerateTokenReq) (*pb.GenerateTokenResp, error) {
	// todo: add your logic here and delete this line
	now := time.Now().Unix()
	accessExpire := l.svcCtx.Config.JwtAuth.AccessExpire
	accessToken, err := l.getJwtToken(l.svcCtx.Config.JwtAuth.AccessSecret, now, accessExpire, in.UserId)
	if err != nil {
		return nil, errors.Wrapf(ErrGenerateTokenError, "getJwtToken err userId:%d , err:%v", in.UserId, err)
	}

	// 存入redis.
	userTokenKey := fmt.Sprintf(globalvar.CacheUserTokenKey, in.UserId)
	err = l.svcCtx.RedisClient.Setex(userTokenKey, accessToken, int(accessExpire))
	if err != nil {
		return nil, errors.Wrapf(ErrGenerateTokenError, "SetnxEx err userId:%d, err:%v", in.UserId, err)
	}
	return &pb.GenerateTokenResp{
		AccessToken:  accessToken,
		AccessExpire: now + accessExpire,
		RefreshAfter: now + accessExpire/2,
	}, nil
	return &pb.GenerateTokenResp{}, nil
}

// 生成token 文档https://pkg.go.dev/github.com/golang-jwt/jwt/v4
// 具体文档https://pkg.go.dev/github.com/golang-jwt/jwt/v4#example-NewWithClaims-CustomClaimsType

func (l *GenerateTokenLogic) getJwtToken(secretKey string, iat, seconds, userId int64) (string, error) {
	// 旧版https://blog.csdn.net/cbmljs/article/details/86072395
	// 新版https://blog.csdn.net/weixin_44294408/article/details/122095919
	claims := make(jwt.MapClaims)
	claims["exp"] = iat + seconds // 过期时间
	claims["iat"] = iat // 当前时间
	claims[ctxdata.CtxKeyJwtUserId] = userId
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims
	return token.SignedString([]byte(secretKey))
}
