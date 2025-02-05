package config

import (
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	// Config -.
	Config struct {
		App   `yaml:"app"`
		HTTP  `yaml:"http"`
		Log   `yaml:"logger"`
		PG    `yaml:"postgres"`
		Kafka `yaml:"kafka"`
	}

	// App -.
	App struct {
		Name    string `env-required:"true" yaml:"name"    env:"APP_NAME"`
		Version string `env-required:"true" yaml:"version" env:"APP_VERSION"`
	}

	// HTTP -.
	HTTP struct {
		Port string `env-required:"true" yaml:"port" env:"HTTP_PORT"`
	}

	// Log -.
	Log struct {
		Level string `env-required:"true" yaml:"log_level"   env:"LOG_LEVEL"`
	}

	// PG -.
	PG struct {
		PoolMax int    `env-required:"true" yaml:"pool_max" env:"PG_POOL_MAX"`
		URL     string `env-required:"true"                 env:"PG_URL"`
	}

	Kafka struct {
		KAFKA_BROKER string `env-required:"true"  yaml:"KAFKA_BROKER"  env:"KAFKA_BROKER"`
		EVENT_TOPIC  string `env-required:"true"  yaml:"EVENT_TOPIC"  env:"EVENT_TOPIC"`
		RETRY_TOPIC  string `env-required:"true"  yaml:"RETRY_TOPIC"  env:"RETRY_TOPIC"`
		DLQ_TOPIC    string `env-required:"true"  yaml:"DLQ_TOPIC"  env:"DLQ_TOPIC"`
	}
)

func NewConfig() (*Config, error) {
	cfg := &Config{}

	workingDir, err2 := os.Getwd()
	if err2 != nil {
		fmt.Println("Error:", err2)

	}

	// Çalışma dizinini yazdır
	fmt.Println("Current working directory:", workingDir)

	err := cleanenv.ReadConfig("./config/config.yml", cfg)
	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
