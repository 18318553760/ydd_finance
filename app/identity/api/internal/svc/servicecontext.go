package svc

import (
	"github.com/zeromicro/go-zero/zrpc"
	"ydd_finance/app/identity/api/internal/config"
	"ydd_finance/app/identity/rpc/identity"
)

type ServiceContext struct {
	Config config.Config
	IdentityRpc identity.Identity
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
		IdentityRpc: identity.NewIdentity(zrpc.MustNewClient(c.IdentityRpc)),
	}
}
