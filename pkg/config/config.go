package config

import (
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"golang.org/x/exp/slog"
)

// Config app configuration
type Config struct {
	// the shared DB ORM object
	DB *gorm.DB
	// the error thrown be GORM when using DB ORM object
	DBErr error
	// the server port.
	ServerPort string `mapstructure:"PORT"`
	// the data source name (DSN) for connecting to the database. required.
	DSN string `mapstructure:"DSN"`
	// migration files location
	MigrationDir string `mapstructure:"MIGRATION_DIR"`
	// database
	Database string `mapstructure:"DATABASE"`
	// log level
	LogLevel slog.Level `mapstructure:"LOG_LEVEL"`
}

// LoadConfig loads config from files
func LoadConfig(target string) (config Config, err error) {
	viper.AddConfigPath("./envs")
	viper.SetConfigName(target)
	viper.SetConfigType("env")
	viper.AutomaticEnv()
	if err = viper.ReadInConfig(); err != nil {
		return
	}
	err = viper.Unmarshal(&config)
	return
}
