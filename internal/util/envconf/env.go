package envconf

import (
	"github.com/kelseyhightower/envconfig"
)

type Spec struct {
	AppName string `envconfig:"APP_NAME" default:"try-go-clickhouse"`
	Port    string `envconfig:"PORT" default:"8080"`

	ClickhouseAddr     []string `envconfig:"CLICKHOUSE_ADDR" default:"127.0.0.1:9000"`
	ClickhouseDatabase string   `envconfig:"CLICKHOUSE_DATABASE" default:"default"`
	ClickhouseUsername string   `envconfig:"CLICKHOUSE_USERNAME" default:"default"`
	ClickhousePassword string   `envconfig:"CLICKHOUSE_PASSWORD"`
}

func New() *Spec {
	var s Spec
	envconfig.MustProcess("", &s)
	return &s
}
