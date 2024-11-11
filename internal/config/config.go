package config

import "github.com/spf13/viper"

type Config struct {
	Environment           string  `mapstructure:"ENVIROMENT"`
	HttpAddress           string  `mapstructure:"HTTP_SERVER_ADDRESS"`
	DbURL                 string  `mapstructure:"DB_URL"`
	DbType                string  `mapstructure:"DB_TYPE"`
	AccessJwtKey          string  `mapstructure:"ACCESS_JWT_KEY"`
	RefreshJwtKey         string  `mapstructure:"REFRESH_JWT_KEY"`
	ApiPrefixStr          string  `mapstructure:"API_V1_PREFIX_STRING"`
	Host                  string  `mapstructure:"HOST"`
	RedisAddress          string  `mapstructure:"REDIS_ADDRESS"`
	RedisPassword         string  `mapstructure:"REDIS_PASSWORD"`
	RedisUsername         string  `mapstructure:"REDIS_USERNAME"`
	CoudinaryURL          string  `mapstructure:"CLOUDINARY_URL"`
	RefreshExpirationHour float64 `mapstructure:"REFRESH_EXPIRATION_HOUR"`
	RedisDB               int     `mapstructure:"REDIS_DB"`
	AccessExpirationHour  float64 `mapstructure:"ACCESS_EXPIRATION_HOUR"`
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
