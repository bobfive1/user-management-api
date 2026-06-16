package db

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/bobfive1/user-management-api/internal/config"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func SetupClientPostgres(configApp config.AppConfig) (*pgxpool.Pool, error) {
	connectTimeout := configApp.PostgresConfig.ConnectTimeout
	if connectTimeout <= 0 {
		connectTimeout = 10 * time.Second
	}

	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout)
	defer cancel()

	connURL := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(configApp.PostgresConfig.Username, configApp.PostgresConfig.Password),
		Host:   net.JoinHostPort(configApp.PostgresConfig.Host, strconv.Itoa(configApp.PostgresConfig.Port)),
		Path:   configApp.PostgresConfig.Database,
	}
	query := connURL.Query()
	if configApp.PostgresConfig.SSLMode != "" {
		query.Set("sslmode", configApp.PostgresConfig.SSLMode)
	}
	connURL.RawQuery = query.Encode()

	poolConfig, err := pgxpool.ParseConfig(connURL.String())
	if err != nil {
		return nil, err
	}

	poolConfig.ConnConfig.ConnectTimeout = connectTimeout
	if configApp.PostgresConfig.MaxConns > 0 {
		poolConfig.MaxConns = int32(configApp.PostgresConfig.MaxConns)
	}
	if configApp.PostgresConfig.MinConns > 0 {
		poolConfig.MinConns = int32(configApp.PostgresConfig.MinConns)
	}

	hostname, _ := os.Hostname()
	if hostname == "" {
		hostname = "Unknow-host"
	}
	poolConfig.ConnConfig.RuntimeParams["application_name"] = fmt.Sprintf("%s-%s", configApp.App.Name, hostname)

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, err
	}

	return pool, nil
}

func StopPostgres(pool *pgxpool.Pool) {
	fmt.Println("Postgres pool closing")
	pool.Close()

}
