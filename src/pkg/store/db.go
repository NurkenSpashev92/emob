package store

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/nurkenspashev92/emob/configs"
)

type Database struct {
	Conn *pgx.Conn
}

var dbInstance *Database

func NewPostgresDb(conf *configs.Config) (*Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if dbInstance != nil {
		return dbInstance, nil
	}

	conn, err := pgx.Connect(ctx, conf.DatabaseURL())
	if err != nil {
		return nil, fmt.Errorf("failed to connect db: %w", err)
	}

	if err := conn.Ping(ctx); err != nil {
		conn.Close(ctx)
		return nil, fmt.Errorf("failed to ping db: %w", err)
	}

	dbInstance = &Database{
		Conn: conn,
	}

	return dbInstance, nil
}
