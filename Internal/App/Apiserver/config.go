package Apiserver

type Config struct {
	BindAddr string `toml:"bind_addr"` // адрес на котором запускаем вебсервер
	LogLevel string `toml:"log_level"`
}

func NewConfig() *Config {
	return &Config{
		BindAddr: ":8085",
		LogLevel: "debug",
	}
}
