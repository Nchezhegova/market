package main

import (
	"context"
	"github.com/Nchezhegova/market/internal/db"
	"github.com/Nchezhegova/market/internal/http/server"
	"github.com/Nchezhegova/market/internal/log"
	"github.com/Nchezhegova/market/internal/service/accrual"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	accrual.RunAccrual(ctx)

	server.StartServer()

	defer log.Logger.Sync()
	defer db.DB.Close()
}
