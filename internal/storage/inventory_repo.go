package storage

import (
	"context"
	"fmt"

	"github.com/omniful/go_commons/log"
	"github.com/singhJasvinder101/go_wms/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type InventoryRepo struct {
	DB *Postgres
}

func NewInventoryRepo(db *Postgres) *InventoryRepo {
	return &InventoryRepo{
		DB: db,
	}
}

func (r *InventoryRepo) Create(ctx context.Context, inventory *models.Inventory) error {
	logTag := "[SKURepo][Create]"
	log.InfofWithContext(ctx, logTag+" creating inventory in db", "inventory", inventory)

	db := r.DB.Cluster.GetMasterDB(ctx)

	if err := db.Create(&inventory).Error; err != nil {
		log.ErrorfWithContext(ctx, logTag+" error when creating inventory in db", err)
		return fmt.Errorf("error when creating inventory in db %v", err)
	}

	log.InfofWithContext(ctx, logTag+" creating inventory in db", inventory)
	return nil
}

func (r *InventoryRepo) Upsert(ctx context.Context, inventory *models.Inventory) error {
	logTag := "[SKURepo][Upsert]"
	log.InfofWithContext(ctx, logTag+" updating inventory in db", "inventory", inventory)
	
	db := r.DB.Cluster.GetMasterDB(ctx)

	if err := db.Clauses(clause.OnConflict{
        Columns:   []clause.Column{{Name: "sku_id"}, {Name: "hub_id"}},
        DoUpdates: clause.AssignmentColumns([]string{"quantity", "tenant_id", "seller_id", "updated_at"}),
    }).Create(inventory).Error; err != nil {
        log.ErrorfWithContext(ctx, logTag+" error when upserting inventory in db", err)
        return fmt.Errorf("error when upserting inventory in db %v", err)
    }

	log.InfofWithContext(ctx, logTag+" updating inventory in db", inventory)
	return nil
}


func (r *InventoryRepo) GetByHubAndSeller(ctx context.Context, hubID int, sellerID string) ([]models.Inventory, error) {
	logTag := "[SKURepo][GetByHubAndSeller]"
	log.InfofWithContext(ctx, logTag+" updating sku in db", "hub_id", hubID, "seller_id", sellerID)
	
	db := r.DB.Cluster.GetSlaveDB(ctx)

	var inventory []models.Inventory
	if err := db.Where("hub_id = ? AND seller_id = ?", hubID, sellerID).Find(&inventory).Error; err != nil {
		if err == gorm.ErrRecordNotFound{
			return nil, fmt.Errorf("no record found with hub_id %d and seller_id %s", hubID, sellerID)
		}
		log.ErrorfWithContext(ctx, logTag+" error when getting inventory by hub_id and seller_id", err)
		return nil, fmt.Errorf("error when getting inventory by hub_id and seller_id %v", err)
	}

	log.InfofWithContext(ctx, logTag+" fetching inventory successfully", inventory)
	return inventory, nil
}

func (r *InventoryRepo) GetByHubSellerSKUs(ctx context.Context, hubID int, sellerID string, skuCodes []string) ([]struct{
	SKU string
	Quantity int64
}, error) {
	logTag := "[SKURepo][GetByHubSellerSKUs]"
	log.InfofWithContext(ctx, logTag+" updating sku in db", "hub_id", hubID, "seller_id", sellerID)
	
	db := r.DB.Cluster.GetSlaveDB(ctx)

	var inventory []struct{
		SKU string
		Quantity int64
	}

	query := db.Table("inventory AS i").
		Select("s.sku_code AS sku, i.quantity").
		Joins("JOIN skus AS s ON s.id = i.sku_id").
		Where("i.hub_id = ? AND i.seller_id = ?", hubID, sellerID)

	if len(skuCodes) > 0 {
        query = query.Where("s.sku_code IN ?", skuCodes)
    }

	if err := query.Find(&inventory).Error; err != nil {
		if err == gorm.ErrRecordNotFound{
			return nil, fmt.Errorf("no record found with hub_id %d and seller_id %s", hubID, sellerID)
		}
		log.ErrorfWithContext(ctx, logTag+" error when getting inventory by hub_id and seller_id", err)
		return nil, fmt.Errorf("error when getting inventory by hub_id and seller_id %v", err)
	}

	log.InfofWithContext(ctx, logTag+" fetching inventory successfully", inventory)
	return inventory, nil
}

func (r *InventoryRepo) UpdateQuantity(ctx context.Context, hubID uint, sellerID string, skuID int, quantity int) error {
	logTag := "[SKURepo][GetByHubAndSeller]"
	log.InfofWithContext(ctx, logTag+" updating sku in db", "hub_id", hubID, "seller_id", sellerID, "sku_code", skuID, "quantity", quantity)
	
	db := r.DB.Cluster.GetMasterDB(ctx)

	var count int64
    err := db.Model(&models.Inventory{}).
        Where("hub_id = ? AND seller_id = ? AND sku_id = ?", hubID, sellerID, skuID).
        Count(&count).Error
    
    if err != nil {
        log.ErrorfWithContext(ctx, logTag+" error checking inventory existence: %v", err)
        return fmt.Errorf("error checking inventory existence: %w", err)
    }

    if count == 0 {
        log.ErrorfWithContext(ctx, logTag+" no inventory record found for hub_id=%d, seller_id=%s, sku_id=%d", 
            hubID, sellerID, skuID)
        return fmt.Errorf("inventory not found for hub_id=%d, seller_id=%s, sku_id=%d", 
            hubID, sellerID, skuID)
    }


	if err := db.Model(&models.Inventory{}).
        Where("hub_id = ? AND seller_id = ? AND sku_id = ?", hubID, sellerID, skuID).
        Update("quantity", gorm.Expr("quantity + ?", quantity)) .Error; err != nil {
			log.ErrorfWithContext(ctx, logTag+" error when updating inventory by hub_id and seller_id", err)
			return fmt.Errorf("error when updating inventory by hub_id, seller_id, skuID, and quanity %v", err)
		}

	log.InfofWithContext(ctx, "inveneotry udpated successfully")
	return nil
}

