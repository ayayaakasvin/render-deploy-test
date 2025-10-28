package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

const (
	configPathEnvKey 	= "CONFIG_PATH"
	_SMTPEnvKey 		= "SMTP"
)

// Config represents the configuration structure
type Config struct {
	HTTPServer 					`yaml:"http_server"																env-required:"true"`
	SMTPConfig					`yaml:"smtp"																	env-required:"true""`
}

type HTTPServer struct {
	Address     string        	`yaml:"address"			env-default:"localhost:8080`
	Timeout     time.Duration 	`yaml:"timeout" 																env-required:"true"`
	IdleTimeout time.Duration 	`yaml:"idle_timeout" 															env-required:"true"`
}


type SMTPConfig struct {
	Username 		string			`yaml:"username"															env-required:"true"`
	Password		string
	Host			string			`yaml:"host"																env-required:"true"`
	Port			int				`yaml:"port"																env-required:"true"`
}

// MustLoadConfig loads the configuration from the specified path
func MustLoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("failed to load .env: %s", err.Error())
	}

	configPath := os.Getenv(configPathEnvKey)
	if configPath == "" {
		log.Fatalf("%s is not set up", configPathEnvKey)
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file %s does not exist: %s", configPath, err.Error())
	}

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("failed to read config file: %s", err.Error())
	}
	
	cfg.SMTPConfig.Password = os.Getenv(_SMTPEnvKey)

	return &cfg
}
