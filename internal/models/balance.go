package models

import (
	"context"
	"github.com/Nchezhegova/market/internal/db"
	"github.com/shopspring/decimal"
)

//type AccrualModel struct {
//	Number  string  `json:"order"`
//	State   string  `json:"status"`
//	Accrual float64 `json:"accrual"`
//}

type BalanceModel struct {
	Sum        decimal.Decimal `json:"current"`
	Withdrawal decimal.Decimal `json:"withdrawn"`
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
	b.Withdrawal = withdrawal
	b.Sum = sum.Sub(b.Withdrawal)
	return nil
}
