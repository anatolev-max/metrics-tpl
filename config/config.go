package config

import "go.uber.org/zap/zapcore"

type ServerConfig struct {
	Host   string
	Port   string
	Scheme string
}

type Config struct {
	Server   ServerConfig
	LogLevel string
}

func NewConfig() Config {
	//goland:noinspection HttpUrlsUsage
	return Config{
		Server: ServerConfig{
			Host:   "localhost",
			Port:   ":8080",
			Scheme: "http://",
		},
		LogLevel: zapcore.InfoLevel.String(),
	}
}
