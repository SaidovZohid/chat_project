package config

import (
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	GrpcPort      string
	Postgres      PostgresConfig
	Redis         Redis
	AuthSecretKey string

	NotificationServiceGrpcPort string
	NotificationServiceHost     string
}

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
}

type Redis struct {
	Addr string
}

func Load(path string) Config {
	godotenv.Load(path + "/.env") // load .env file if it exists

	conf := viper.New()
	conf.AutomaticEnv()

	cfg := Config{
		GrpcPort: conf.GetString("GRPC_PORT"),
		Postgres: PostgresConfig{
			Host:     conf.GetString("POSTGRES_HOST"),
			Port:     conf.GetString("POSTGRES_PORT"),
			User:     conf.GetString("POSTGRES_USER"),
			Password: conf.GetString("POSTGRES_PASSWORD"),
			Database: conf.GetString("POSTGRES_DATABASE"),
		},
		Redis: Redis{
			Addr: conf.GetString("REDIS_ADDR"),
		},
		AuthSecretKey:               conf.GetString("AUTH_SECRET_KEY"),
		NotificationServiceHost:     conf.GetString("NOTIFICATION_SERVICE_HOST"),
		NotificationServiceGrpcPort: conf.GetString("NOTIFICATION_SERVICE_GRPC_PORT"),
	}

	return cfg
}
