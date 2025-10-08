package services

import (
	"context"
	"fmt"

	"github.com/omniful/go_commons/log"
	"github.com/singhJasvinder101/go_wms/internal/storage"
	"github.com/singhJasvinder101/go_wms/models"
	"gorm.io/datatypes"
)

type HubService struct {
    HubRepo *storage.HubRepo
}

func NewHubService(hubRepo *storage.HubRepo) *HubService {
    return &HubService{
        HubRepo: hubRepo,
    }
}

func (s *HubService) CreateHub(ctx context.Context, tenantId string, name string, location datatypes.JSON) (*models.Hub, error) {
    logTag := "[HubService][CreateHub]"
    log.InfofWithContext(ctx, logTag+" creating hub for tenant %s", tenantId)

 
    hub := &models.Hub{
        TenantID: tenantId,
        Name:     name,
        Location: location,
    }

    if err := s.HubRepo.Create(ctx, hub); err != nil {
        log.ErrorfWithContext(ctx, logTag+" failed to create hub in database: %v", err)
        return nil, fmt.Errorf("failed to create hub %w", err)
    }

    log.InfofWithContext(ctx, logTag+" hub created successfully with ID: %d", hub.ID)
    return hub, nil
}

func (s *HubService) GetHubByID(ctx context.Context, id uint) (*models.Hub, error) {
    logTag := "[HubService][GetHubByID]"
    log.InfofWithContext(ctx, logTag+" fetching hub by ID %d", id)

    hub, err := s.HubRepo.GetByID(ctx, id)

    if err != nil {
        log.ErrorfWithContext(ctx, logTag+" failed to fetch hub %v", err)
        return nil, fmt.Errorf("failed to fetch hub %w", err)
    }
    
    return hub, nil
}

func (s *HubService) GetAllHubs(ctx context.Context) ([]models.Hub, error) {
    logTag := "[HubService][GetAllHubs]"
    log.InfofWithContext(ctx, logTag+" fetching all hubs")

    hubs, err := s.HubRepo.GetAll(ctx)
    if err != nil {
        log.ErrorfWithContext(ctx, logTag+" failed to fetch hubs: %v", err)
        return nil, fmt.Errorf("failed to fetch hubs %w", err)
    }

    log.InfofWithContext(ctx, logTag+" found %d hubs", len(hubs))
    return hubs, nil
}