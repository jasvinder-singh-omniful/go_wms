package utils

import (
	"context"

	"github.com/omniful/go_commons/log"
	"github.com/singhJasvinder101/go_wms/internal/config"
)

func InitLogger(ctx context.Context){
	dev_env := config.GetConfig().Environment == "local"

	logLevel := "info"
	logFormat := "json"


	if dev_env {
		logLevel = "debug"
		logFormat = "text"
	}

	err := log.InitializeLogger(
		log.Formatter(logFormat),
		log.Level(logLevel),
		log.ColoredLevelEncoder(),
	)

	if err != nil {
		log.Error("logger init error")
		panic(err)
	}

	log.InfofWithContext(ctx, "logger initialized successfully")
}
