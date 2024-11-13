package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	RedisAddress          string  `mapstructure:"REDIS_ADDRESS"`
	AccessJwtKey          string  `mapstructure:"ACCESS_JWT_KEY"`
	Environment           string  `mapstructure:"ENVIROMENT"`
	DbType                string  `mapstructure:"DB_TYPE"`
	RedisPassword         string  `mapstructure:"REDIS_PASSWORD"`
	RefreshJwtKey         string  `mapstructure:"REFRESH_JWT_KEY"`
	ApiPrefixStr          string  `mapstructure:"API_V1_PREFIX_STRING"`
	RedisUsername         string  `mapstructure:"REDIS_USERNAME"`
	DbURL                 string  `mapstructure:"DB_URL"`
	HttpAddress           string  `mapstructure:"HTTP_SERVER_ADDRESS"`
	Host                  string  `mapstructure:"HOST"`
	CoudinaryURL          string  `mapstructure:"CLOUDINARY_URL"`
	FrontendURL           string  `mapstructure:"FRONTEND_URL"`
	GoogleClientSecret    string  `mapstructure:"GOOGLE_CLIENT_SECRET"`
	GoogleClientID        string  `mapstructure:"GOOGLE_CLIENT_ID"`
	GoogleSigningKey      string  `mapstructure:"GOOGLE_SIGNING_KEY"`
	GoogleMaxAge          int     `mapstructure:"GOOGLE_MAX_AGE"`
	AccessExpirationHour  float64 `mapstructure:"ACCESS_EXPIRATION_HOUR"`
	RedisDB               int     `mapstructure:"REDIS_DB"`
	RefreshExpirationHour float64 `mapstructure:"REFRESH_EXPIRATION_HOUR"`
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
