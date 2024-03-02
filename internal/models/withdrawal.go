package models

import (
	"context"
	"fmt"
	"github.com/Nchezhegova/market/internal/db"
	"github.com/Nchezhegova/market/internal/service/luhn"
	"github.com/shopspring/decimal"
	"time"
)

type WithdrawalModel struct {
	Order     string          `json:"order"`
	Sum       decimal.Decimal `json:"sum"`
	Processed string          `json:"processed_at,omitempty"`
}

type WithdrawalResp struct {
	Order     string  `json:"order"`
	Sum       float64 `json:"sum"`
	Processed string  `json:"processed_at,omitempty"`
}

type Withdrawal interface {
	AddWithdrawal(context.Context, int) error
	CheckOrder(context.Context, int) bool
}

func (w *WithdrawalModel) CheckNumber(ctx context.Context, number string) error {
	if !luhn.Luhn(number) {
		err := fmt.Errorf("not valid number")
		return err
	}
	return nil
}

func (w *WithdrawalModel) AddWithdrawal(ctx context.Context, uid int) error {
	if err := w.CheckNumber(ctx, w.Order); err != nil {
		return err
	}

	w.Processed = time.Now().Format(time.RFC3339)
	if err := db.AddWithdrawal(ctx, uid, w.Order, w.Sum, w.Processed); err != nil {
		return err
	}
	return nil
}

func GetWithdrawal(ctx context.Context, uid int) (error, []WithdrawalResp) {
	var err error
	var w []WithdrawalResp
	var DBw []db.WithdrawalDB
	err, DBw = db.GetWithdrawals(ctx, uid)
	if err != nil {
		return err, nil
	}
	for i := range DBw {
		var witdrawal WithdrawalResp
		witdrawal.Order = DBw[i].Order
		witdrawal.Sum, _ = DBw[i].Sum.Float64()
		witdrawal.Processed = DBw[i].Processed

		w = append(w, witdrawal)
	}
	return nil, w
}
