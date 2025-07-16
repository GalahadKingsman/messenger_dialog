package config

// Config содержит всю конфигурацию приложения
type DBConfig struct {
	Host     string `env:"HOST"`
	Port     string `env:"PORT"`
	User     string `env:"USER"`
	Password string `env:"PASSWORD"`
	Name     string `env:"NAME"`
}

type Config struct {
	DB       DBConfig `envPrefix:"DB_"`
	GRPCPort int      `env:"GRPC_PORT"`
}
