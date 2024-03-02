package models

import (
	"context"
	"fmt"
	"github.com/Nchezhegova/market/internal/db"
	"github.com/Nchezhegova/market/internal/log"
	"github.com/Nchezhegova/market/internal/service/luhn"
	"github.com/shopspring/decimal"
	"strconv"
	"time"
)

type OrderModel struct {
	ID      int             `json:"ID"`
	Number  int             `json:"order"`
	UserID  int             `json:"user_id"`
	State   string          `json:"status"`
	Accrual decimal.Decimal `json:"accrual"`
	Upload  string          `json:"uploaded_at"`
}

type OrderWithdrawal struct {
	Number  string  `json:"number"`
	Status  string  `json:"status"`
	Accrual float64 `json:"accrual,omitempty"`
	Upload  string  `json:"uploaded_at"`
}

type Order interface {
	CheckNumber(context.Context, string) error
	AddOrder(context.Context, string) error
	CheckOrder(context.Context, int) bool
}

func (o *OrderModel) CheckNumber(ctx context.Context, number string) error {
	if !luhn.Luhn(number) {
		err := fmt.Errorf("not valid number")
		return err
	}
	return nil
}

func (o *OrderModel) AddOrder(ctx context.Context, number string, uid int) error {
	var err error
	if err = o.CheckNumber(ctx, number); err != nil {
		return err
	}
	o.Number, err = strconv.Atoi(number)
	if err != nil {
		return err
	}

	if o.CheckOrder(ctx, uid) {
		err = fmt.Errorf("Order already exists")
		return err
	}
	o.State = "NEW"
	o.Upload = time.Now().Format(time.RFC3339)
	db.AddOrder(ctx, o.Number, o.State, uid, o.Upload)
	return nil
}

func (o *OrderModel) CheckOrder(ctx context.Context, uid int) bool {
	if exists := db.CheckOrder(ctx, o.Number, uid); exists {
		log.Logger.Info("Order already exists")
		return exists
	}
	return false
}

func UpdateOrder(ctx context.Context, number string, status string, accrual decimal.Decimal) {
	var o OrderModel
	var err error
	o.Number, err = strconv.Atoi(number)
	if err != nil {
		return
	}
	o.State = status
	o.Accrual = accrual
	db.UpdateOrder(ctx, o.Number, o.State, o.Accrual)
}

func GetOrders(ctx context.Context, uid int) []OrderWithdrawal {
	var DBorders []db.DBOrder
	var o []OrderWithdrawal
	DBorders = db.GetOrders(ctx, uid)
	for i := range DBorders {
		var order OrderWithdrawal
		order.Number = strconv.Itoa(DBorders[i].Number)
		order.Status = DBorders[i].Status
		if order.Status == "PROCESSED" {
			order.Accrual, _ = DBorders[i].Accrual.Float64()
		}
		order.Upload = DBorders[i].Upload
		o = append(o, order)
	}
	return o
}
