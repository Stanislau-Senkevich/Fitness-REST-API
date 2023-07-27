package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Port string
	PostgresConfig
}

type PostgresConfig struct {
	DBHost     string `mapstructure:"postgres_host"`
	DBPort     string `mapstructure:"postgres_port"`
	DBName     string `mapstructure:"postgres_db_name"`
	DBUser     string `mapstructure:"postgres_user"`
	DBPassword string
}

func InitConfig() (*Config, error) {
	viper.SetConfigFile("configs/config.yml")

	var cfg Config

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	if err := viper.UnmarshalKey("postgres_config", &cfg.PostgresConfig); err != nil {
		return nil, err
	}

	if err := parseEnv(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func parseEnv(cfg *Config) error {
	//if err := gotenv.Load(); err != nil {
	//	return err
	//}

	if err := viper.BindEnv("postgres_password"); err != nil {
		return err
	}

	cfg.DBPassword = viper.GetString("postgres_password")
	return nil
}
