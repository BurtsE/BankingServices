package config

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	ServerPort string `json:"port" yaml:"port"`
	LogLevel   string `json:"log_level" yaml:"log_level"`
	Postgres   `json:"postgres" yaml:"postgres"`
}

type Postgres struct {
	Host     string `json:"host" yaml:"host"`
	Port     string `json:"port" yaml:"port"`
	Username string `json:"username" yaml:"username"`
	Password string `json:"password" yaml:"password"`
	Database string `json:"database" yaml:"database"`
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

func GetCardSecretKey() string {
	return getEnv("CARD_SECRET_KEY")
}

func GetEncryptionPublicKey() string {
	return getEnv("ENCRYPTION_PUBLIC_KEY")
}
func GetEncryptionPrivateKey() string {
	return getEnv("ENCRYPTION_PRIVATE_KEY")
}

func GetBankingServiceGrpcURI() string {
	return getEnv("BANKING_SERVICE_GRPC_URI")
}

func getEnv(param string) string {
	val := os.Getenv(param)
	if val == "" {
		log.Fatalf("env not set: %s", param)
	}
	return val
}
