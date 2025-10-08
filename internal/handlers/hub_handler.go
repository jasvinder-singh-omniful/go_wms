package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/omniful/go_commons/log"
	"github.com/omniful/go_commons/validator"
	"github.com/singhJasvinder101/go_wms/internal/services"
	"gorm.io/datatypes"

	"github.com/omniful/go_commons/http"
)

type HubHandler struct {
	HubService *services.HubService
}

func (h *HubHandler) CreateHub(c *gin.Context) {
	ctx := c.Request.Context()

	logTag := "[HubHandler][CreateHub]"
	log.InfofWithContext(ctx, logTag+" creating hub")

	var body struct {
		TenantId string `json:"tenant_id" validate:"required,alpha"`
		Name     string `json:"name" validate:"required,alpha"`
		Location string `json:"location" validate:"required"`
	}
	
	if err := c.ShouldBindJSON(&body); err != nil {
		log.ErrorfWithContext(ctx, logTag+" failed to bind json%v", err)
		c.JSON(http.StatusBadRequest.Code(), gin.H{
			"message": "please write valid input",
			"error":   err.Error(),
		})
	}
	
	if err := validator.ValidateStruct(ctx, body); err.Exists() {
		log.ErrorfWithContext(ctx, logTag+" please enter valid input %v", err)
		c.JSON(http.StatusBadRequest.Code(), gin.H{
			"message": "please write valid input",
			"error":   err.Error(),
		})
	}

	hub, err := h.HubService.CreateHub(ctx, body.TenantId, body.Name, datatypes.JSON(body.Location))
	if err != nil {
		log.ErrorfWithContext(ctx, logTag+" failed to get hub: %v", err)
        c.JSON(http.StatusNotFound.Code(), gin.H{
            "error": "Hub not found",
        })
        return
	}

	response := gin.H{
        "Id":       hub.ID,
        "TenantID": hub.TenantID,
        "Name":     hub.Name,
        "Location": hub.Location,
    }

    c.JSON(http.StatusOK.Code(), response)
}

func (h *HubHandler) GetHub(c *gin.Context){
	ctx := c.Request.Context()
    logTag := "[HubHandler][GetHub]"
    
    var body struct {
        ID uint `uri:"id" validate:"required,min=1"`
    }

	if err := c.ShouldBindJSON(&body); err != nil {
        log.ErrorfWithContext(ctx, logTag+" failed to bind json%v", err)
		c.JSON(http.StatusBadRequest.Code(), gin.H{
			"error": err.Error(),
		})
        return
    }
	
	if err := validator.ValidateStruct(ctx, body); err.Exists() {
		log.ErrorfWithContext(ctx, logTag+" please enter valid input %v", err)
        c.JSON(http.StatusBadRequest.Code(), gin.H{
            "error": err.ErrorMessage(),
            "errors": err.ErrorMap(),
        })
		return
	}

	hub, err := h.HubService.GetHubByID(ctx, body.ID)
    if err != nil {
        log.ErrorfWithContext(ctx, logTag+" failed to get hub: %v", err)
        c.JSON(http.StatusNotFound.Code(), gin.H{
            "error": "Hub not found",
        })
        return
    }

    response := &gin.H{
        "ID":       hub.ID,
        "TenantID": hub.TenantID,
        "Name":     hub.Name,
        "Location": hub.Location,
    }

    c.JSON(http.StatusOK.Code(), response)
}

func (h *HubHandler) GetAllHubs(c *gin.Context){
	ctx := c.Request.Context()
    logTag := "[HubHandler][GetAllHubs]"
    log.InfofWithContext(ctx, logTag+" fetching all hubs")

	hubs, err := h.HubService.GetAllHubs(ctx)
    if err != nil {
        log.ErrorfWithContext(ctx, logTag+" failed to get hubs: %v", err)
        c.JSON(http.StatusInternalServerError.Code(), gin.H{
            "error": "Failed to fetch hubs",
        })
        return
    }

	var response []gin.H
	for _, hub := range hubs {
		response = append(response, gin.H{
			"ID":       hub.ID,
			"TenantID": hub.TenantID,
			"Name":     hub.Name,
			"Location": hub.Location,
		})
	}

    c.JSON(http.StatusOK.Code(), gin.H{
        "hubs": response,
        "count": len(response),
    })
}


