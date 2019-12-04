package config

import (
	"github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/spf13/viper"
)

// MysqlConnect represents mysql connection config
type MysqlConnect struct {
	Host     string
	Port     int
	Username string
	Password string
	Database string
}

// Config represents application config
type Config struct {
	Listen string

	MySQL struct {
		Write MysqlConnect
		Read  MysqlConnect
	}
}

func init() {
	viper.SetConfigName("finance-app")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/etc/otvet/")
}

// Load reads config
func Load() (*Config, error) {
	cfg := new(Config)
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}
	err = viper.Unmarshal(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, validateConfig(cfg)
}

func validateConfig(cfg *Config) error {
	return validation.Errors{
		"listen":               validation.Validate(cfg.Listen, validation.Required, is.DialString),
		"mysql.write.host":     validation.Validate(cfg.MySQL.Write.Host, validation.Required),
		"mysql.write.database": validation.Validate(cfg.MySQL.Write.Database, validation.Required),
		"mysql.write.username": validation.Validate(cfg.MySQL.Write.Username, validation.Required),
		"mysql.write.password": validation.Validate(cfg.MySQL.Write.Password, validation.Required),

		"mysql.read.host":     validation.Validate(cfg.MySQL.Read.Host, validation.Required),
		"mysql.read.database": validation.Validate(cfg.MySQL.Read.Database, validation.Required),
		"mysql.read.username": validation.Validate(cfg.MySQL.Read.Username, validation.Required),
		"mysql.read.password": validation.Validate(cfg.MySQL.Read.Password, validation.Required),
	}.Filter()
}
