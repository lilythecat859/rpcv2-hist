package config

import (
	"fmt"
	"net"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	LogLevel      string
	JSONRPCListen string
	RESTListen    string
	GRPCListen    string
	ClickHouse    ClickHouseConfig
}

type ClickHouseConfig struct {
	Addr     string
	Database string
	User     string
	Password string
	AsyncInsert bool
	MaxOpenConns int
	MaxIdleConns int
	ConnMaxLifetime time.Duration
}

func Load() (*Config, error) {
	v := viper.New()
	v.SetEnvPrefix("RPCV2")
	v.AutomaticEnv()

	setDefaults(v)

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshal: %w", err)
	}

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("validate: %w", err)
	}
	return &cfg, nil
}

func setDefaults(v *viper.Viper) {
	v.SetDefault("LogLevel", "info")
	v.SetDefault("JSONRPCListen", "0.0.0.0:8899")
	v.SetDefault("RESTListen", "0.0.0.0:8080")
	v.SetDefault("GRPCListen", "0.0.0.0:9090")

	v.SetDefault("ClickHouse.Addr", "clickhouse://127.0.0.1:9000")
	v.SetDefault("ClickHouse.Database", "solana")
	v.SetDefault("ClickHouse.User", "default")
	v.SetDefault("ClickHouse.Password", "")
	v.SetDefault("ClickHouse.AsyncInsert", true)
	v.SetDefault("ClickHouse.MaxOpenConns", 32)
	v.SetDefault("ClickHouse.MaxIdleConns", 16)
	v.SetDefault("ClickHouse.ConnMaxLifetime", 30*time.Minute)
}

func (c *Config) validate() error {
	for _, addr := range []string{c.JSONRPCListen, c.RESTListen, c.GRPCListen} {
		if _, err := net.Listen("tcp", addr); err == nil {
			continue
		}
		if _, _, err := net.SplitHostPort(addr); err != nil {
			return fmt.Errorf("invalid addr %q: %w", addr, err)
		}
	}
	return nil
}