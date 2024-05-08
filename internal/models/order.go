package models

import (
	"context"
	"fmt"
	"github.com/Nchezhegova/market/internal/db"
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
	var err error
	o.Number, err = strconv.Atoi(number)
	if err != nil {
		return err
	}
	return nil
}

func (o *OrderModel) AddOrder(ctx context.Context, uid int) error {
	o.State = "NEW"
	o.Upload = time.Now().Format(time.RFC3339)
	if err := db.AddOrder(ctx, o.Number, o.State, uid, o.Upload); err != nil {
		return err
	}
	return nil
}

func (o *OrderModel) CheckOrder(ctx context.Context) (int, error) {
	uid, err := db.CheckOrder(ctx, o.Number)
	return uid, err
}

func UpdateOrder(ctx context.Context, number string, status string, accrual decimal.Decimal) error {
	var o OrderModel
	var err error
	o.Number, err = strconv.Atoi(number)
	if err != nil {
		return err
	}
	o.State = status
	o.Accrual = accrual
	if err = db.UpdateOrder(ctx, o.Number, o.State, o.Accrual); err != nil {
		return err
	}
	return nil
}

func GetOrders(ctx context.Context, uid int) ([]OrderWithdrawal, error) {
	var DBorders []db.DBOrder
	var o []OrderWithdrawal
	var err error
	DBorders, err = db.GetOrders(ctx, uid)
	if err != nil {
		return nil, err
	}
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
	return o, nil
}
