package config

import (
	"github.com/zeromicro/go-zero/zrpc"
	"github.com/zeromicro/go-zero/rest"
)


type Config struct {
	rest.RestConf
	IdentityRpc zrpc.RpcClientConf
	NoAuthUrls []string
	JwtAuth struct {
		AccessSecret string
	}
}
