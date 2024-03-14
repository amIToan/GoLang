package util

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Environment              string        `mapstructure:"ENVIRONMENT"`
	DbDriver                 string        `mapstructure:"DB_DRIVER"`
	DbSource                 string        `mapstructure:"DB_SOURCE"`
	MigrationUrl             string        `mapstructure:"MIGRATION_URL"`
	HttpServerAddress        string        `mapstructure:"HTTP_SERVER_ADDRESS"`
	GrpcServerAddress        string        `mapstructure:"GRPC_SERVER_ADDRESS"`
	TokenSymmetricKey        string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	RedisAddress             string        `mapstructure:"REDIS_ADDRESS"`
	EmailSenderAddress       string        `mapstructure:"EMAIL_SENDER_ADDRESS"`
	EmailSenderPassword      string        `mapstructure:"EMAIL_SENDER_PASSWORD"`
	EmailSenderName          string        `mapstructure:"EMAIL_SENDER_NAME"`
	ValidDurationTime        time.Duration `mapstructure:"VALID_TOKEN_DURATION"`
	RefreshTokenDurationTime time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
}

func LoadConfigDB_Server(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AutomaticEnv()
	err = viper.ReadInConfig()
	if err != nil {
		return
	}
	err = viper.Unmarshal(&config)
	return
}
