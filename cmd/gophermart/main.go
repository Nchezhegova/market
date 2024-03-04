package main

import (
	"context"
	"github.com/Nchezhegova/market/internal/config"
	"github.com/Nchezhegova/market/internal/db"
	"github.com/Nchezhegova/market/internal/http/server"
	"github.com/Nchezhegova/market/internal/log"
	"github.com/Nchezhegova/market/internal/service/accrual"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg := config.Config{}
	cfg.GenerationConfig()

	db.RunDB(cfg.Database)
	accrual.RunAccrual(ctx, cfg.Accrual)
	server.StartServer(cfg.Service)

	defer log.Logger.Sync()
	defer db.DB.Close()
}
