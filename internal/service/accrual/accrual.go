package accrual

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Nchezhegova/market/internal/config"
	"github.com/Nchezhegova/market/internal/db"
	"github.com/Nchezhegova/market/internal/log"
	"github.com/Nchezhegova/market/internal/models"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
	"io"
	"net/http"
	"strconv"
	"time"
)

var OrdersList chan int

type AccrualModel struct {
	Number  string          `json:"order"`
	State   string          `json:"status"`
	Accrual decimal.Decimal `json:"accrual"`
}

var TooManyRequests = errors.New("too many requests for accrual service")

func RunAccrual(ctx context.Context, addr string) {
	OrdersList = make(chan int)
	go func() {
		for {
			if err := GenerateOrdersList(ctx); err != nil {
				log.Logger.Info("Problem with generate order list")
			}
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
		retry, err := GetOrderInformation(ctx, order, addr)
		if errors.Is(err, TooManyRequests) {
			time.Sleep(time.Duration(retry) * time.Second)
		}
	}
}

func GetOrderInformation(ctx context.Context, number int, addr string) (int, error) {
	url := fmt.Sprintf("%s/api/orders/%v", addr, number)
	resp, err := http.Get(url)
	if err != nil {
		log.Logger.Info("Error get information from accrual", zap.Error(err))
		return 0, err
	}
	if resp.StatusCode == http.StatusTooManyRequests {
		retry := resp.Header.Get("Retry-After")
		retryInt, err := strconv.Atoi(retry)
		if err != nil {
			log.Logger.Info("Can't convert Retry-After", zap.Error(err))
			return 0, err
		}
		return retryInt, TooManyRequests
	}
	if resp.StatusCode != http.StatusOK {
		log.Logger.Info("Response status not OK", zap.Error(err))
		return 0, err
	}
	var order AccrualModel
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Logger.Info("Error reading information from accrual", zap.Error(err))
		return 0, err
	}
	if err = json.Unmarshal(b, &order); err != nil {
		log.Logger.Info("Error convert information from accrual", zap.Error(err))
		return 0, err
	}
	if err = order.UpdateOrderInformation(ctx); err != nil {
		log.Logger.Info("Error update order", zap.Error(err))
		return 0, err
	}
	err = resp.Body.Close()
	if err != nil {
		log.Logger.Info("Error closing body:", zap.Error(err))
		return 0, err
	}
	return 0, nil
}

func GenerateOrdersList(ctx context.Context) error {
	number, user, err := db.GetNewOrder(ctx)
	if number == 0 {
		return err
	}
	OrdersList <- number
	if err = db.OrderProcessing(ctx, number, user); err != nil {
		return err
	}
	return nil
}

func (a *AccrualModel) UpdateOrderInformation(ctx context.Context) error {
	err := models.UpdateOrder(ctx, a.Number, a.State, a.Accrual)
	if err != nil {
		return err
	}
	return nil
}
