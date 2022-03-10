/**
* @program: src
*
* @description:grpc拦截器 参考grpc UnaryServerInterceptor，多个可以github.com/grpc-ecosystem/go-grpc-middleware
*
* @author: Mr.chen
*
* @create: 2022-03-10 10:26
**/

package grpcinterceptor
import (
	"context"
	"ydd_finance/common/perr"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

/**
* @Description - logging：RPC 方法的入参出参的日志输出
*	其类型是一个函数，这个函数有，4 个入参，两个出参，介绍如下
*	ctx context.Context 上下文
*	req interface {} 用户请求的参数
*	info UnaryServerInfo RPC 方法的所有信息，定义如下
**/

func LoggerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	resp, err = handler(ctx, req)
	if err != nil {
		causeErr := errors.Cause(err)                // err类型,获取errors.Wrapf下的第一个参数Cause
		if e, ok := causeErr.(*perr.CodeError); ok { //自定义错误类型，perr库产生的错误，重新定义err
			// 将日记添加前缀打印在控制台上，可以利用kakfa将控制台输出的日记上云
			logx.WithContext(ctx).Errorf("【RPC-SRV-ERR】 %+v", err)
			//重新将错误转换成grpc err, 错误格式为rpc error: code = Code(100001) desc = 账号不存在
			err = status.Error(codes.Code(e.GetErrCode()), e.GetErrMsg())
		} else {
			logx.WithContext(ctx).Errorf("【RPC-SRV-ERR】 %+v", err)
		}
	}
	return resp, err
}