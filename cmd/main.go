package main

import (
	"try-go-clickhouse/internal/repository"
	"try-go-clickhouse/internal/service/ginrouter"
	"try-go-clickhouse/internal/util/envconf"
	"try-go-clickhouse/internal/util/zaplog"
)

func main() {
	env := envconf.New()
	log := zaplog.New()

	repo, err := repository.New(env, log)
	if err != nil {
		log.Sugar().Panic(err)
	}

	r := ginrouter.New(repo, env, log)
	log.Info(env.AppName + " is starting")

	if err := r.Run(":" + env.Port); err != nil {
		log.Sugar().Panic(err)
	}
}

/*
https://github.com/ClickHouse/clickhouse-go/tree/main/examples/clickhouse_api
https://www.pixelstech.net/article/1676119587-A-simple-tutorial-on-GoLang-connecting-to-Clickhouse
*/
