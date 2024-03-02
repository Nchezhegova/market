package db

import (
	"context"
	"github.com/Nchezhegova/market/internal/log"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

type DBOrder struct {
	Number  int             `json:"number"`
	Status  string          `json:"status"`
	Accrual decimal.Decimal `json:"accrual,omitempty"`
	Upload  string          `json:"uploaded_at"`
}

func AddOrder(ctx context.Context, onumber int, ostate string, uid int, upload string) {
	_, err := DB.ExecContext(ctx, "INSERT INTO orders (number, status, user_id, uploaded_at) VALUES ($1, $2, $3, $4)",
		onumber, ostate, uid, upload)
	if err != nil {
		log.Logger.Fatal("Problem with adding order", zap.Error(err))
	}
}

func CheckOrder(ctx context.Context, onumber int) int {
	var uid int
	err := DB.QueryRowContext(ctx, "SELECT COALESCE(user_id,0) FROM orders WHERE number = $1",
		onumber).Scan(&uid)
	if err != nil {
		log.Logger.Info("Problem with checking order", zap.Error(err))
	}
	return uid
}

func GetNewOrder(ctx context.Context) (int, int) {
	var number int
	var user int
	err := DB.QueryRowContext(ctx, "SELECT number,user_id FROM orders WHERE status = $1 OR status =$2",
		"NEW", "PROCESSING").Scan(&number, &user)
	if err != nil {
		log.Logger.Info("Problem with getting order", zap.Error(err))
		return 0, 0
	}
	return number, user
}

func OrderProcessing(ctx context.Context, number int, user int) {
	_, err := DB.ExecContext(ctx, "UPDATE orders SET status =$1 WHERE number = $2 AND user_id = $3",
		"PROCESSING", number, user)
	if err != nil {
		log.Logger.Info("Problem with update order", zap.Error(err))
	}
}

func UpdateOrder(ctx context.Context, number int, status string, accrual decimal.Decimal) {
	_, err := DB.ExecContext(ctx, "UPDATE orders SET status =$1, accrual =$2 WHERE number = $3",
		status, accrual, number)
	if err != nil {
		log.Logger.Info("Problem with update order", zap.Error(err))
	}
}

func GetOrders(ctx context.Context, uid int) []DBOrder {
	var DBorders []DBOrder
	rows, err := DB.QueryContext(ctx, "SELECT number,status,COALESCE(accrual,0),uploaded_at FROM orders WHERE user_id=$1", uid)
	if err != nil {
		log.Logger.Info("Error DB:", zap.Error(err))
	}
	defer rows.Close()

	for rows.Next() {
		var d DBOrder
		if err := rows.Scan(&d.Number, &d.Status, &d.Accrual, &d.Upload); err != nil {
			log.Logger.Info("Error DB:", zap.Error(err))
		}
		DBorders = append(DBorders, d)
	}
	return DBorders
}
