package main

import (
	"flag"
	"fmt"
	"ydd_finance/app/identity/rpc/internal/config"
	"ydd_finance/app/identity/rpc/internal/server"
	"ydd_finance/app/identity/rpc/internal/svc"
	"ydd_finance/app/identity/rpc/pb"
	"ydd_finance/common/interceptor/grpcinterceptor"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/identity.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)
	srv := server.NewIdentityServer(ctx)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		pb.RegisterIdentityServer(grpcServer, srv)

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	// 添加拦截器，将入参和出参的日记输出,可以定义多个拦截器，参考https://blog.csdn.net/EDDYCJY/article/details/102426265
	s.AddUnaryInterceptors(grpcinterceptor.LoggerInterceptor)

	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
