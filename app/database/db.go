package database

import (
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

var (
	ErrorFailedSQLReq = fmt.Errorf("failed to exec sql request")
)

type Database struct {
	psqlConnection *PostgresConnection
}

type PostgresConnection struct {
	SQLBuilder squirrel.StatementBuilderType
	Pool       *pgxpool.Pool
	maxConns   int
}

func NewPostgresConnection(dbURL string, maxConns int) (*PostgresConnection, error) {
	var psqlConnection PostgresConnection
	psqlConnection.maxConns = maxConns
	psqlConnection.SQLBuilder = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	config, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		return nil, err
	}
	config.MaxConns = int32(psqlConnection.maxConns)

	psqlConnection.Pool, err = pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}

	return &psqlConnection, nil
}

func NewDatabase(psql *PostgresConnection) *Database {
	return &Database{
		psqlConnection: psql,
	}
}

func (db *Database) AddUserReq(ctx context.Context, userID int64) error {
	sql, args, _ := db.psqlConnection.SQLBuilder.
		Insert("requests").Columns("user_id, req_time").Values(userID, time.Now()).ToSql()

	_, err := db.psqlConnection.Pool.Exec(ctx, sql, args)
	if err != nil {
		return ErrorFailedSQLReq
	}

	return nil
}

func (db *Database) CheckAmountRequestInN(ctx context.Context, userID int64, N int) (int, error) {
	sql, args, _ := db.psqlConnection.SQLBuilder.
		Select("COUNT(*)").From("requests").
		Where("user_id = ? AND req_time > (NOW() - interval '? seconds')", userID, N).ToSql()

	var amount int
	err := db.psqlConnection.Pool.QueryRow(ctx, sql, args).Scan(&amount)
	if err != nil {
		return 0, err
	}

	return amount, nil
}
