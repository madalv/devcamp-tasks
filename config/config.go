package config

import "github.com/spf13/viper"

type Config struct {
	DBSource string `mapstructure:"db_source"`
	DBDriver string `mapstructure:"db_driver"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.SetConfigFile(path)

	err = viper.ReadInConfig()
	if err != nil {
		return
	}
	err = viper.Unmarshal(&config)
	return
}
