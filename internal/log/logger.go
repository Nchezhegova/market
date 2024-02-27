package log

import (
	"fmt"
	"go.uber.org/zap"
)

var Logger *zap.Logger

func init() {
	var err error
	Logger, err = zap.NewProduction()
	if err != nil {
		panic(fmt.Sprintf("Не удалось создать логгер: %v", err))
	}

}
