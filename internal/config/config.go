package config

// Config содержит всю конфигурацию приложения
type DBConfig struct {
	Host     string `env:"HOST" envDefault:"localhost"`
	Port     string `env:"PORT" envDefault:"5432"`
	User     string `env:"USER" envDefault:"postgres"`
	Password string `env:"PASSWORD"`
	Name     string `env:"NAME" envDefault:"messenger_users"`
}

type Config struct {
	DB       DBConfig `envPrefix:"DB_"`
	GRPCPort int      `env:"GRPC_PORT" envDefault:"9001"`
}
