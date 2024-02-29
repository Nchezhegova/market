package db

import (
	"context"
	"github.com/Nchezhegova/market/internal/log"
	"go.uber.org/zap"
)

func CheckUser(ctx context.Context, name string) bool {
	var count int
	err := DB.QueryRowContext(ctx, "SELECT COUNT(*) FROM users WHERE name = $1",
		name).Scan(&count)
	if err != nil {
		log.Logger.Fatal("Problem with checking user", zap.Error(err))
	}
	return count > 0
}

func AddUser(ctx context.Context, name string, password string) {
	_, err := DB.ExecContext(ctx, "INSERT INTO users (name,password) VALUES ($1, $2)",
		name, password)
	if err != nil {
		log.Logger.Fatal("Problem with adding user", zap.Error(err))
	}
}

func CheckPassword(ctx context.Context, name string) (string, int) {
	var p string
	var id int
	err := DB.QueryRowContext(ctx, "SELECT password, id FROM users WHERE name = $1",
		name).Scan(&p, &id)
	if err != nil {
		log.Logger.Fatal("Problem with checking password", zap.Error(err))
	}
	return p, id
}
