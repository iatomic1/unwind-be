package config

import "github.com/spf13/viper"

type Config struct {
	HttpAddress string `mapstructure:"HTTP_SERVER_ADDRESS"`
	DbUrl       string `mapstructure:"DB_URL"`
	DbType      string `mapstructure:"DB_TYPE"`
	JwtKey      string `mapstructure:"JWT_KEY"`
	Environment string `mapstructure:"ENVIROMENT"`
	Host        string `mapstructure:"HOST"`
	// REDIS                  string `mapstructure:"REDIS"`
	RedisAddress  string `mapstructure:"REDIS_ADDRESS"`
	RedisPassword string `mapstructure:"REDIS_PASSWORD"`
	RedisUsername string `mapstructure:"REDIS_USERNAME"`
	RedisDB       int    `mapstructure:"REDIS_DB"`
}

func Load(path string) (*Config, error) {
	return LoadEnvironmentVariables(path, ".env")
}

func LoadTest(path string) (*Config, error) {
	return LoadEnvironmentVariables(path, ".env.test")
}

func LoadEnvironmentVariables(p string, env string) (*Config, error) {
	cfg := Config{}

	viper.AddConfigPath(p)
	viper.SetConfigName(env)
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	err = viper.Unmarshal(&cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
