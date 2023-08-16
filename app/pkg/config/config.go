package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"log"
	"sync"
)

/// Конфигурация приложения \\\

type Config struct {
	HTTP struct {
		Host         string `yaml:"host" env:"HTTP-HOST"`
		Port         string `yaml:"port" env:"HTTP-PORT"`
		ReadTimeout  int    `yaml:"read_timeout" env:"HTTP-READ-TIMEOUT"`
		WriteTimeout int    `yaml:"write_timeout" env:"HTTP-WRITE-TIMEOUT"`
	} `yaml:"http"`
	PostgreSQL struct {
		DSN               string `env:"DATABASE_DSN" env-required:"true"`
		RequestTimeout    int    `yaml:"request_timeout" env-default:"5"`
		ConnectionTimeout int    `yaml:"connection_timeout" env-default:"10"`
		ShutdownTimeout   int    `yaml:"shutdown_timeout" env-default:"5"`
	} `yaml:"postgresql" env-required:"true"`
	JWT struct {
		AccessExpirationMinutes int16  `yaml:"access_expiration_minutes"`
		RefreshExpirationDays   int16  `yaml:"refresh_expiration_days"`
		AccessTokenSecretKey    string `yaml:"access_token_secret_key"`
		RefreshTokenSecretKey   string `yaml:"refresh_token_secret_key"`
	} `yaml:"jwt"`
	MAIL struct {
		MailAddress  string `env:"MAIL_ADD" env-required:"true""`
		MailPassword string `env:"MAIL_PAS" env-required:"true"`
	}
}

// / Функция для получения конфигурации приложения из файла config.yml \\\
var cfg Config

func GetConfig(configPath string, dotenvPath string) *Config {
	var once sync.Once

	if err := godotenv.Load(dotenvPath); err != nil {
		log.Fatalf("could not load .env file: %v", err)
	}

	once.Do(func() {
		if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
			log.Fatalf("config file does not exist: %v", err)
		}
	})
	return &cfg
}
