package main

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/omniful/go_commons/http"
	"github.com/omniful/go_commons/log"
	"github.com/singhJasvinder101/go_wms/internal/config"
	"github.com/singhJasvinder101/go_wms/internal/handlers"
	"github.com/singhJasvinder101/go_wms/internal/services"
	"github.com/singhJasvinder101/go_wms/internal/setup"
	"github.com/singhJasvinder101/go_wms/internal/storage"
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


	server := http.InitializeServer(
		":3001", 0, 0, 0, true,
	)

	server.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK.Code(), gin.H{
			"message": "service is healthy",
			"data":    map[string]string{"status": "healthy"},
		})
	})


	//databse initialization
	cluster := storage.NewPostgres(ctx)
	log.InfofWithContext(ctx, "database initialized successfully %v", cluster)

	// repos
	hubRepo := storage.NewHubRepo(cluster)
	skuRepo := storage.NewSKURepo(cluster)
	inventoryRepo := storage.NewInventoryRepo(cluster)

	//services
	hubService := services.NewHubService(hubRepo)
	skuService := services.NewSKUService(skuRepo)
	inventoryService := services.NewInventoryService(inventoryRepo, skuRepo, hubRepo)

	//handlers
	hubHandler := handlers.NewHubHandler(hubService)
	skuHandler := handlers.NewSKUHandler(skuService)
	inventoryHandler := handlers.NewInventoryHandler(inventoryService)

	setup.SetupRoutes(server, hubHandler, skuHandler, inventoryHandler)

	log.InfofWithContext(ctx, "starting server on port 3001")
	if err := server.StartServer("wms-service"); err != nil {
		log.Error("error while starting the server")
	}
}
