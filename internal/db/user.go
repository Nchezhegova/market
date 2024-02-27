package db

import (
	"context"
	"github.com/Nchezhegova/market/internal/log"
	"go.uber.org/zap"
)

func CheckUser(ctx context.Context, name string, email string) bool {
	var count int
	err := DB.QueryRowContext(ctx, "SELECT COUNT(*) FROM users WHERE name = $1 OR email = $2",
		name, email).Scan(&count)
	if err != nil {
		log.Logger.Fatal("Problem with checking user", zap.Error(err))
	}
	return count > 0
}

func AddUser(ctx context.Context, name string, email string, password string) {
	_, err := DB.ExecContext(ctx, "INSERT INTO users (name, email,password) VALUES ($1, $2, $3)",
		name, email, password)
	if err != nil {
		log.Logger.Fatal("Problem with adding user", zap.Error(err))
	}
}

func CheckPassword(ctx context.Context, email string) (string, int) {
	var p string
	var id int
	err := DB.QueryRowContext(ctx, "SELECT password, id FROM users WHERE email = $1",
		email).Scan(&p, &id)
	if err != nil {
		log.Logger.Fatal("Problem with checking password", zap.Error(err))
	}
	return p, id
}
