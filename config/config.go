package config

type ServerConfig struct {
	Host   string
	Port   string
	Scheme string
}

type Config struct {
	Server ServerConfig
}

func NewConfig() Config {
	//goland:noinspection HttpUrlsUsage
	return Config{
		Server: ServerConfig{
			Host:   "localhost",
			Port:   ":8080",
			Scheme: "http://",
		},
	}
}
