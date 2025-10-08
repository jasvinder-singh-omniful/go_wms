package services

import (
	"context"
	"fmt"

	"github.com/omniful/go_commons/log"
	"github.com/singhJasvinder101/go_wms/internal/storage"
	"github.com/singhJasvinder101/go_wms/models"
)

type InventoryService struct {
	InventoryRepo *storage.InventoryRepo
	SKURepo       *storage.SKURepo
	HubRepo       *storage.HubRepo
}

func NewInventoryService(inventoryRepo *storage.InventoryRepo, skuRepo *storage.SKURepo, hubRepo *storage.HubRepo) *InventoryService {
	return &InventoryService{
		InventoryRepo: inventoryRepo,
		SKURepo:       skuRepo,
		HubRepo:       hubRepo,
	}
}

func (s *InventoryService) CreateInventory(ctx context.Context, tenantId, sellerId string, skuCode string, hubId int, quantity int64) (*models.Inventory, error) {
	logTag := "[InventoryService][CreateInventory]"
	log.InfofWithContext(ctx, logTag+" creating inventory for hub %d, seller %s, SKU %s", tenantId, sellerId, skuCode)

	skus, err := s.SKURepo.GetByCodes(ctx, tenantId, sellerId, []string{skuCode})
	if err != nil {
		log.ErrorfWithContext(ctx, logTag+" failed to get SKU ID %v", err)
		return nil, fmt.Errorf("failed to get SKU ID %w", err)
	}

	if len(skus) == 0 {
		return nil, fmt.Errorf("SKU not found: %s", skuCode)
	}
	skuID := skus[0].ID

	inventory := &models.Inventory{
		SKUID:    int(skuID),
		HubID:    hubId,
		Quantity: quantity,
	}

	if err := s.InventoryRepo.Create(ctx, inventory); err != nil {
		log.ErrorfWithContext(ctx, logTag+" failed to create inventory in database %v", err)
		return nil, fmt.Errorf("failed to create inventory %w", err)
	}

	log.InfofWithContext(ctx, logTag+" inventory created successfully with ID: %d", inventory.ID)
	return inventory, nil
}

func (s *InventoryService) UpsertInventory(ctx context.Context, tenantID, sellerID string, skuCode string, hubId int, quantity int64) (*models.Inventory, error) {
	logTag := "[InventoryService][UpsertInventory]"
	log.InfofWithContext(ctx, logTag+" upserting inventory for hub %d, seller %s, SKU %s", tenantID, sellerID, skuCode)

	skus, err := s.SKURepo.GetByCodes(ctx, tenantID, sellerID, []string{skuCode})
	if err != nil {
		log.ErrorfWithContext(ctx, logTag+" failed to get SKU ID %v", err)
		return nil, fmt.Errorf("failed to get SKU ID %w", err)
	}

	if len(skus) == 0 {
		return nil, fmt.Errorf("SKU not found: %s", skuCode)
	}
	skuID := skus[0].ID

	_, err = s.HubRepo.GetByID(ctx, uint(hubId))
	if err != nil {
		log.ErrorfWithContext(ctx, logTag+" failed to get Hub ID %v", err)
		return nil, fmt.Errorf("failed to get Hub ID %w", err)
	}

	inventory := &models.Inventory{
		SKUID:    int(skuID),
		HubID:    hubId,
		Quantity: quantity,
	}

	if err := s.InventoryRepo.Upsert(ctx, inventory); err != nil {
		log.ErrorfWithContext(ctx, logTag+" failed to upsert inventory in database: %v", err)
		return nil, fmt.Errorf("failed to upsert inventory: %w", err)
	}

	log.InfofWithContext(ctx, logTag+" inventory upserted successfully")
	return inventory, nil
}

func (s *InventoryService) UpdateInventoryQuantity(ctx context.Context, hubID uint, sellerID string, skuCode string, quantity int) error {
	logTag := "[InventoryService][UpdateInventoryQuantity]"
	log.InfofWithContext(ctx, logTag+" updating inventory quantities for hub %d, seller %s", hubID, sellerID)


	if err := s.InventoryRepo.UpdateQuantity(ctx, hubID, sellerID, skuCode, int(quantity)); err != nil {
		log.ErrorfWithContext(ctx, logTag+" failed to update inventory for SKU %s: %v", skuCode, err)
		return fmt.Errorf("failed to update inventory for SKU %s: %w", skuCode, err)
	}
	log.InfofWithContext(ctx, logTag+" updated inventory for SKU %s, quantity: %d", skuCode, quantity)

	return nil
}
