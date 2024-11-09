package config

type ServerConfig struct {
	Host string
	Port string
}

type Config struct {
	Server ServerConfig
}

func NewConfig() Config {
	return Config{
		Server: ServerConfig{
			Host: "http://localhost",
			Port: ":8000",
		},
	}
}
