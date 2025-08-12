package clickhouse

import (
	"context"
	"crypto/tls"
	"fmt"
	"gateway/internal/config"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

type ClickHouseStorage struct {
	conn driver.Conn
}

func NewClickHouseStorage(cfg *config.Config) *ClickHouseStorage {
	conn, err := connect(cfg)
	if err != nil {
		panic(err)
	}
	return &ClickHouseStorage{
		conn: conn,
	}
}

func connect(cfg *config.Config) (driver.Conn, error) {
	var (
		ctx       = context.Background()
		conn, err = clickhouse.Open(&clickhouse.Options{
			Addr: []string{cfg.ClickHouse.Address},
			Auth: clickhouse.Auth{
				Database: cfg.ClickHouse.Database,
				Username: cfg.ClickHouse.Username,
				Password: cfg.ClickHouse.Password,
			},
			ClientInfo: clickhouse.ClientInfo{
				Products: []struct {
					Name    string
					Version string
				}{
					{Name: "an-example-go-client", Version: "0.1"},
				},
			},
			Debugf: func(format string, v ...any) {
				fmt.Printf(format, v)
			},
			TLS: &tls.Config{
				InsecureSkipVerify: true,
			},
		})
	)

	if err != nil {
		return nil, err
	}

	if err := conn.Ping(ctx); err != nil {
		if exception, ok := err.(*clickhouse.Exception); ok {
			fmt.Printf("Exception [%d] %s \n%s\n", exception.Code, exception.Message, exception.StackTrace)
		}
		return nil, err
	}
	return conn, nil
}
