package main

import (
	"context"

	"github.com/omniful/go_commons/log"
	"github.com/singhJasvinder101/go_wms/internal/config"
	"github.com/singhJasvinder101/go_wms/utils"
)


func main(){
	ctx := context.Background()

	//initialize configs
	config.InitConfig(ctx)
	
	cfg := config.GetConfig()
	log.InfofWithContext(ctx, "config initialized succesffully %v", cfg)


	//initialize logger
	utils.InitLogger(ctx)



}
