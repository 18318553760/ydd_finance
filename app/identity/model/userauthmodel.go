package model

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/stores/builder"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/stringx"
)

var (
	userAuthFieldNames          = builder.RawFieldNames(&UserAuth{})
	userAuthRows                = strings.Join(userAuthFieldNames, ",")
	userAuthRowsExpectAutoSet   = strings.Join(stringx.Remove(userAuthFieldNames, "`id`", "`create_time`", "`update_time`"), ",")
	userAuthRowsWithPlaceHolder = strings.Join(stringx.Remove(userAuthFieldNames, "`id`", "`create_time`", "`update_time`"), "=?,") + "=?"

	cacheIdentityUserAuthIdPrefix              = "cache:identity:userAuth:id:"
	cacheIdentityUserAuthAuthTypeAuthKeyPrefix = "cache:identity:userAuth:authType:authKey:"
	cacheIdentityUserAuthUserIdAuthTypePrefix  = "cache:identity:userAuth:userId:authType:"
)

type (
	UserAuthModel interface {
		Insert(data *UserAuth) (sql.Result, error)
		FindOne(id int64) (*UserAuth, error)
		FindOneByAuthTypeAuthKey(authType string, authKey string) (*UserAuth, error)
		FindOneByUserIdAuthType(userId int64, authType string) (*UserAuth, error)
		Update(data *UserAuth) error
		Delete(id int64) error
	}

	defaultUserAuthModel struct {
		sqlc.CachedConn
		table string
	}

	UserAuth struct {
		Id         int64     `db:"id"`
		CreateTime time.Time `db:"create_time"`
		UpdateTime time.Time `db:"update_time"`
		DeleteTime time.Time `db:"delete_time"`
		DelState   int64     `db:"del_state"`
		UserId     int64     `db:"user_id"`
		AuthKey    string    `db:"auth_key"`  // 平台唯一id
		AuthType   string    `db:"auth_type"` // 平台类型
	}
)

func NewUserAuthModel(conn sqlx.SqlConn, c cache.CacheConf) UserAuthModel {
	return &defaultUserAuthModel{
		CachedConn: sqlc.NewConn(conn, c),
		table:      "`user_auth`",
	}
}

func (m *defaultUserAuthModel) Insert(data *UserAuth) (sql.Result, error) {
	identityUserAuthIdKey := fmt.Sprintf("%s%v", cacheIdentityUserAuthIdPrefix, data.Id)
	identityUserAuthAuthTypeAuthKeyKey := fmt.Sprintf("%s%v:%v", cacheIdentityUserAuthAuthTypeAuthKeyPrefix, data.AuthType, data.AuthKey)
	identityUserAuthUserIdAuthTypeKey := fmt.Sprintf("%s%v:%v", cacheIdentityUserAuthUserIdAuthTypePrefix, data.UserId, data.AuthType)
	ret, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?)", m.table, userAuthRowsExpectAutoSet)
		return conn.Exec(query, data.DeleteTime, data.DelState, data.UserId, data.AuthKey, data.AuthType)
	}, identityUserAuthUserIdAuthTypeKey, identityUserAuthIdKey, identityUserAuthAuthTypeAuthKeyKey)
	return ret, err
}

func (m *defaultUserAuthModel) FindOne(id int64) (*UserAuth, error) {
	identityUserAuthIdKey := fmt.Sprintf("%s%v", cacheIdentityUserAuthIdPrefix, id)
	var resp UserAuth
	err := m.QueryRow(&resp, identityUserAuthIdKey, func(conn sqlx.SqlConn, v interface{}) error {
		query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", userAuthRows, m.table)
		return conn.QueryRow(v, query, id)
	})
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultUserAuthModel) FindOneByAuthTypeAuthKey(authType string, authKey string) (*UserAuth, error) {
	identityUserAuthAuthTypeAuthKeyKey := fmt.Sprintf("%s%v:%v", cacheIdentityUserAuthAuthTypeAuthKeyPrefix, authType, authKey)
	var resp UserAuth
	err := m.QueryRowIndex(&resp, identityUserAuthAuthTypeAuthKeyKey, m.formatPrimary, func(conn sqlx.SqlConn, v interface{}) (i interface{}, e error) {
		query := fmt.Sprintf("select %s from %s where `auth_type` = ? and `auth_key` = ? limit 1", userAuthRows, m.table)
		if err := conn.QueryRow(&resp, query, authType, authKey); err != nil {
			return nil, err
		}
		return resp.Id, nil
	}, m.queryPrimary)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultUserAuthModel) FindOneByUserIdAuthType(userId int64, authType string) (*UserAuth, error) {
	identityUserAuthUserIdAuthTypeKey := fmt.Sprintf("%s%v:%v", cacheIdentityUserAuthUserIdAuthTypePrefix, userId, authType)
	var resp UserAuth
	err := m.QueryRowIndex(&resp, identityUserAuthUserIdAuthTypeKey, m.formatPrimary, func(conn sqlx.SqlConn, v interface{}) (i interface{}, e error) {
		query := fmt.Sprintf("select %s from %s where `user_id` = ? and `auth_type` = ? limit 1", userAuthRows, m.table)
		if err := conn.QueryRow(&resp, query, userId, authType); err != nil {
			return nil, err
		}
		return resp.Id, nil
	}, m.queryPrimary)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultUserAuthModel) Update(data *UserAuth) error {
	identityUserAuthAuthTypeAuthKeyKey := fmt.Sprintf("%s%v:%v", cacheIdentityUserAuthAuthTypeAuthKeyPrefix, data.AuthType, data.AuthKey)
	identityUserAuthUserIdAuthTypeKey := fmt.Sprintf("%s%v:%v", cacheIdentityUserAuthUserIdAuthTypePrefix, data.UserId, data.AuthType)
	identityUserAuthIdKey := fmt.Sprintf("%s%v", cacheIdentityUserAuthIdPrefix, data.Id)
	_, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, userAuthRowsWithPlaceHolder)
		return conn.Exec(query, data.DeleteTime, data.DelState, data.UserId, data.AuthKey, data.AuthType, data.Id)
	}, identityUserAuthIdKey, identityUserAuthAuthTypeAuthKeyKey, identityUserAuthUserIdAuthTypeKey)
	return err
}

func (m *defaultUserAuthModel) Delete(id int64) error {
	data, err := m.FindOne(id)
	if err != nil {
		return err
	}

	identityUserAuthUserIdAuthTypeKey := fmt.Sprintf("%s%v:%v", cacheIdentityUserAuthUserIdAuthTypePrefix, data.UserId, data.AuthType)
	identityUserAuthIdKey := fmt.Sprintf("%s%v", cacheIdentityUserAuthIdPrefix, id)
	identityUserAuthAuthTypeAuthKeyKey := fmt.Sprintf("%s%v:%v", cacheIdentityUserAuthAuthTypeAuthKeyPrefix, data.AuthType, data.AuthKey)
	_, err = m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
		return conn.Exec(query, id)
	}, identityUserAuthIdKey, identityUserAuthAuthTypeAuthKeyKey, identityUserAuthUserIdAuthTypeKey)
	return err
}

func (m *defaultUserAuthModel) formatPrimary(primary interface{}) string {
	return fmt.Sprintf("%s%v", cacheIdentityUserAuthIdPrefix, primary)
}

func (m *defaultUserAuthModel) queryPrimary(conn sqlx.SqlConn, v, primary interface{}) error {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", userAuthRows, m.table)
	return conn.QueryRow(v, query, primary)
}
