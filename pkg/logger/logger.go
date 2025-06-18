package logger

import (
	"log"

	"go.uber.org/zap"
)

var Log *zap.Logger

func Init() {
	var err error
	Log, err = zap.NewDevelopment()
	if err != nil {
		log.Fatalf("failed to init zap looger: %v", err)
	}
}
