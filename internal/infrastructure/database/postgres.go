package database

import (
	"context"
	"fmt"
	"time"

	"github.com/damndelion/test_task_kami/configs"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgresDB(config configs.Postgres) (*pgxpool.Pool, error) {
	if config.Host == "" || config.Port == "" || config.Username == "" || config.Password == "" || config.DBName == "" || config.SSLMode == "" {
		return nil, fmt.Errorf("invalid database configuration")
	}

	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s",
		config.Username, config.Password, config.Host, config.Port, config.DBName, config.SSLMode)

	configPool, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database configuration: %w", err)
	}

	//Default values adjust if needed
	configPool.MaxConns = 4
	configPool.MinConns = 0
	configPool.MaxConnLifetime = 0
	configPool.MaxConnIdleTime = 0
	configPool.HealthCheckPeriod = 30 * time.Second

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	dbPool, err := pgxpool.NewWithConfig(ctx, configPool)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	err = dbPool.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return dbPool, nil
}
