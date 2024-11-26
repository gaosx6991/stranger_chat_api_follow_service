package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Server ServerConfig `mapstructure:"server"`

	GrpcServer ServerConfig `mapstructure:"grpc_server"`

	MongoDB     MongoDBConfig `mapstructure:"mongodb"`
	UserService ServiceConfig `mapstructure:"user_service"`
	PostService ServiceConfig `mapstructure:"post_service"`
}

type ServerConfig struct {
	Port int `mapstructure:"port"`
}

type MongoDBConfig struct {
	URI        string `mapstructure:"uri"`
	Database   string `mapstructure:"database"`
	Collection string `mapstructure:"collection"`
}

type ServiceConfig struct {
	Host string `mapstructure:"host"`
}

func LoadConfig(path string) (*Config, error) {
	viper.SetConfigFile(path)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	config := &Config{}
	if err := viper.Unmarshal(config); err != nil {
		return nil, err
	}

	return config, nil
}
