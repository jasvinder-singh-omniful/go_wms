package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/omniful/go_commons/http"
	"github.com/omniful/go_commons/log"
	"github.com/omniful/go_commons/validator"
	"github.com/singhJasvinder101/go_wms/internal/services"
)

type InventoryHandler struct {
    InventoryService *services.InventoryService
}

func NewInventoryHandler(inventoryService *services.InventoryService) *InventoryHandler {
    return &InventoryHandler{
        InventoryService: inventoryService,
    }
}

type CreateInventoryRequest struct {
    TenantID string `json:"tenant_id" validate:"required"`
    SellerID string `json:"seller_id" validate:"required"`
    SKUCode  string `json:"sku_code" validate:"required,min=1"`
    HubID    int    `json:"hub_id" validate:"required,min=1"`
    Quantity int64  `json:"quantity" validate:"required,min=0"`
}

type UpsertInventoryRequest struct {
    TenantID string `json:"tenant_id" validate:"required"`
    SellerID string `json:"seller_id" validate:"required"`
    SKUCode  string `json:"sku_code" validate:"required,min=1"`
    HubID    int    `json:"hub_id" validate:"required,min=1"`
    Quantity int64  `json:"quantity" validate:"required,min=0"`
}

type GetInventoryRequest struct {
    TenantID string   `json:"tenant_id" validate:"required"`
    SellerID string   `json:"seller_id" validate:"required"`
    HubID    int      `json:"hub_id" validate:"required,min=1"`
    SKUCodes []string `json:"sku_codes,omitempty" validate:"omitempty,min=1,max=100,dive,required,min=1"`
}

type UpdateInventoryQuantityRequest struct {
    HubID    uint   `json:"hub_id" validate:"required,min=1"`
    SellerID string `json:"seller_id" validate:"required"`
    SKUCode  string `json:"sku_code" validate:"required,min=1"`
    Quantity int    `json:"quantity" validate:"required"`
}

type InventoryResponse struct {
    ID       uint   `json:"id"`
    SKUID    int    `json:"sku_id"`
    HubID    int    `json:"hub_id"`
    Quantity int64  `json:"quantity"`
}

type InventoryItem struct {
    SKU       string `json:"sku"`
    Inventory int64  `json:"inventory"`
}

func (h *InventoryHandler) CreateInventory(c *gin.Context) {
    ctx := c.Request.Context()
    logTag := "[InventoryHandler][CreateInventory]"
    log.InfofWithContext(ctx, logTag+" creating inventory")

	var body struct {
		TenantID string `json:"tenant_id" validate:"required"`
		SellerID string `json:"seller_id" validate:"required"`
		SKUCode  string `json:"sku_code" validate:"required,min=1"`
		HubID    int    `json:"hub_id" validate:"required,min=1"`
		Quantity int64  `json:"quantity" validate:"required,min=0"`
	}

    if err := c.ShouldBindJSON(&body); err != nil {
		log.ErrorfWithContext(ctx, logTag+" failed to bind JSON %v", err)
		c.JSON(http.StatusBadRequest.Code(), gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := validator.ValidateStruct(ctx, body); err.Exists() {
		log.ErrorfWithContext(ctx, logTag+" plaes enter vlaid input %v", err)
		c.JSON(http.StatusBadRequest.Code(), gin.H{
			"error":  err.ErrorMessage(),
			"errors": err.ErrorMap(),
		})
	}

    inventory, err := h.InventoryService.CreateInventory(ctx, body.TenantID, body.SellerID, body.SKUCode, body.HubID, body.Quantity)
    if err != nil {
        log.ErrorfWithContext(ctx, logTag+" failed to create inventory: %v", err)
        c.JSON(http.StatusInternalServerError.Code(), gin.H{
            "error": "Failed to create inventory",
        })
        return
    }

    response := &gin.H{
        "ID":       inventory.ID,
        "SKUID":    inventory.SKUID,
        "HubID":    inventory.HubID,
        "Quantity": inventory.Quantity,
    }

    log.InfofWithContext(ctx, logTag+" inventory created successfully")
    c.JSON(http.StatusCreated.Code(), response)
}

func (h *InventoryHandler) UpsertInventory(c *gin.Context) {
    ctx := c.Request.Context()

    logTag := "[InventoryHandler][UpsertInventory]"
    log.InfofWithContext(ctx, logTag+" upserting inventory")

    var body struct {
		TenantID string `json:"tenant_id" validate:"required"`
		SellerID string `json:"seller_id" validate:"required"`
		SKUCode  string `json:"sku_code" validate:"required,min=1"`
		HubID    int    `json:"hub_id" validate:"required,min=1"`
		Quantity int64  `json:"quantity" validate:"required,min=0"`
	}
    if err := c.ShouldBindJSON(&body); err != nil {
		log.ErrorfWithContext(ctx, logTag+" failed to bind JSON %v", err)
		c.JSON(http.StatusBadRequest.Code(), gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := validator.ValidateStruct(ctx, body); err.Exists() {
		log.ErrorfWithContext(ctx, logTag+" plaes enter vlaid input %v", err)
		c.JSON(http.StatusBadRequest.Code(), gin.H{
			"error":  err.ErrorMessage(),
			"errors": err.ErrorMap(),
		})
	}


    inventory, err := h.InventoryService.UpsertInventory(ctx, body.TenantID, body.SellerID, body.SKUCode, body.HubID, body.Quantity)
    if err != nil {
        log.ErrorfWithContext(ctx, logTag+" failed to upsert inventory %v", err)
        c.JSON(http.StatusInternalServerError.Code(), gin.H{
            "error": "Failed to upsert inventory",
        })
        return
    }

    response := &gin.H{
        "ID":       inventory.ID,
        "SKUID":    inventory.SKUID,
        "HubID":    inventory.HubID,
        "Quantity": inventory.Quantity,
    }

    log.InfofWithContext(ctx, logTag+" inventory upserted successfully")
    c.JSON(http.StatusOK.Code(), response)
}

func (h *InventoryHandler) GetInventory(c *gin.Context) {
    ctx := c.Request.Context()
    logTag := "[InventoryHandler][GetInventory]"
    log.InfofWithContext(ctx, logTag+" getting inventory")

    var body struct {
		TenantID string   `json:"tenant_id" validate:"required"`
		SellerID string   `json:"seller_id" validate:"required"`
		HubID    int      `json:"hub_id" validate:"required,min=1"`
		SKUCodes []string `json:"sku_codes,omitempty" validate:"omitempty,min=1,max=100,dive,required,min=1"`
	}

    if err := c.ShouldBindJSON(&body); err != nil {
		log.ErrorfWithContext(ctx, logTag+" failed to bind JSON %v", err)
		c.JSON(http.StatusBadRequest.Code(), gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := validator.ValidateStruct(ctx, body); err.Exists() {
		log.ErrorfWithContext(ctx, logTag+" plaes enter vlaid input %v", err)
		c.JSON(http.StatusBadRequest.Code(), gin.H{
			"error":  err.ErrorMessage(),
			"errors": err.ErrorMap(),
		})
	}

    inventoryList, err := h.InventoryService.InventoryRepo.GetByHubAndSeller(ctx, body.HubID, body.SellerID)
    if err != nil {
        log.ErrorfWithContext(ctx, logTag+" failed to get inventory: %v", err)
        c.JSON(http.StatusInternalServerError.Code(), gin.H{
            "error": "Failed to fetch inventory",
        })
        return
    }

    c.JSON(http.StatusOK.Code(), gin.H{
        "items": inventoryList,
        "count": len(inventoryList),
    })
}

func (h *InventoryHandler) UpdateInventoryQuantity(c *gin.Context) {
    ctx := c.Request.Context()
    logTag := "[InventoryHandler][UpdateInventoryQuantity]"
    log.InfofWithContext(ctx, logTag+" updating inventory quantity")

    var body struct {
		HubID    uint   `json:"hub_id" validate:"required,min=1"`
		SellerID string `json:"seller_id" validate:"required"`
		SKUCode  string `json:"sku_code" validate:"required,min=1"`
		Quantity int    `json:"quantity" validate:"required"`
	}
    if err := c.ShouldBindJSON(&body); err != nil {
		log.ErrorfWithContext(ctx, logTag+" failed to bind JSON %v", err)
		c.JSON(http.StatusBadRequest.Code(), gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := validator.ValidateStruct(ctx, body); err.Exists() {
		log.ErrorfWithContext(ctx, logTag+" plaes enter vlaid input %v", err)
		c.JSON(http.StatusBadRequest.Code(), gin.H{
			"error":  err.ErrorMessage(),
			"errors": err.ErrorMap(),
		})
	}

    err := h.InventoryService.UpdateInventoryQuantity(ctx, body.HubID, body.SellerID, body.SKUCode, body.Quantity)
    if err != nil {
        log.ErrorfWithContext(ctx, logTag+" failed to update inventory quantity %v", err)
        c.JSON(http.StatusInternalServerError.Code(), gin.H{
            "error": "Failed to update inventory quantity",
        })
        return
    }

    log.InfofWithContext(ctx, logTag+" inventory quantity updated successfully")
    c.JSON(http.StatusOK.Code(), gin.H{
        "message": "Inventory quantity updated successfully",
    })
}