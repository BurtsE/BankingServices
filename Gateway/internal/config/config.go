package config

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	ServerPort string `json:"port" yaml:"port"`
	LogLevel   string `json:"log_level" yaml:"log_level"`
	Redis      `json:"redis" yaml:"redis"`
}

type Redis struct {
	Host     string `json:"host" yaml:"host"`
	Port     string `json:"port" yaml:"port"`
	Password string `json:"password" yaml:"password"`
	Database int    `json:"database" yaml:"database"`
}

func InitConfig() (*Config, error) {
	c := &Config{}
	data, err := os.ReadFile("./configs/app_config.json")
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, c)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func GetUserServiceHttpURI() string {
	return getEnv("USER_SERVICE_HTTP_URI")
}

func GetUserServiceGrpcURI() string {
	return getEnv("USER_SERVICE_GRPC_URI")
}

func GetBankingServiceURI() string {
	return getEnv("BANKING_SERVICE_URI")
}

func getEnv(param string) string {
	val := os.Getenv(param)
	if val == "" {
		log.Fatalf("env not set: %s", param)
	}
	return val
}
