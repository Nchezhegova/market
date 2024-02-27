package models

import (
	"context"
	"github.com/Nchezhegova/market/internal/db"
)

//type AccrualModel struct {
//	Number  string  `json:"order"`
//	State   string  `json:"status"`
//	Accrual float64 `json:"accrual"`
//}

type BalanceModel struct {
	Sum        float64 `json:"current"`
	Withdrawal float64 `json:"withdrawn"`
}

type Balance interface {
	GetBalance(context.Context, int) error
}

func (b *BalanceModel) GetBalance(ctx context.Context, uid int) error {
	var err error
	var sum int
	sum, err = db.GetAccrual(ctx, uid)
	if err != nil {
		return err
	}

	var withdrawal int
	withdrawal, err = db.GetWithdrawal(ctx, uid)
	if err != nil {
		return err
	}
	b.Withdrawal = float64(withdrawal / 100)
	b.Sum = float64(sum/100) - b.Withdrawal
	return nil
}
