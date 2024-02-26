package repository

import (
	"context"
	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"go.uber.org/zap"
	"try-go-clickhouse/internal/model"
	"try-go-clickhouse/internal/util/envconf"
)

var _ model.Repository = (*Repo)(nil)

type Repo struct {
	conn driver.Conn
	log  *zap.Logger
}

func New(env *envconf.Spec, log *zap.Logger) (model.Repository, error) {
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: env.ClickhouseAddr,
		Auth: clickhouse.Auth{
			Database: env.ClickhouseDatabase,
			Username: env.ClickhouseUsername,
			Password: env.ClickhousePassword,
		},
	})
	if err != nil {
		log.Sugar().Error(err)
		return nil, err
	}

	if err = conn.Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS users (
				id UInt32,
				name String,
				age UInt8
			) ENGINE = Memory `); err != nil {
		log.Sugar().Error(err)
		return nil, err
	}

	return &Repo{conn, log}, err
}

func (r *Repo) ListUsers(ctx context.Context) ([]model.User, error) {
	var users []model.User
	if err := r.conn.Select(ctx, &users,
		"select id, name, age from users"); err != nil {
		r.log.Sugar().Error(err)
		return users, err
	}

	return users, nil
}

func (r *Repo) AddUser(user model.User) (model.User, error) {
	if err := r.conn.Exec(context.Background(),
		"insert into users (id, name, age) values ($1, $2, $3)",
		user.Id, user.Name, user.Age); err != nil {
		r.log.Sugar().Error(err)
		return user, err
	}

	return user, nil
}

func (r *Repo) GetUser(id int) (model.User, error) {
	var user model.User
	if err := r.conn.QueryRow(context.Background(),
		"select id, name, age from users where id = $1",
		id).Scan(&user.Id, &user.Name, &user.Age); err != nil {
		r.log.Sugar().Error(err)
		return user, err
	}

	return user, nil
}
