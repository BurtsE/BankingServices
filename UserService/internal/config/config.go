package config

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	ServerPort      string `json:"http_port" yaml:"http_port"`
	GrpcPort        int    `json:"grpc_port" yaml:"grpc_port"`
	JWTSecret       string `json:"jwt_secret" yaml:"jwt_secret"`
	LogLevel        string `json:"log_level" yaml:"log_level"`
	Postgres        `json:"postgres" yaml:"postgres"`
	PostgresTracing `json:"postgres_tracing" yaml:"postgres_tracing"`
}

type Postgres struct {
	Host     string `json:"host" yaml:"host"`
	Port     string `json:"port" yaml:"port"`
	Username string `json:"username" yaml:"username"`
	Password string `json:"password" yaml:"password"`
	Database string `json:"database" yaml:"database"`
}

type PostgresTracing struct {
	EnableTracing            bool `json:"enable_tracing" yaml:"enable_tracing"`
	IncludeConnectionDetails bool `json:"include_connection_details" yaml:"include_connection_details"`
	IncludeSqlStatements     bool `json:"include_sql_statements" yaml:"include_sql_statements"`
	IncludeQueryPrefix       bool `json:"include_query_prefix" yaml:"include_query_prefix"`
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

func GetJWTSecretKey() string {
	return getEnv("JWT_SECRET_KEY")
}

func GetEnv() string {
	return getEnv("ENV")
}

func GetJaegerUrl() string {
	return getEnv("JAEGER_URL")
}

func getEnv(param string) string {
	val := os.Getenv(param)
	if val == "" {
		log.Fatalf("env not set: %s", param)
	}
	return val
}
