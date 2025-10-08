package setup

import (
	"github.com/omniful/go_commons/http"
	"github.com/singhJasvinder101/go_wms/internal/handlers"
)

func SetupRoutes(server *http.Server, hubHandler *handlers.HubHandler, skuHandler *handlers.SKUHandler, inventoryHandler *handlers.InventoryHandler){
	v1 := server.Group("/api/v1")
	{
		//hub routes
		hubRoutes := v1.Group("/hubs")
		{
			hubRoutes.POST("/create", hubHandler.CreateHub)
			hubRoutes.POST("/get", hubHandler.GetHub)
			hubRoutes.GET("/getall", hubHandler.GetAllHubs)
		}
		

		//sku routes
		skuRoutes := v1.Group("/skus")
		{
			skuRoutes.POST("/create", skuHandler.CreateSKU)
			skuRoutes.POST("/get", skuHandler.GetSKUsByCodes)
		}

		//inventory routes
		inventoryRoutes := v1.Group("/inventory")
		{
			inventoryRoutes.POST("/create", inventoryHandler.CreateInventory)
			inventoryRoutes.PATCH("/upsert", inventoryHandler.UpsertInventory)
			inventoryRoutes.POST("/get", inventoryHandler.GetInventory)
			inventoryRoutes.PATCH("/update-quantity", inventoryHandler.UpdateInventoryQuantity)
		}
	}
}

