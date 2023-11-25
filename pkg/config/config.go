package config

import (
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type Config struct {
	DBHost          string `mapstructure:"DB_HOST"`
	DBName          string `mapstructure:"DB_NAME"`
	DBUser          string `mapstructure:"DB_USER"`
	DBPort          string `mapstructure:"DB_PORT"`
	DBPassword      string `mapstructure:"DB_PASSWORD"`
	AUTHTOKEN       string `mapstructure:"DB_AUTHTOKEN"`
	ACCOUNTSID      string `mapstructure:"DB_ACCOUNTSID"`
	SERVICESID      string `mapstructure:"DB_SERVICESID"`
	AdminAccessKey  string `mapstructure:"AdminAccessKey"`
	AdminRefreshKey string `mapstructure:"AdminRefreshKey"`
	UserAccessKey   string `mapstructure:"UserAccessKey"`
	UserRefreshKey  string `mapstructure:"UserRefreshKey"`
}

var envs = []string{
	"DB_HOST", "DB_NAME", "DB_USER", "DB_PORT", "DB_PASSWORD",
	"DB_AUTHTOKEN", "DB_ACCOUNTSID", "DB_SERVICESID",
	"AdminAccessKey", "AdminRefreshKey",
	"UserAccessKey", "UserRefreshKey",
}

func LoadConfig() (Config, error) {
	var config Config
	viper.AddConfigPath("./")
	viper.SetConfigFile(".env")

	err := viper.ReadInConfig()
	if err != nil {
		return Config{}, err
	}
	for _, env := range envs {
		if err := viper.BindEnv(env); err != nil {
			return config, err
		}
	}
	if err := viper.Unmarshal(&config); err != nil {
		return config, err
	}
	if err := validator.New().Struct(&config); err != nil {
		return config, err
	}

	return config, nil

}
