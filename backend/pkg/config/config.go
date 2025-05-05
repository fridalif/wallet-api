package config

import (
	"backend/pkg/customerror"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DbHost     string
	DbPort     string
	DbUser     string
	DbPassword string
	DbName     string
	WebHost    string
	WebPort    string
}

func NewConfig(dotenvPath string) (*Config, error) {
	err := godotenv.Load(dotenvPath)
	if err != nil {
		return &Config{}, customerror.NewError("config.NewConfig", "", err.Error())
	}
	var config Config
	config.DbHost = os.Getenv("DB_HOST")
	if config.DbHost == "" {
		return &Config{}, customerror.NewError("config.NewConfig", "", "DB_HOST incorrect")
	}
	config.DbPort = os.Getenv("DB_PORT")
	if config.DbPort == "" {
		return &Config{}, customerror.NewError("config.NewConfig", "", "DB_PORT incorrect")
	}
	config.DbUser = os.Getenv("DB_USER")
	if config.DbUser == "" {
		return &Config{}, customerror.NewError("config.NewConfig", "", "DB_USER incorrect")
	}
	config.DbPassword = os.Getenv("DB_PASSWORD")
	if config.DbPassword == "" {
		return &Config{}, customerror.NewError("config.NewConfig", "", "DB_PASSWORD incorrect")
	}
	config.DbName = os.Getenv("DB_NAME")
	if config.DbName == "" {
		return &Config{}, customerror.NewError("config.NewConfig", "", "DB_NAME incorrect")
	}
	config.WebHost = os.Getenv("WEB_HOST")
	if config.WebHost == "" {
		return &Config{}, customerror.NewError("config.NewConfig", "", "WEB_HOST incorrect")
	}
	config.WebPort = os.Getenv("WEB_PORT")
	if config.WebPort == "" {
		return &Config{}, customerror.NewError("config.NewConfig", "", "WEB_PORT incorrect")
	}
	return &config, nil
}
