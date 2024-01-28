package config

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

const formatJSON = "json"

type Config struct {
	Server struct {
		Host        string `envconfig:"SERVER_HOST" default:":9000"`
		MetricsBind string `envconfig:"BIND_METRICS" default:":9090"`
		HealthHost  string `envconfig:"BIND_HEALTH" default:":9091"`
	}

	Service struct {
		LogLevel  string `envconfig:"LOGGER_LEVEL" default:"debug"`
		LogFormat string `envconfig:"LOGGER_FORMAT" default:"console"`
	}

	PostgresDB struct {
		Addr        string `envconfig:"PG_ADDR"`
		Port        uint16 `envconfig:"PG_PORT"`
		ReplicaAddr string `envconfig:"PG_REPLICA_ADDR"`
		DB          string `envconfig:"PG_DB"`
		User        string `envconfig:"PG_USER"`
		Password    string `envconfig:"PG_PASSWORD"`
	}

	PgMaxConn               int    `envconfig:"PG_MAX_CONN" default:"10"`
	PgMaxIdleLifetime       string `envconfig:"PG_MAX_IDLE_LIFETIME" default:"30s"`
	PgMaxLifetime           string `envconfig:"PG_MAX_LIFETIME" default:"3m"`
	PgPrepareCacheCap       int    `envconfig:"PG_PREPARE_CACHE_CAP" default:"128"`
	PgAllowReadFromReplicas bool   `envconfig:"PG_ALLOW_READ_FROM_REPLICAS" default:"true"`
}

func (cfg Config) Logger() (logger zerolog.Logger) {
	level := zerolog.InfoLevel
	if newLevel, err := zerolog.ParseLevel(cfg.Service.LogLevel); err == nil {
		level = newLevel
	}

	var out io.Writer = os.Stdout
	if cfg.Service.LogFormat != formatJSON {
		out = zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.StampMicro}
	}
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	return zerolog.New(out).Level(level).With().Caller().Timestamp().Logger()
}

func (cfg Config) PgxConnConfig() *pgxpool.Config {
	c, err := pgxpool.ParseConfig(fmt.Sprintf(
		"host=%s port=%d dbname=%s sslmode=disable user=%s password=%s pool_max_conns=%d",
		cfg.PostgresDB.Addr,
		cfg.PostgresDB.Port,
		cfg.PostgresDB.DB,
		cfg.PostgresDB.User,
		cfg.PostgresDB.Password,
		cfg.PgMaxConn,
	))
	if err != nil {
		panic(err)
	}

	return c
}

func Parse() (*Config, error) {
	var cfg = &Config{}
	err := envconfig.Process("", cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
