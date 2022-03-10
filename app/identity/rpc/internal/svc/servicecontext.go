package svc

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"ydd_finance/app/identity/model"
	//"strconv"
	"ydd_finance/app/identity/rpc/internal/config"
)

type ServiceContext struct {
	Config config.Config
	RedisClient *redis.Redis
	UserModel     model.UserModel
	UserAuthModel model.UserAuthModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	// redis.WithPass(strconv.Itoa(c.Redis.Pass))
	redisPool := redis.New(c.Redis.Host,redis.WithPass(c.Redis.Pass))
	return &ServiceContext{
		Config: c,
		RedisClient: redisPool,
		UserAuthModel: model.NewUserAuthModel(sqlx.NewMysql(c.DB.DataSource), c.Cache),
		UserModel:     model.NewUserModel(sqlx.NewMysql(c.DB.DataSource), c.Cache),
	}
}
