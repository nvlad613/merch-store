package config

import (
	"fmt"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"go.uber.org/zap"
	"strings"
)

type Config struct {
	Database DatabaseConnection
	Logger   zap.Config
	Server   ServerParams
}

type ServerParams struct {
	Hostname           string `koanf:"host"`
	Port               int    `koanf:"port"`
	IdleTimeoutSec     int    `koanf:"idle-timeout"`
	ShutdownTimeoutSec int    `koanf:"shutdown-timeout"`
}

type DatabaseConnection struct {
	Hostname   string    `koanf:"host"`
	Port       int       `koanf:"port"`
	Database   string    `koanf:"name"`
	User       BasicAuth `koanf:"auth"`
	TLSEnabled bool      `koanf:"tls-enabled"`
}

type BasicAuth struct {
	Login    string `koanf:"login"`
	Password string `koanf:"password"`
}

func Load() (Config, error) {
	var (
		zero Config
		k    = koanf.New(".")
	)

	// Load main config
	if err := k.Load(file.Provider("application.yml"), yaml.Parser()); err != nil {
		return zero, fmt.Errorf("failed to load configuration: %w", err)
	}

	// Merge with local config
	_ = k.Load(file.Provider("local.yml"), yaml.Parser())

	// Merge with env variables
	_ = k.Load(env.Provider("MERCHSTORE_", ".", func(s string) string {
		trimmed := strings.TrimPrefix(s, "MERCHSTORE_")
		lowerCased := strings.Replace(strings.ToLower(trimmed), "_", ".", -1)
		return "application." + lowerCased
	}), nil)

	var (
		loggerConfig zap.Config
		dbConfig     DatabaseConnection
		serverConfig ServerParams
	)
	if err := k.Unmarshal("application.db", &dbConfig); err != nil {
		return zero, fmt.Errorf("failed to unmarshal db configuration: %w", err)
	}
	if err := k.Unmarshal("application.server", &serverConfig); err != nil {
		return zero, fmt.Errorf("failed to unmarshal server configuration: %w", err)
	}
	if err := k.UnmarshalWithConf("application.logger", &loggerConfig, koanf.UnmarshalConf{Tag: "yaml"}); err != nil {
		return zero, fmt.Errorf("failed to unmarshal logger configuration: %w", err)
	}

	return Config{
		Database: dbConfig,
		Logger:   loggerConfig,
		Server:   serverConfig,
	}, nil
}
