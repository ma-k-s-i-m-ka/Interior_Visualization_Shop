package postgres

import (
	"Interior_Visualization_Shop/app/pkg/config"
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"time"
)

func ConnectDB(cfg config.Config) (*pgx.Conn, error) {

	pgxConfig, err := pgx.ParseConfig(cfg.PostgreSQL.DSN)
	if err != nil {
		return nil, fmt.Errorf("cannot parse database config from dsn %v", err)
	}

	dbTimeout, dbCancel := context.WithTimeout(context.Background(), time.Duration(cfg.PostgreSQL.ConnectionTimeout)*time.Second)
	defer dbCancel()

	dbConn, err := pgx.ConnectConfig(dbTimeout, pgxConfig)
	if err != nil {
		return nil, fmt.Errorf("cannot connect to database: %v", err)
	}

	if err = dbConn.Ping(dbTimeout); err != nil {
		fmt.Errorf("cannot ping database: %v", err)
	}
	return dbConn, nil
}
