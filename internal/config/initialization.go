package config

import (
	"flag"
	"os"
)

type Config struct {
	Service  string
	Database string
	Accrual  string
}

func (c *Config) GenerationConfig() {
	flag.StringVar(&c.Service, "a", ADDRSERV, "input addr serv")
	flag.StringVar(&c.Database, "d", DATEBASE, "input addr DB")
	flag.StringVar(&c.Accrual, "r", ACCRUALSYSTEMADDRESS, "input addr accrual")
	flag.Parse()

	if envRunAddr := os.Getenv("RUN_ADDRESS"); envRunAddr != "" {
		c.Service = envRunAddr
	}
	if envAddrDB := os.Getenv("DATABASE_URI"); envAddrDB != "" {
		c.Database = envAddrDB
	}
	if envAccrual := os.Getenv("ACCRUAL_SYSTEM_ADDRESS"); envAccrual != "" {
		c.Accrual = envAccrual
	}

}
