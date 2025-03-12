package api

import (
	"fmt"
	"strings"

	"github.com/1755/bookstore-api/api/routers"
	"github.com/1755/bookstore-api/internal/pg"
	"github.com/go-playground/validator/v10"
	"github.com/google/wire"
	"github.com/spf13/viper"
)

type ConfigPath string

type Config struct {
	Logger     *LoggerConfig     `validate:"required"`
	Server     *ServerConfig     `validate:"required"`
	Monitoring *MonitoringConfig `validate:"required"`
	Routers    *routers.Config   `validate:"required"`
	Postgres   *pg.Config        `validate:"required"`
}

func NewConfig(configPath ConfigPath) (Config, error) {

	v := viper.New()
	v.SetConfigFile(string(configPath))
	v.SetEnvPrefix("")
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := v.ReadInConfig(); err != nil {
		return Config{}, fmt.Errorf("failed to read config: %w", err)
	}

	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return Config{}, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	validate := validator.New()
	if err := validate.Struct(config); err != nil {
		return Config{}, fmt.Errorf("config validation failed: %w", err)
	}

	return config, nil
}

var ConfigModule = wire.NewSet(
	NewConfig,
	wire.FieldsOf(
		new(Config),
		"Logger",
		"Server",
		"Monitoring",
		"Routers",
		"Postgres",
	),
)
