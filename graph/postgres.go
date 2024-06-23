package graph

import (
	"context"
	"database/sql"
	"fmt"
	"t-ozon/graph/model"

	"github.com/jackc/pgx/v4/pgxpool"
)

func ConnectToPostgres(config map[string]string) (*pgxpool.Pool, error) {
	host := config["host"]
	port := config["port"]
	user := config["user"]
	password := config["password"]
	database := config["database"]
	connectionString := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s database=%s",
		host, port, user, password, database,
	)
	pool, err := pgxpool.Connect(context.Background(), connectionString)
	if err != nil {
		return nil, err
	}
	return pool, nil
}

func (db *PostgresDB) GetUserByIDBD(ctx context.Context, userID string) (*model.User, error) {
	var user model.User
	err := db.pool.QueryRow(ctx, "SELECT id, name FROM public.\"User\" WHERE id = $1", userID).
		Scan(&user.ID, &user.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found: %w", err)
		}
		return nil, fmt.Errorf("failed to fetch user from database: %w", err)
	}
	return &user, nil
}
