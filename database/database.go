package database

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"stew/types"
	"time"
)

var Pool *pgxpool.Pool = nil

func SetTimeout(seconds time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), time.Second*seconds)
}

func LoadDatabase(conf types.DatabaseConfig) *pgxpool.Pool {
	poolConf, err := pgxpool.ParseConfig(
		fmt.Sprintf("postgres://%s:%s@%s:%d/%s", conf.SQLUsername, conf.SQLPassword, conf.SQLHost, conf.SQLPort, conf.SQLDatabase))
	if err != nil {
		panic(err)
	}
	poolConf.MinConns = conf.MinConns
	poolConf.MaxConns = conf.MaxConns
	poolConf.MaxConnIdleTime = 5 * time.Minute

	pool, err := pgxpool.NewWithConfig(context.Background(), poolConf)
	Pool = pool
	if err != nil {
		panic(err)
	}
	return pool
}

func ConnectDatabase(db *pgxpool.Pool) {
	c, cancel := SetTimeout(5)
	defer cancel()

	err := db.Ping(c)
	if err != nil {
		panic(err)
	}
}
