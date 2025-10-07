package storage

import (
	"context"
	"fmt"

	"github.com/omniful/go_commons/log"
	"github.com/singhJasvinder101/go_wms/models"
)


type SKURepo struct {
	DB *Postgres
}


func (r *SKURepo) Create(ctx context.Context, sku *models.SKU) error {
	logTag := "[SKURepo][Create]"
	log.InfofWithContext(ctx, logTag+" creating sku in db", sku)

	db := r.DB.Cluster.GetMasterDB(ctx)

	if err := db.Create(&sku).Error; err != nil {
		log.ErrorfWithContext(ctx, logTag+" error when creating sku in db", err)
		return fmt.Errorf("error when creating sku in db %v", err)
	}

	log.InfofWithContext(ctx, logTag+" sk created successfully")
	return nil
}

func (r *SKURepo) GetByIDs(ctx context.Context, tenantID, sellerID string, skuCodes []string) ([]models.SKU, error){
	logTag := "[SKURepo][GetByIDs]"
	log.InfofWithContext(ctx, logTag+" geting sku by ids in database ", "tenant_id", tenantID, "seller_id", sellerID, "sku_codes", skuCodes)

	db := r.DB.Cluster.GetSlaveDB(ctx)

	var skus []models.SKU
	if err := db.Where("tenant_id = ? AND seller_id = ? AND sku_codes = (?)", tenantID, sellerID, skuCodes).Find(&skus).Error; err != nil {
		log.ErrorfWithContext(ctx, logTag+" error when getting sku by ids in db", err)
		return nil, fmt.Errorf("error when getting sku by ids in db %v", err)	
	}

	return skus, nil
}

