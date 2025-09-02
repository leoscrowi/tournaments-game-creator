package config

import "os"

type Config struct {
	GrpcConfig     GrpcConfig
	DatabaseConfig DatabaseConfig
}

type GrpcConfig struct {
	Auth    string
	Storage string
	Port    string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SslMode  string
}

func MustLoad() *Config {
	return &Config{
		GrpcConfig: GrpcConfig{
			Auth:    os.Getenv("GRPC_AUTH"),
			Storage: os.Getenv("GRPC_STORAGE"),
			Port:    os.Getenv("GRPC_PORT"),
		},
		DatabaseConfig: DatabaseConfig{
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			Name:     os.Getenv("DB_NAME"),
			SslMode:  os.Getenv("DB_SSL_MODE"),
		},
	}
}
