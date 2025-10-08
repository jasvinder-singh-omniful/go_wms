package services

import (
	"context"
	"fmt"

	"github.com/omniful/go_commons/log"
	"github.com/singhJasvinder101/go_wms/internal/storage"
	"github.com/singhJasvinder101/go_wms/models"
	"gorm.io/datatypes"
)

type SKUService struct {
    SKURepo *storage.SKURepo
}

func NewSKUService(skuRepo *storage.SKURepo) *SKUService {
    return &SKUService{
        SKURepo: skuRepo,
    }
}

func (s *SKUService) CreateSKU(ctx context.Context, tenantId, sellerId string, skuCode, name string, metadata datatypes.JSON) (*models.SKU, error) {
    logTag := "[SKUService][CreateSKU]"
    log.InfofWithContext(ctx, logTag+" creating SKU for tenant: %s, seller: %s", tenantId, sellerId)


    sku := &models.SKU{
        TenantID: tenantId,
        SellerID: sellerId,
        SKUCode:  skuCode,
        Name:     name,
        MetaData: metadata,
    }

    if err := s.SKURepo.Create(ctx, sku); err != nil {
        log.ErrorfWithContext(ctx, logTag+" failed to create SKU in database: %v", err)
        return nil, fmt.Errorf("failed to create SKU %w", err)
    }

    log.InfofWithContext(ctx, logTag+" SKU created successfully with ID: %d", sku.ID)
    return sku, nil
}

func (s *SKUService) GetSKUsByCodes(ctx context.Context, tenantID, sellerID string, skuCodes []string) ([]models.SKU, error) {
    logTag := "[SKUService][GetSKUsByIDs]"
    log.InfofWithContext(ctx, logTag+" fetching SKUs for tenant %s, seller: %s", tenantID, sellerID)


    skus, err := s.SKURepo.GetByCodes(ctx, tenantID, sellerID, skuCodes)
    if err != nil {
        log.ErrorfWithContext(ctx, logTag+" failed to fetch SKUs %v", err)
        return nil, fmt.Errorf("failed to fetch SKUs: %w", err)
    }

    log.InfofWithContext(ctx, logTag+" found %d SKUs out of %d requested", len(skus), len(skuCodes))
    return skus, nil
}