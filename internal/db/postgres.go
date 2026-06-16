package db

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/bobfive1/user-management-api/internal/config"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func SetupClientPostgres(configApp config.AppConfig) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s", configApp.PostgresConfig.Username, configApp.PostgresConfig.Password, configApp.PostgresConfig.Host, configApp.PostgresConfig.Port, configApp.PostgresConfig.Database, configApp.PostgresConfig.SSLMode)

	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, err
	}

	config.MaxConns = 10
	config.MinConns = 1

	hostname, _ := os.Hostname()
	if hostname == "" {
		hostname = "Unknow-host"
	}
	config.ConnConfig.RuntimeParams["application_name"] = fmt.Sprintf("%s-%s", configApp.App.Name, hostname)

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, err
	}

	return pool, nil
}

func StopPostgres(pool *pgxpool.Pool) {
	fmt.Println("Postgres pool closing")
	pool.Close()

}
