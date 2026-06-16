package config

import (
	"bytes"
	"strings"
	"time"

	_ "embed"

	"github.com/spf13/viper"
)

var (
	//go:embed config.yaml
	conf []byte
)

type AppConfig struct {
	App            AppInfo          `mapstructure:"app"`
	ServerAPI      ServerAPIConfig  `mapstructure:"apiServer"`
	Logging        LoggingConfig    `mapstructure:"logging"`
	PostgresConfig PostgreSQLConfig `mapstructure:"postgresql"`
}

type AppInfo struct {
	Name string `mapstructure:"name"`
}

type ServerAPIConfig struct {
	Port            string        `mapstructure:"address"`
	ReadTimeout     time.Duration `mapstructure:"readTimeout"`
	WriteTimeout    time.Duration `mapstructure:"writeTimeout"`
	IdleTimeout     time.Duration `mapstructure:"idleTimeout"`
	ShutdownTimeout time.Duration `mapstructure:"shutdownTimeout"`
}

type LoggingConfig struct {
	Level string `mapstructure:"level"`
}

type PostgreSQLConfig struct {
	Host           string        `mapstructure:"host"`
	Port           int           `mapstructure:"port"`
	Username       string        `mapstructure:"username"`
	Password       string        `mapstructure:"password"`
	Database       string        `mapstructure:"database"`
	SSLMode        string        `mapstructure:"sslmode"`
	MaxConns       int           `mapstructure:"max-conns"`
	MinConns       int           `mapstructure:"min-conns"`
	ConnectTimeout time.Duration `mapstructure:"connect-timeout"`
}

func LoadConfig() (*AppConfig, error) {
	viper.SetConfigType("yaml")
	viper.SetEnvPrefix("")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	viper.AutomaticEnv()
	viper.AllowEmptyEnv(true)
	viper.RegisterAlias("posgresql", "postgresql")

	// read config
	if err := viper.ReadConfig(bytes.NewBuffer(conf)); err != nil {
		return nil, err
	}

	appConfig := &AppConfig{}
	if err := viper.Unmarshal(appConfig); err != nil {
		return nil, err
	}

	return appConfig, nil
}
