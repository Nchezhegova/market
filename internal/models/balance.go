package models

import (
	"context"
	"github.com/Nchezhegova/market/internal/db"
	"github.com/shopspring/decimal"
)

type BalanceModel struct {
	Sum        float64 `json:"current"`
	Withdrawal float64 `json:"withdrawn"`
}

type Balance interface {
	GetBalance(context.Context, int) error
}

func (b *BalanceModel) GetBalance(ctx context.Context, uid int) error {
	var err error
	var sum decimal.Decimal
	sum, err = db.GetAccrual(ctx, uid)
	if err != nil {
		return err
	}

	var withdrawal decimal.Decimal
	withdrawal, err = db.GetWithdrawal(ctx, uid)
	if err != nil {
		return err
	}
	b.Withdrawal, _ = withdrawal.Float64()
	b.Sum, _ = sum.Sub(withdrawal).Float64()
	return nil
}
