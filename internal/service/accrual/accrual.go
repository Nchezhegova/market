package accrual

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/Nchezhegova/market/internal/config"
	"github.com/Nchezhegova/market/internal/db"
	"github.com/Nchezhegova/market/internal/log"
	"github.com/Nchezhegova/market/internal/models"
	"go.uber.org/zap"
	"net/http"
	"time"
)

var OrdersList chan int

type AccrualModel struct {
	Number  string  `json:"order"`
	State   string  `json:"status"`
	Accrual float64 `json:"accrual"`
}

func RunAccrual(ctx context.Context, addr string) {
	OrdersList = make(chan int)
	go func() {
		for {
			GenerateOrdersList(ctx)
			time.Sleep(time.Duration(10) * time.Second)
		}
	}()
	for i := 0; i < config.ELEMENTS; i++ {
		go Worker(ctx, addr)
	}

}

func Worker(ctx context.Context, addr string) {
	for {
		order, ok := <-OrdersList
		if !ok {
			log.Logger.Info("Problem with channel")
			return
		}
		GetOrderInformation(ctx, order, addr)
	}
}

func GetOrderInformation(ctx context.Context, number int, addr string) {
	url := fmt.Sprintf("http://%s/api/orders/%v", addr, number)
	resp, err := http.Get(url)
	if err != nil {
		log.Logger.Info("Error get information from accrual", zap.Error(err))
	}
	if resp.StatusCode != http.StatusOK {
		log.Logger.Info("Response status not OK", zap.Error(err))
		return
	}

	var order AccrualModel
	var buf bytes.Buffer
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		log.Logger.Info("Error reading information from accrual", zap.Error(err))
		return
	}
	if err := json.Unmarshal(buf.Bytes(), &order); err != nil {
		log.Logger.Info("Error convert information from accrual", zap.Error(err))
		return
	}
	order.UpdateOrderInformation(ctx)
}

func GenerateOrdersList(ctx context.Context) {
	number, user := db.GetNewOrder(ctx)
	if number == 0 {
		return
	}
	OrdersList <- number
	db.OrderProcessing(ctx, number, user)
}

func (a *AccrualModel) UpdateOrderInformation(ctx context.Context) {
	models.UpdateOrder(ctx, a.Number, a.State, a.Accrual)
}
