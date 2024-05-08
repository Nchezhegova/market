package db

import (
	"context"
	"github.com/Nchezhegova/market/internal/log"
	"go.uber.org/zap"
)

func CheckUser(ctx context.Context, name string) (bool, error) {
	var count int
	err := DB.QueryRowContext(ctx, "SELECT COUNT(*) FROM users WHERE name = $1",
		name).Scan(&count)
	if err != nil {
		log.Logger.Info("Problem with checking user", zap.Error(err))
		return false, err
	}
	return count > 0, nil
}

func AddUser(ctx context.Context, name string, password string) error {
	_, err := DB.ExecContext(ctx, "INSERT INTO users (name,password) VALUES ($1, $2)",
		name, password)
	if err != nil {
		log.Logger.Info("Problem with adding user", zap.Error(err))
		return err
	}
	return nil
}

func CheckPassword(ctx context.Context, name string) (string, int, error) {
	var p string
	var id int
	err := DB.QueryRowContext(ctx, "SELECT password, id FROM users WHERE name = $1",
		name).Scan(&p, &id)
	if err != nil {
		log.Logger.Info("Problem with checking password", zap.Error(err))
		return "", 0, err
	}
	return p, id, nil
}
