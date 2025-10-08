package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/omniful/go_commons/http"
	"github.com/omniful/go_commons/log"
	"github.com/omniful/go_commons/validator"
	"github.com/singhJasvinder101/go_wms/internal/services"
	"gorm.io/datatypes"
)

type SKUHandler struct {
	SKUService *services.SKUService
}

func NewSKUHandler(skuService *services.SKUService) *SKUHandler {
	return &SKUHandler{
		SKUService: skuService,
	}
}

func (h *SKUHandler) CreateSKU(c *gin.Context) {
	ctx := c.Request.Context()

	logTag := "[SKUHandler][CreateSKU]"
	log.InfofWithContext(ctx, logTag+" creating SKU")

	var body struct {
		TenantID string         `json:"tenant_id" validate:"required"`
		SellerID string         `json:"seller_id" validate:"required"`
		SKUCode  string         `json:"sku_code" validate:"required,min=1,max=50"`
		Name     string         `json:"name" validate:"required,min=2,max=200"`
		Metadata datatypes.JSON `json:"metadata,omitempty"`
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

	sku, err := h.SKUService.CreateSKU(ctx, body.TenantID, body.SellerID, body.SKUCode, body.Name, body.Metadata)
	if err != nil {
		log.ErrorfWithContext(ctx, logTag+" failed to create SKU: %v", err)
		c.JSON(http.StatusInternalServerError.Code(), gin.H{
			"error": "Failed to create SKU",
		})
		return
	}

	response := &gin.H{
		"ID":       sku.ID,
		"TenantID": sku.TenantID,
		"SellerID": sku.SellerID,
		"SKUCode":  sku.SKUCode,
		"Name":     sku.Name,
		"Metadata": sku.MetaData,
	}

	log.InfofWithContext(ctx, logTag+" SKU created successfully")
	c.JSON(http.StatusCreated.Code(), response)
}

func (h *SKUHandler) GetSKUsByCodes(c *gin.Context) {
	ctx := c.Request.Context()
	logTag := "[SKUHandler][GetSKUsByCodes]"
	log.InfofWithContext(ctx, logTag+" getting SKUs by codes")

	var body struct {
		TenantID string   `json:"tenant_id" validate:"required"`
		SellerID string   `json:"seller_id" validate:"required"`
		SKUCodes []string `json:"sku_codes" validate:"required,min=1,max=100,dive,required,min=1"`
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

	skus, err := h.SKUService.GetSKUsByCodes(ctx, body.TenantID, body.SellerID, body.SKUCodes)
	if err != nil {
		log.ErrorfWithContext(ctx, logTag+" failed to get SKUs %v", err)
		c.JSON(http.StatusInternalServerError.Code(), gin.H{
			"error": "Failed to fetch SKUs",
		})
		return
	}

	var response []gin.H
	for _, sku := range skus {
		response = append(response, gin.H{
			"ID":       sku.ID,
			"TenantID": sku.TenantID,
			"SellerID": sku.SellerID,
			"SKUCode":  sku.SKUCode,
			"Name":     sku.Name,
			"Metadata": sku.MetaData,
		})
	}

	c.JSON(http.StatusOK.Code(), gin.H{
		"skus":  response,
		"count": len(response),
	})
}
