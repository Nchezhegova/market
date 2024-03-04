package db

import (
	"context"
	"github.com/Nchezhegova/market/internal/log"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

type WithdrawalDB struct {
	Order     string          `json:"order"`
	Sum       decimal.Decimal `json:"sum"`
	Processed string          `json:"processed_at,omitempty"`
}

func GetAccrual(ctx context.Context, uid int) (decimal.Decimal, error) {
	var sum decimal.Decimal
	err := DB.QueryRowContext(ctx, "SELECT COALESCE(SUM(accrual),0) FROM orders WHERE user_id = $1",
		uid).Scan(&sum)
	if err != nil {
		log.Logger.Info("Problem with getting accrual", zap.Error(err))
		return decimal.New(0, 0), err
	}
	return sum, nil
}

func GetWithdrawal(ctx context.Context, uid int) (decimal.Decimal, error) {
	var sum decimal.Decimal
	err := DB.QueryRowContext(ctx, "SELECT COALESCE(SUM(withdrawal),0) FROM withdrawals WHERE user_id = $1",
		uid).Scan(&sum)
	if err != nil {
		log.Logger.Info("Problem with getting withdrawal", zap.Error(err))
		return decimal.New(0, 0), err
	}
	return sum, nil
}

func AddWithdrawal(ctx context.Context, uid int, order string, w decimal.Decimal, p string) error {
	_, err := DB.ExecContext(ctx, "INSERT INTO withdrawals (order_id, user_id, withdrawal,processed_at) VALUES ($1, $2, $3, $4)",
		order, uid, w, p)
	if err != nil {
		log.Logger.Info("Problem with adding withdrawal", zap.Error(err))
		return err
	}
	return nil
}

func GetWithdrawals(ctx context.Context, uid int) ([]WithdrawalDB, error) {
	var WithdrawalList []WithdrawalDB
	rows, err := DB.QueryContext(ctx, "SELECT order_id,COALESCE(withdrawal,0),processed_at FROM withdrawals WHERE user_id=$1", uid)
	if err != nil {
		log.Logger.Info("Error DB:", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	if err := rows.Err(); err != nil {
		log.Logger.Info("Error DB:", zap.Error(err))
	}
	for rows.Next() {
		var w WithdrawalDB
		if err := rows.Scan(&w.Order, &w.Sum, &w.Processed); err != nil {
			log.Logger.Info("Error DB:", zap.Error(err))
			return nil, err
		}
		WithdrawalList = append(WithdrawalList, w)
	}

	return WithdrawalList, nil
}
