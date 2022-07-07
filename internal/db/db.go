package db

import (
	"context"
	"github.com/abrbird/orders-tracking/config"
	"github.com/jackc/pgx/v4/pgxpool"
)

func New(dbConfig config.Database, ctx context.Context) (*pgxpool.Pool, error) {
	//pool, err := pgxpool.Connect(ctx, dbConfig.String())

	cfg, _ := pgxpool.ParseConfig("")
	cfg.ConnConfig.Host = dbConfig.Host
	cfg.ConnConfig.Port = uint16(dbConfig.Port)
	cfg.ConnConfig.User = dbConfig.User
	cfg.ConnConfig.Password = dbConfig.Password
	cfg.ConnConfig.Database = dbConfig.DB
	cfg.ConnConfig.PreferSimpleProtocol = true
	cfg.MaxConns = 20
	pool, err := pgxpool.ConnectConfig(ctx, cfg)

	if err != nil {
		return nil, err
	}

	return pool, nil
}
